package api

import (
	"archive/tar"
	"bufio"
	"bytes"
	"fmt"
	"io"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"text/template"
	"time"

	log "github.com/sirupsen/logrus"
)

var fileToUpload = `
FROM busybox:latest
ADD {{.Src}} /.source/{{.Src}}
CMD ["cp", "-r", "/.source/{{.Src}}", "/.destination/"]
`

var folderToUpload = `
FROM busybox:latest
ADD {{.Src}} /.source/{{.Src}}
CMD ["cp", "-r", "/.source/{{.Src}}/.", "/.destination/"]
`

func addDockerfileToTar(srcPath string, tarWriter *tar.Writer, tmpWriter *bufio.Writer) error {
	fi, err := os.Lstat(srcPath)
	if err != nil {
		log.Errorf("Can't get file info: %s, error: %s", srcPath, err)
		return err
	}

	dockerfileToUpload := fileToUpload
	if fi.Mode().IsDir() {
		dockerfileToUpload = folderToUpload
	}

	tmpl, err := template.New("Dockerfile").Parse(dockerfileToUpload)
	if err != nil {
		return err
	}

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, struct {
		Src string
	}{
		Src: filepath.Base(srcPath),
	}); err != nil {
		log.Errorf("Can't execute template to upload: %s", err)
		return err
	}

	hdr, err := tar.FileInfoHeader(fi, "")
	if err != nil {
		log.Errorf("Can't get file info header: %s, error: %s", srcPath, err)
		return err
	}

	hdr.Name = DOCKERFILE
	hdr.Mode = 0100644 // Regular file + rw-r--r--
	hdr.Size = int64(buf.Len())
	hdr.ModTime = time.Now()
	hdr.Typeflag = tar.TypeReg
	hdr.Linkname = ""

	if err := tarWriter.WriteHeader(hdr); err != nil {
		log.Errorf("Can't write tar header: %s", err)
		return err
	}

	tmpWriter.Reset(tarWriter)
	defer tmpWriter.Reset(nil)
	if _, err := io.Copy(tmpWriter, &buf); err != nil {
		log.Errorf("Can't write Dockerfile to tar: %s", err)
		return err
	}
	if err := tmpWriter.Flush(); err != nil {
		log.Errorf("Can't flush Dockerfile to tar: %s", err)
		return err
	}

	return nil
}

func (client *DockerClient) Upload(srcPath string, quiet bool) (string, error) {
	v := url.Values{}
	v.Set("rm", "1")
	if quiet {
		v.Set("q", "1")
	}

	uri := fmt.Sprintf("/v%s/build?%s", API_VERSION, v.Encode())

	srcPath = filepath.Clean(srcPath)

	var (
		rootDir  = filepath.Dir(srcPath)
		filename = filepath.Base(srcPath)
	)

	srcFi, err := os.Lstat(srcPath)
	if err != nil {
		return "", err
	}

	fmt.Fprintf(client.out, "Sending build context to Docker daemon\n")
	if !quiet && (log.GetLevel() < log.InfoLevel) {
		fmt.Fprintf(client.out, "---> ")
	}

	pipeReader, pipeWriter := io.Pipe()

	go func() {
		var (
			files int64 = 0
			total int64 = 0
		)

		bufWriter := bufio.NewWriterSize(pipeWriter, 32*1024)
		tarWriter := tar.NewWriter(bufWriter)
		tmpWriter := bufio.NewWriterSize(nil, 32*1024)
		defer tmpWriter.Reset(nil)

		seen := make(map[string]bool)

		if err := addDockerfileToTar(srcPath, tarWriter, tmpWriter); err != nil {
			log.Debugf("Can't add Dockerfile: %s", err)
			return
		}

		addFileToTar := func(filePath string, f os.FileInfo, err error) error {
			if err != nil {
				log.Debugf("Can't stat file %s, error: %s", filePath, err)
				return nil
			}

			relFilePath, err := filepath.Rel(rootDir, filePath)
			if err != nil || (relFilePath == "." && f.IsDir()) {
				return nil
			}

			skip := false

			switch relFilePath {
			default:
				skip = !strings.HasPrefix(relFilePath, filename)
			case DOCKERFILE:
				skip = true
			}

			if skip {
				if f.IsDir() {
					return filepath.SkipDir
				}
				return nil
			}

			if seen[relFilePath] {
				return nil
			}
			seen[relFilePath] = true

			var size int64

			if err := func() error { // Adding a file to tar
				fi, err := os.Lstat(filePath)
				if err != nil {
					log.Errorf("Can't get file info: %s, error: %s", filePath, err)
					return err
				}

				size = fi.Size()

				link := ""
				if (fi.Mode() & os.ModeSymlink) != 0 {
					if link, err = os.Readlink(filePath); err != nil {
						log.Errorf("Can't read link to tar: %s, error: %s", filePath, err)
						return err
					}
				}

				hdr, err := tar.FileInfoHeader(fi, link)
				if err != nil {
					log.Errorf("Can't get file info header to tar: %s, error: %s", filePath, err)
					return err
				}

				name := relFilePath
				if fi.IsDir() && !strings.HasSuffix(name, "/") {
					name = name + "/"
				}
				hdr.Name = name

				if err := tarWriter.WriteHeader(hdr); err != nil {
					log.Errorf("Can't write tar header, error: %s", err)
					return err
				}

				if hdr.Typeflag == tar.TypeReg {
					file, err := os.Open(filePath)
					if err != nil {
						log.Errorf("Can't open file: %s, error: %s", filePath, err)
						return err
					}

					tmpWriter.Reset(tarWriter)
					defer tmpWriter.Reset(nil)
					_, err = io.Copy(tmpWriter, file)
					file.Close()
					if err != nil {
						log.Errorf("Can't write file to tar: %s, error: %s", filePath, err)
						return err
					}
					if err := tmpWriter.Flush(); err != nil {
						log.Errorf("Can't flush file to tar, error: %s", err)
						return err
					}
				}

				return nil
			}(); err != nil {
				log.Debugf("Can't add file %s to tar, error: %s", filePath, err)
			}

			files++
			total += size

			if !quiet && (log.GetLevel() < log.InfoLevel) {
				fmt.Fprintf(client.out, ".")
			}

			if srcFi.Mode().IsDir() {
				relFilePath, err = filepath.Rel(filename, relFilePath)
				if err != nil || (relFilePath == "." && f.IsDir()) {
					return nil
				}
			}

			log.WithFields(log.Fields{
				"": fmt.Sprintf(" %7.2f KB", float64(size)/1000),
			}).Infof("---> %s", relFilePath)

			return nil
		}

		if srcFi.Mode().IsDir() {
			filepath.Walk(filepath.Join(rootDir, "."), addFileToTar)
		} else {
			if err := addFileToTar(srcPath, srcFi, nil); err != nil {
				log.Debugf("Can't add file %s to tar, error: %s", srcPath, err)
			}
		}

		if err := tarWriter.Close(); err != nil {
			log.Debugf("Can't close tar writer: %s", err)
		}

		bufWriter.Flush()
		if err := pipeWriter.Close(); err != nil {
			log.Debugf("Can't close pipe writer: %s", err)
		}

		if !quiet && (log.GetLevel() < log.InfoLevel) {
			fmt.Fprintf(client.out, "\n")
		}
		fmt.Fprintf(client.out, "---> Sent %d file(s), %.2f KB\n", files, float64(total)/1000)
	}()

	headers := map[string]string{}
	headers["Content-type"] = "application/tar"

	return client.doStreamRequest("POST", uri, pipeReader, headers, quiet)
}
