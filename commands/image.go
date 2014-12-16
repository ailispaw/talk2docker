package commands

import (
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/spf13/cobra"
	api "github.com/yungsang/dockerclient"
	"github.com/yungsang/tablewriter"
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

// Define at ps.go
// var boolAll, boolQuiet, boolNoHeader bool

func init() {
	cmdIs.Flags().BoolVarP(&boolAll, "all", "a", false, "Show all images. Only named/taged and leaf images are shown by default.")
	cmdIs.Flags().BoolVarP(&boolQuiet, "quiet", "q", false, "Only display numeric IDs")
	cmdIs.Flags().BoolVarP(&boolNoHeader, "no-header", "n", false, "Omit the header")

	cmdListImages.Flags().BoolVarP(&boolAll, "all", "a", false, "Show all images. Only named/taged and leaf images are shown by default.")
	cmdListImages.Flags().BoolVarP(&boolQuiet, "quiet", "q", false, "Only display numeric IDs")
	cmdListImages.Flags().BoolVarP(&boolNoHeader, "no-header", "n", false, "Omit the header")

	cmdImage.AddCommand(cmdListImages)
}

func listImages(ctx *cobra.Command, args []string) {
	docker, err := client.GetDockerClient(configPath, hostName)
	if err != nil {
		log.Fatal(err)
	}

	images, err := docker.ListImages(boolAll)
	if err != nil {
		log.Fatal(err)
	}

	var matchName = ""
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
		roots := make([]*api.Image, 0)
		parents := make(map[string][]*api.Image)
		for _, image := range images {
			if image.ParentId == "" {
				roots = append(roots, image)
			} else {
				if children, exists := parents[image.ParentId]; exists {
					parents[image.ParentId] = append(children, image)
				} else {
					children := make([]*api.Image, 0)
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
					fmt.Sprintf("%.3f", float64(image.VirtualSize)/1000000.0),
					FormatDateTime(time.Unix(image.Created, 0)),
				}
				items = append(items, out)
			}
		}
	}

	var header = []string{
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
	}
	table.SetBorder(false)
	table.AppendBulk(items)
	table.Render()
}

func walkTree(images []*api.Image, parents map[string][]*api.Image, prefix string, items [][]string) [][]string {
	printImage := func(prefix string, image *api.Image, isLeaf bool) {
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
			fmt.Sprintf("%.3f", float64(image.VirtualSize)/1000000.0),
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
