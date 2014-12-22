package commands

import (
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/spf13/cobra"
	"github.com/yungsang/tablewriter"
	"github.com/yungsang/talk2docker/api"
	"github.com/yungsang/talk2docker/client"
)

var cmdIs = &cobra.Command{
	Use:     "ls [NAME[:TAG]]",
	Aliases: []string{"images"},
	Short:   "List images",
	Long:    appName + " ls - List images",
	Run:     listImages,
}

var cmdImage = &cobra.Command{
	Use:     "image [command]",
	Aliases: []string{"img"},
	Short:   "Manage images",
	Long:    appName + " image - Manage images",
	Run: func(ctx *cobra.Command, args []string) {
		ctx.Usage()
	},
}

var cmdListImages = &cobra.Command{
	Use:     "list [NAME[:TAG]]",
	Aliases: []string{"ls"},
	Short:   "List images",
	Long:    appName + " image list - List images",
	Run:     listImages,
}

var cmdPullImage = &cobra.Command{
	Use:   "pull <NAME[:TAG]>",
	Short: "Pull an image",
	Long:  appName + " image pull - Pull an image",
	Run:   pullImage,
}

func init() {
	cmdIs.Flags().BoolVarP(&boolAll, "all", "a", false, "Show all images. Only named/taged and leaf images are shown by default.")
	cmdIs.Flags().BoolVarP(&boolQuiet, "quiet", "q", false, "Only display numeric IDs")
	cmdIs.Flags().BoolVarP(&boolNoHeader, "no-header", "n", false, "Omit the header")

	cmdListImages.Flags().BoolVarP(&boolAll, "all", "a", false, "Show all images. Only named/taged and leaf images are shown by default.")
	cmdListImages.Flags().BoolVarP(&boolQuiet, "quiet", "q", false, "Only display numeric IDs")
	cmdListImages.Flags().BoolVarP(&boolNoHeader, "no-header", "n", false, "Omit the header")

	cmdPullImage.Flags().BoolVarP(&boolAll, "all", "a", false, "Pull all tagged images in the repository. Only the \"latest\" tagged image is pulled by default.")

	cmdImage.AddCommand(cmdListImages)
	cmdImage.AddCommand(cmdPullImage)
}

func listImages(ctx *cobra.Command, args []string) {
	docker, err := client.NewDockerClient(configPath, hostName)
	if err != nil {
		log.Fatal(err)
	}

	images, err := docker.ListImages(boolAll, nil)
	if err != nil {
		log.Fatal(err)
	}

	matchName := ""
	if len(args) > 0 {
		matchName = args[0]
	}

	matchImageByName := func(tags []string, name string) bool {
		_name := strings.Split(name, ":")

		for _, tag := range tags {
			_tag := strings.Split(tag, ":")
			if _tag[0] == _name[0] {
				if (len(_name) < 2) || (_tag[1] == _name[1]) {
					return true
				}
			}
		}

		return false
	}

	if boolQuiet {
		for _, image := range images {
			if (matchName == "") || matchImageByName(image.RepoTags, matchName) {
				fmt.Println(Truncate(image.Id, 12))
			}
		}
		return
	}

	var items [][]string

	if boolAll {
		roots := []api.Image{}
		parents := map[string][]api.Image{}
		for _, image := range images {
			if image.ParentId == "" {
				roots = append(roots, image)
			} else {
				if children, exists := parents[image.ParentId]; exists {
					parents[image.ParentId] = append(children, image)
				} else {
					children := []api.Image{}
					parents[image.ParentId] = append(children, image)
				}
			}
		}

		items = walkTree(roots, parents, "", items)
	} else {
		for _, image := range images {
			if (matchName == "") || matchImageByName(image.RepoTags, matchName) {
				name := strings.Join(image.RepoTags, ", ")
				if name == "<none>:<none>" {
					name = "<none>"
				}
				out := []string{
					Truncate(image.Id, 12),
					FormatNonBreakingString(name),
					FormatFloat(float64(image.VirtualSize) / 1000000),
					FormatDateTime(time.Unix(image.Created, 0)),
				}
				items = append(items, out)
			}
		}
	}

	header := []string{
		"ID",
		"Name:Tags",
		"Size(MB)",
	}
	if !boolAll {
		header = append(header, "Created at")
	}

	table := tablewriter.NewWriter(os.Stdout)
	if !boolNoHeader {
		table.SetHeader(header)
	} else {
		table.SetBorder(false)
	}
	table.AppendBulk(items)
	table.Render()
}

func walkTree(images []api.Image, parents map[string][]api.Image, prefix string, items [][]string) [][]string {
	printImage := func(prefix string, image api.Image, isLeaf bool) {
		name := strings.Join(image.RepoTags, ", ")
		if name == "<none>:<none>" {
			if isLeaf {
				name = "<none>"
			} else {
				name = ""
			}
		}
		out := []string{
			FormatNonBreakingString(fmt.Sprintf("%s %s", prefix, Truncate(image.Id, 12))),
			FormatNonBreakingString(name),
			FormatFloat(float64(image.VirtualSize) / 1000000),
		}
		items = append(items, out)
	}

	length := len(images)
	if length > 1 {
		for index, image := range images {
			if (index + 1) == length {
				subimages, exists := parents[image.Id]
				printImage(prefix+"└", image, !exists)
				if exists {
					items = walkTree(subimages, parents, prefix+" ", items)
				}
			} else {
				subimages, exists := parents[image.Id]
				printImage(prefix+"├", image, !exists)
				if exists {
					items = walkTree(subimages, parents, prefix+"│", items)
				}
			}
		}
	} else {
		for _, image := range images {
			subimages, exists := parents[image.Id]
			printImage(prefix+"└", image, !exists)
			if exists {
				items = walkTree(subimages, parents, prefix+" ", items)
			}
		}
	}
	return items
}

func pullImage(ctx *cobra.Command, args []string) {
	if len(args) < 1 {
		fmt.Println("Needs an argument <NAME> to pull")
		ctx.Usage()
		return
	}

	var (
		name = args[0]
		tag  = ""
	)

	n := strings.LastIndex(name, ":")
	if n >= 0 {
		if !strings.Contains(name[n+1:], "/") {
			tag = name[n+1:]
			name = name[:n]
		}
	}

	if tag == "" {
		tag = "latest"
	}

	repository := name + ":" + tag

	if boolAll {
		repository = name
	}

	docker, err := client.NewDockerClient(configPath, hostName)
	if err != nil {
		log.Fatal(err)
	}

	info, err := docker.Info()
	if err != nil {
		log.Fatal(err)
	}

	config, err := client.LoadConfig(configPath)
	if err != nil {
		log.Fatal(err)
	}

	server, err := config.GetIndexServer(info.IndexServerAddress)
	if (err != nil) || (server.Auth == "") {
		log.Fatal("Please login prior to pulling an image.")
	}

	err = docker.PullImage(repository, server.Auth)
	if err != nil {
		log.Fatal(err)
	}
}
