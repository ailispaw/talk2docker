package commands

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"

	"github.com/ailispaw/talk2docker/api"
	"github.com/ailispaw/talk2docker/client"
)

var cmdUploadToContainer = &cobra.Command{
	Use:     "upload <PATH> <(NAME|ID):PATH>",
	Aliases: []string{"import"},
	Short:   "Upload a file/folder to a container",
	Long:    APP_NAME + " container upload - Upload a file/folder to a container",
	Run:     uploadToContainer,
}

func uploadToContainer(ctx *cobra.Command, args []string) {
	if len(args) < 2 {
		ErrorExit(ctx, "Needs two arguments <PATH> to upload into <(NAME|ID):PATH>")
	}

	srcPath, err := filepath.Abs(args[0])
	if err != nil {
		log.Fatal(err)
	}

	arr := strings.Split(args[1], ":")
	if len(arr) < 2 || (arr[1] == "") {
		ErrorExit(ctx, fmt.Sprint("Needs <(NAME|ID):PATH> for the second argument"))
	}

	var (
		name = arr[0]
		path = arr[1]
	)

	f, err := os.Open(os.DevNull)
	if err != nil {
		log.Fatal(err)
	}

	docker, err := client.NewDockerClient(configPath, hostName, f)
	if err != nil {
		log.Fatal(err)
	}

	info, err := docker.Info()
	if err != nil {
		log.Fatal(err)
	}

	rootDir := ""

	for _, pair := range info.DriverStatus {
		if pair[0] == "Root Dir" {
			rootDir = pair[1]
		}
	}

	if rootDir == "" {
		log.Fatal("Can't get root dir for the container")
	}

	containerInfo, err := docker.InspectContainer(name)
	if err != nil {
		log.Fatal(err)
	}

	rootDir = filepath.Join(rootDir, "/mnt/"+containerInfo.Id)

	dstPath := filepath.Join(rootDir, path)

	ctx.Printf("Uploading %s into %s\n", args[0], args[1])

	message, err := docker.Upload(srcPath, true)
	if err != nil {
		log.Fatal(err)
	}

	var (
		config     api.Config
		hostConfig api.HostConfig
	)

	if _, err := fmt.Sscanf(message, "Successfully built %s", &config.Image); err != nil {
		log.Fatal(err)
	}

	defer docker.RemoveImage(config.Image, true, false)

	hostConfig.Binds = []string{dstPath + ":/.destination"}

	cid, err := docker.CreateContainer("", config, hostConfig)
	if err != nil {
		log.Fatal(err)
	}
	defer docker.RemoveContainer(cid, true)

	if err := docker.StartContainer(cid); err != nil {
		log.Fatal(err)
	}

	if _, err := docker.WaitContainer(cid); err != nil {
		log.Fatal(err)
	}

	ctx.Print("Successfully uploaded\n")
}

var cmdUploadToVolume = &cobra.Command{
	Use:     "upload <PATH> <ID:PATH>",
	Aliases: []string{"import"},
	Short:   "Upload a file/folder to a volume",
	Long:    APP_NAME + " volume upload - Upload a file/folder to a volume",
	Run:     uploadToVolume,
}

func uploadToVolume(ctx *cobra.Command, args []string) {
	if len(args) < 2 {
		ErrorExit(ctx, "Needs two arguments <PATH> to upload into <ID:PATH>")
	}

	srcPath, err := filepath.Abs(args[0])
	if err != nil {
		log.Fatal(err)
	}

	arr := strings.Split(args[1], ":")
	if len(arr) < 2 || (arr[1] == "") {
		ErrorExit(ctx, fmt.Sprint("Needs <ID:PATH> for the second argument"))
	}

	var (
		id   = arr[0]
		path = arr[1]
	)

	volumes, err := getVolumes(ctx)
	if err != nil {
		log.Fatal(err)
	}

	volume := volumes.Find(id)
	if volume == nil {
		log.Fatalf("No such volume: %s", id)
	}

	dstPath := filepath.Join(volume.Path, path)

	f, err := os.Open(os.DevNull)
	if err != nil {
		log.Fatal(err)
	}

	docker, err := client.NewDockerClient(configPath, hostName, f)
	if err != nil {
		log.Fatal(err)
	}

	ctx.Printf("Uploading %s into %s\n", args[0], args[1])

	message, err := docker.Upload(srcPath, true)
	if err != nil {
		log.Fatal(err)
	}

	var (
		config     api.Config
		hostConfig api.HostConfig
	)

	if _, err := fmt.Sscanf(message, "Successfully built %s", &config.Image); err != nil {
		log.Fatal(err)
	}

	defer docker.RemoveImage(config.Image, true, false)

	hostConfig.Binds = []string{dstPath + ":/.destination"}

	cid, err := docker.CreateContainer("", config, hostConfig)
	if err != nil {
		log.Fatal(err)
	}
	defer docker.RemoveContainer(cid, true)

	if err := docker.StartContainer(cid); err != nil {
		log.Fatal(err)
	}

	if _, err := docker.WaitContainer(cid); err != nil {
		log.Fatal(err)
	}

	ctx.Print("Successfully uploaded\n")
}
