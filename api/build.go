package api

import (
	"archive/tar"
	"bufio"
	"fmt"
	"io"
	"io/ioutil"
	"net/url"
	"os"
	"path"
	"path/filepath"
	"strings"

	log "github.com/sirupsen/logrus"
)

func (client *DockerClient) BuildImage(root, tag string) error {
	v := url.Values{}
	v.Set("rm", "1")
	if tag != "" {
		v.Set("t", tag)
	}

	uri := fmt.Sprintf("/v%s/build?%s", API_VERSION, v.Encode())

	filename := path.Join(root, "Dockerfile")
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		return fmt.Errorf("No Dockerfile found in %s", root)
	}

	ignore, err := ioutil.ReadFile(path.Join(root, ".dockerignore"))
	if err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("Error reading .dockerignore: %s", err)
	}

	var excludes []string
	for _, pattern := range strings.Split(string(ignore), "\n") {
		pattern = strings.TrimSpace(pattern)
		if pattern == "" {
			continue
		}
		pattern = filepath.Clean(pattern)
		ok, err := filepath.Match(pattern, "Dockerfile")
		if err != nil {
			return fmt.Errorf("Bad .dockerignore pattern: %s, error: %s", pattern, err)
		}
		if ok {
			return fmt.Errorf("Dockerfile was excluded by .dockerignore pattern %s", pattern)
		}
		excludes = append(excludes, pattern)
	}

	fmt.Fprintf(client.out, "Sending build context to Docker daemon\n")
	if log.GetLevel() < log.InfoLevel {
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

		filepath.Walk(filepath.Join(root, "."), func(filePath string, f os.FileInfo, err error) error {
			if err != nil {
				log.Debugf("Can't stat file %s, error: %s", filePath, err)
				return nil
			}

			relFilePath, err := filepath.Rel(root, filePath)
			if err != nil || (relFilePath == "." && f.IsDir()) {
				return nil
			}

			skip, err := func() (bool, error) { // Excluding
				for _, exclude := range excludes {
					matched, err := filepath.Match(exclude, relFilePath)
					if err != nil {
						log.Errorf("Error matching: %s, pattern: %s", relFilePath, exclude)
						return false, err
					}
					if matched {
						if filepath.Clean(relFilePath) == "." {
							log.Errorf("Can't exclude whole path, excluding pattern: %s", exclude)
							continue
						}
						return true, nil
					}
				}
				return false, nil
			}()
			if err != nil {
				log.Debugf("Error matching: %s, %s", relFilePath, err)
				return err
			}
			if skip {
				log.WithField("", " Skipped").Debugf("---> %s", relFilePath)
				if f.IsDir() {
					return filepath.SkipDir
				}
				return nil
			}

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
					err = tmpWriter.Flush()
					if err != nil {
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

			if log.GetLevel() < log.InfoLevel {
				fmt.Fprintf(client.out, ".")
			}
			log.WithFields(log.Fields{
				"": fmt.Sprintf(" %7.2f KB", float64(size)/1000),
			}).Infof("---> %s", relFilePath)

			return nil
		})

		if err := tarWriter.Close(); err != nil {
			log.Debugf("Can't close tar writer: %s", err)
		}

		bufWriter.Flush()
		if err := pipeWriter.Close(); err != nil {
			log.Debugf("Can't close pipe writer: %s", err)
		}

		if log.GetLevel() < log.InfoLevel {
			fmt.Fprintf(client.out, "\n")
		}
		fmt.Fprintf(client.out, "---> Sent %d file(s), %.2f KB\n", files, float64(total)/1000)
	}()

	headers := map[string]string{}
	headers["Content-type"] = "application/tar"

	return client.doStreamRequest("POST", uri, pipeReader, headers)
}