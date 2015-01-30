package commands

import (
	"encoding/json"
	"fmt"
	"log"
	"path/filepath"
	"sort"
	"strings"

	"github.com/spf13/cobra"
	"github.com/yungsang/tablewriter"

	"github.com/ailispaw/talk2docker/api"
	"github.com/ailispaw/talk2docker/client"
)

// https://github.com/docker/docker/blob/master/daemon%2Fvolumes.go#L21
type Mount struct {
	hostPath      string
	MountToPath   string
	ContainerId   string
	ContainerName string
	Writable      bool
}

// https://github.com/docker/docker/blob/master/volumes%2Fvolume.go#L16
type Volume struct {
	ID          string
	Path        string
	IsBindMount bool
	Writable    bool

	MountedOn []*Mount

	configPath string
}

type Volumes []*Volume

func (volumes Volumes) Len() int {
	return len(volumes)
}

func (volumes Volumes) Swap(i, j int) {
	volumes[i], volumes[j] = volumes[j], volumes[i]
}

func (volumes Volumes) Less(i, j int) bool {
	return volumes[i].Path < volumes[j].Path
}

func (volumes Volumes) Find(id string) *Volume {
	l := len(id)
	for _, volume := range volumes {
		if len(volume.ID) < l {
			continue
		}
		if id == volume.ID[:l] {
			return volume
		}
	}
	return nil
}

var cmdVs = &cobra.Command{
	Use:     "vs",
	Aliases: []string{"volumes"},
	Short:   "List volumes",
	Long:    APP_NAME + " vs - List volumes",
	Run:     listVolumes,
}

var cmdVolume = &cobra.Command{
	Use:     "volume [command]",
	Aliases: []string{"vol"},
	Short:   "Manage volumes",
	Long:    APP_NAME + " volume - Manage volumes",
	Run: func(ctx *cobra.Command, args []string) {
		ctx.Help()
	},
}

var cmdListVolumes = &cobra.Command{
	Use:     "list",
	Aliases: []string{"ls"},
	Short:   "List volumes",
	Long:    APP_NAME + " volume list - List volumes",
	Run:     listVolumes,
}

var cmdInspectVolumes = &cobra.Command{
	Use:     "inspect <ID>...",
	Aliases: []string{"ins", "info"},
	Short:   "Inspect volumes",
	Long:    APP_NAME + " volume inspect - Inspect volumes",
	Run:     inspectVolumes,
}

var cmdRemoveVolumes = &cobra.Command{
	Use:     "remove <ID>...",
	Aliases: []string{"rm"},
	Short:   "Remove volumes",
	Long:    APP_NAME + " volume remove - Remove volumes",
	Run:     removeVolumes,
}

var cmdExportVolume = &cobra.Command{
	Use:   "export <ID>",
	Short: "Stream the contents of a volume as a tar archive",
	Long:  APP_NAME + " volume export - Stream the contents of a volume as a tar archive",
	Run:   exportVolume,
}

func init() {
	flags := cmdVs.Flags()
	flags.BoolVarP(&boolAll, "all", "a", false, "Show all volumes. Only active volumes are shown by default.")
	flags.BoolVarP(&boolQuiet, "quiet", "q", false, "Only display numeric IDs")
	flags.BoolVarP(&boolNoHeader, "no-header", "n", false, "Omit the header")

	flags = cmdListVolumes.Flags()
	flags.BoolVarP(&boolAll, "all", "a", false, "Show all volumes. Only active volumes are shown by default.")
	flags.BoolVarP(&boolQuiet, "quiet", "q", false, "Only display numeric IDs")
	flags.BoolVarP(&boolNoHeader, "no-header", "n", false, "Omit the header")
	cmdVolume.AddCommand(cmdListVolumes)

	cmdVolume.AddCommand(cmdInspectVolumes)

	cmdVolume.AddCommand(cmdRemoveVolumes)

	cmdVolume.AddCommand(cmdExportVolume)
}

func listVolumes(ctx *cobra.Command, args []string) {
	volumes, err := getVolumes(ctx)
	if err != nil {
		log.Fatal(err)
	}

	sort.Sort(volumes)

	var _volumes Volumes
	for _, volume := range volumes {
		if len(volume.MountedOn) > 0 {
			_volumes = append(_volumes, volume)
		} else if boolAll && !volume.IsBindMount {
			_volumes = append(_volumes, volume)
		}
	}
	volumes = _volumes

	if boolQuiet {
		for _, volume := range volumes {
			ctx.Println(Truncate(volume.ID, 12))
		}
		return
	}

	if boolYAML || boolJSON {
		if err := FormatPrint(ctx.Out(), volumes); err != nil {
			log.Fatal(err)
		}
		return
	}

	formatNames := func(mounts []*Mount) string {
		names := []string{}
		for _, mount := range mounts {
			var name string
			if mount.Writable {
				name = fmt.Sprintf("%s:%s", mount.ContainerName, mount.MountToPath)
			} else {
				name = fmt.Sprintf("%s:%s:ro", mount.ContainerName, mount.MountToPath)
			}
			names = append(names, name)
		}
		return strings.Join(names, ", ")
	}

	var items [][]string
	for _, volume := range volumes {
		out := []string{
			Truncate(volume.ID, 12),
			formatNames(volume.MountedOn),
			volume.Path,
		}
		items = append(items, out)
	}

	header := []string{
		"ID",
		"Mounted On",
		"Path",
	}

	PrintInTable(ctx.Out(), header, items, 0, tablewriter.ALIGN_DEFAULT)
}

