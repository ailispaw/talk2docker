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

func ListImages(ctx *cobra.Command, args []string) {
	docker, err := client.GetDockerClient(ctx)
	if err != nil {
		log.Fatal(err)
	}

	images, err := docker.ListImages(GetBoolFlag(ctx, "all"))
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

	if GetBoolFlag(ctx, "quiet") {
		for _, image := range images {
			if (matchName == "") || matchImageByName(image.RepoTags, matchName) {
				fmt.Println(Truncate(image.Id, 12))
			}
		}
		return
	}

	var items [][]string

	if GetBoolFlag(ctx, "all") {
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
				name := strings.Join(image.RepoTags, ",\u00a0")
				if name == "<none>:<none>" {
					name = "<none>"
				}
				out := []string{
					Truncate(image.Id, 12),
					name,
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
	if !GetBoolFlag(ctx, "all") {
		header = append(header, "Created at")
	}

	table := tablewriter.NewWriter(os.Stdout)
	if !GetBoolFlag(ctx, "no-header") {
		table.SetHeader(header)
	}
	table.SetBorder(false)
	table.AppendBulk(items)
	table.Render()
}

func walkTree(images []*api.Image, parents map[string][]*api.Image, prefix string, items [][]string) [][]string {
	printImage := func(prefix string, image *api.Image, isLeaf bool) {
		name := strings.Join(image.RepoTags, ",\u00a0")
		if name == "<none>:<none>" {
			if isLeaf {
				name = "<none>"
			} else {
				name = ""
			}
		}
		out := []string{
			fmt.Sprintf("%s%s%s", prefix, "\u00a0", Truncate(image.Id, 12)),
			name,
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
					items = walkTree(subimages, parents, prefix+"\u00a0", items)
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
				items = walkTree(subimages, parents, prefix+"\u2063", items)
			}
		}
	}
	return items
}