func inspectVolumes(ctx *cobra.Command, args []string) {
	if len(args) < 1 {
		ctx.Println("Needs an argument <ID> at least to inspect")
		ctx.Usage()
		return
	}

	volumes, err := getVolumes(ctx)
	if err != nil {
		log.Fatal(err)
	}

	var _volumes Volumes
	var gotError = false

	for _, id := range args {
		if volume := volumes.Find(id); volume == nil {
			log.Printf("No such volume: %s\n", id)
			gotError = true
		} else {
			_volumes = append(_volumes, volume)
		}
	}

	if len(_volumes) > 0 {
		if err := FormatPrint(ctx.Out(), _volumes); err != nil {
			log.Fatal(err)
		}
	}

	if gotError {
		log.Fatal("Error: failed to inspect one or more volumes")
	}
}

func removeVolumes(ctx *cobra.Command, args []string) {
	if len(args) < 1 {
		ctx.Println("Needs an argument <ID> at least to inspect")
		ctx.Usage()
		return
	}

	volumes, err := getVolumes(ctx)
	if err != nil {
		log.Fatal(err)
	}

	var gotError = false

	for _, id := range args {
		volume := volumes.Find(id)
		if volume == nil {
			log.Printf("No such volume: %s\n", id)
			gotError = true
			continue
		}

		if len(volume.MountedOn) > 0 {
			log.Printf("The volume is in use, cannot remove: %s\n", volume.ID)
			gotError = true
			continue
		}

		if volume.IsBindMount {
			log.Printf("The volume is bound, cannot remove: %s\n", volume.ID)
			gotError = true
			continue
		}

		if err := removeVolume(ctx, volume); err != nil {
			log.Println(err)
			gotError = true
		} else {
			ctx.Println(volume.ID)
		}
	}

	if gotError {
		log.Fatal("Error: failed to remove one or more volumes")
	}
}

func getVolumes(ctx *cobra.Command) (Volumes, error) {
	docker, err := client.NewDockerClient(configPath, hostName, ctx.Out())
	if err != nil {
		return nil, err
	}

	info, err := docker.Info()
	if err != nil {
		return nil, err
	}

	rootDir := "/var/lib/docker"

	if (info.Debug != 0) && (info.DockerRootDir != "") {
		rootDir = info.DockerRootDir
	} else {
		for _, pair := range info.DriverStatus {
			if pair[0] == "Root Dir" {
				rootDir = filepath.Dir(pair[0])
			}
		}
	}

	path := filepath.Join(rootDir, "/volumes")

	var (
		config     api.Config
		hostConfig api.HostConfig
	)

	config.Cmd = []string{"/bin/sh", "-c", "awk '{print $0}' /.docker_volumes/*/config.json"}
	config.Image = "busybox:latest"

	hostConfig.Binds = []string{path + ":/.docker_volumes:ro"}

	var cid string
	cid, err = docker.CreateContainer("", config, hostConfig)
	if err != nil {
		if apiErr, ok := err.(api.Error); ok && (apiErr.StatusCode == 404) {
			if err := pullImageInSilence(ctx, config.Image); err != nil {
				return nil, err
			}

			cid, err = docker.CreateContainer("", config, hostConfig)
			if err != nil {
				return nil, err
			}
		} else {
			return nil, err
		}
	}
	defer docker.RemoveContainer(cid, true)

	if err := docker.StartContainer(cid); err != nil {
		return nil, err
	}

	if _, err := docker.WaitContainer(cid); err != nil {
		return nil, err
	}

	logs, err := docker.GetContainerLogs(cid, false, true, true, false, 0)
	if err != nil {
		return nil, err
	}

	if logs[0] == "" {
		return nil, nil
	}

	jsonVolumes := strings.Split(strings.TrimSpace(logs[0]), "\n")

	var volumes Volumes
	for _, v := range jsonVolumes {
		volume := &Volume{}
		if err := json.Unmarshal([]byte(v), volume); err != nil {
			return nil, err
		}
		volume.configPath = filepath.Join(path, "/"+volume.ID)
		volumes = append(volumes, volume)
	}

	if err := docker.RemoveContainer(cid, true); err != nil {
		return nil, err
	}

	mounts, err := getMounts(ctx)
	if err != nil {
		return nil, err
	}

	for _, volume := range volumes {
		for _, mount := range mounts {
			if mount.hostPath == volume.Path {
				volume.MountedOn = append(volume.MountedOn, mount)
			}
		}
	}

	return volumes, nil
}

func getMounts(ctx *cobra.Command) ([]*Mount, error) {
	docker, err := client.NewDockerClient(configPath, hostName, ctx.Out())
	if err != nil {
		return nil, err
	}

	containers, err := docker.ListContainers(true, false, 0, "", "", nil)
	if err != nil {
		return nil, err
	}

	var mounts []*Mount

	for _, container := range containers {
		localMounts := map[string]*Mount{}

		containerInfo, err := docker.InspectContainer(container.Id)
		if err != nil {
			log.Println(err)
			continue
		}

		for _, bind := range containerInfo.HostConfig.Binds {
			var (
				arr   = strings.Split(bind, ":")
				mount Mount
			)

			mount.ContainerId = containerInfo.Id

			switch len(arr) {
			case 1:
				mount.MountToPath = bind
				mount.Writable = true
			case 2:
				mount.hostPath = arr[0]
				mount.MountToPath = arr[1]
				mount.Writable = true
			case 3:
				mount.hostPath = arr[0]
				mount.MountToPath = arr[1]
				switch arr[2] {
				case "ro":
					mount.Writable = false
				case "rw":
					mount.Writable = true
				default:
					continue
				}
			default:
				continue
			}

			mount.ContainerName = strings.TrimPrefix(containerInfo.Name, "/")

			localMounts[mount.MountToPath] = &mount
		}

		for mountToPath, hostPath := range containerInfo.Volumes {
			if _, exists := localMounts[mountToPath]; !exists {
				localMounts[mountToPath] = &Mount{
					hostPath:      hostPath,
					MountToPath:   mountToPath,
					ContainerId:   containerInfo.Id,
					ContainerName: strings.TrimPrefix(containerInfo.Name, "/"),
					Writable:      containerInfo.VolumesRW[mountToPath],
				}
			}
		}

		for _, mount := range localMounts {
			mounts = append(mounts, mount)
		}
	}

	return mounts, nil
}

func removeVolume(ctx *cobra.Command, volume *Volume) error {
	var (
		config     api.Config
		hostConfig api.HostConfig
	)

	config.Cmd = []string{"/bin/sh", "-c", "rm -rf /.docker_volume_config/" + volume.ID}
	config.Image = "busybox:latest"

	hostConfig.Binds = []string{filepath.Dir(volume.configPath) + ":/.docker_volume_config"}

	if !volume.IsBindMount {
		config.Cmd[2] = config.Cmd[2] + (" && rm -rf /.docker_volume/" + volume.ID)

		hostConfig.Binds = append(hostConfig.Binds, filepath.Dir(volume.Path)+":/.docker_volume")
	}

	docker, err := client.NewDockerClient(configPath, hostName, ctx.Out())
	if err != nil {
		return err
	}

	var cid string
	cid, err = docker.CreateContainer("", config, hostConfig)
	if err != nil {
		if apiErr, ok := err.(api.Error); ok && (apiErr.StatusCode == 404) {
			if err := pullImageInSilence(ctx, config.Image); err != nil {
				return err
			}

			cid, err = docker.CreateContainer("", config, hostConfig)
			if err != nil {
				return err
			}
		} else {
			return err
		}
	}
	defer docker.RemoveContainer(cid, true)

	if err := docker.StartContainer(cid); err != nil {
		return err
	}

	if _, err := docker.WaitContainer(cid); err != nil {
		return err
	}

	return nil
}

func exportVolume(ctx *cobra.Command, args []string) {
	if len(args) < 1 {
		ctx.Println("Needs an argument <ID> to export")
		ctx.Usage()
		return
	}

	volumes, err := getVolumes(ctx)
	if err != nil {
		log.Fatal(err)
	}

	volume := volumes.Find(args[0])
	if volume == nil {
		log.Fatalf("No such volume: %s\n", args[0])
	}

	var (
		config     api.Config
		hostConfig api.HostConfig
	)

	config.Image = "busybox:latest"

	hostConfig.Binds = []string{volume.Path + ":/" + volume.ID}

	docker, err := client.NewDockerClient(configPath, hostName, ctx.Out())
	if err != nil {
		log.Fatal(err)
	}

	var cid string
	cid, err = docker.CreateContainer("", config, hostConfig)
	if err != nil {
		if apiErr, ok := err.(api.Error); ok && (apiErr.StatusCode == 404) {
			if err := pullImageInSilence(ctx, config.Image); err != nil {
				log.Fatal(err)
			}

			cid, err = docker.CreateContainer("", config, hostConfig)
			if err != nil {
				log.Fatal(err)
			}
		} else {
			log.Fatal(err)
		}
	}
	defer docker.RemoveContainer(cid, true)

	if err := docker.CopyContainer(cid, "/"+volume.ID); err != nil {
		log.Fatal(err)
	}
}
