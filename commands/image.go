package commands

import (
	"fmt"
	"log"
	"os"
	"regexp"
	"strings"
	"time"

	"github.com/spf13/cobra"
	"github.com/yungsang/tablewriter"
	"github.com/yungsang/talk2docker/api"
	"github.com/yungsang/talk2docker/client"
)

var (
	boolForce, boolNoPrune bool
)

var cmdIs = &cobra.Command{
	Use:     "ls [NAME[:TAG]]",
	Aliases: []string{"images"},
	Short:   "List images",
	Long:    APP_NAME + " ls - List images",
	Run:     listImages,
}

var cmdImage = &cobra.Command{
	Use:     "image [command]",
	Aliases: []string{"img"},
	Short:   "Manage images",
	Long:    APP_NAME + " image - Manage images",
	Run: func(ctx *cobra.Command, args []string) {
		ctx.Usage()
	},
}

var cmdListImages = &cobra.Command{
	Use:     "list [NAME[:TAG]]",
	Aliases: []string{"ls"},
	Short:   "List images",
	Long:    APP_NAME + " image list - List images",
	Run:     listImages,
}

var cmdPullImage = &cobra.Command{
	Use:   "pull <NAME[:TAG]>",
	Short: "Pull an image",
	Long:  APP_NAME + " image pull - Pull an image",
	Run:   pullImage,
}

var cmdTagImage = &cobra.Command{
	Use:   "tag <NAME[:TAG]|ID> <NAME[:TAG]>",
	Short: "Tag an image",
	Long:  APP_NAME + " image tag - Tag an image",
	Run:   tagImage,
}

var cmdShowImageHistory = &cobra.Command{
	Use:     "history <NAME[:TAG]|ID>",
	Aliases: []string{"hist"},
	Short:   "Show the histry of an image",
	Long:    APP_NAME + " history tag - Show the histry of an image",
	Run:     showImageHistory,
}

var cmdRemoveImages = &cobra.Command{
	Use:     "remove <NAME[:TAG]|ID>...",
	Aliases: []string{"rm"},
	Short:   "Remove images",
	Long:    APP_NAME + " image remove - Remove images",
	Run:     removeImages,
}

func init() {
	cmdIs.Flags().BoolVarP(&boolAll, "all", "a", false, "Show all images. Only named/taged and leaf images are shown by default.")
	cmdIs.Flags().BoolVarP(&boolQuiet, "quiet", "q", false, "Only display numeric IDs")
	cmdIs.Flags().BoolVarP(&boolNoHeader, "no-header", "n", false, "Omit the header")

	cmdListImages.Flags().BoolVarP(&boolAll, "all", "a", false, "Show all images. Only named/taged and leaf images are shown by default.")
	cmdListImages.Flags().BoolVarP(&boolQuiet, "quiet", "q", false, "Only display numeric IDs")
	cmdListImages.Flags().BoolVarP(&boolNoHeader, "no-header", "n", false, "Omit the header")

	cmdPullImage.Flags().BoolVarP(&boolAll, "all", "a", false, "Pull all tagged images in the repository. Only the \"latest\" tagged image is pulled by default.")

	cmdTagImage.Flags().BoolVarP(&boolForce, "force", "f", false, "Force to tag")

	cmdShowImageHistory.Flags().BoolVarP(&boolAll, "all", "a", false, "Show all build instructions")
	cmdShowImageHistory.Flags().BoolVarP(&boolNoHeader, "no-header", "n", false, "Omit the header")

	cmdRemoveImages.Flags().BoolVarP(&boolForce, "force", "f", false, "Force removal of the images")
	cmdRemoveImages.Flags().BoolVarP(&boolNoPrune, "no-prune", "n", false, "Do not delete untagged parents")

	cmdImage.AddCommand(cmdListImages)
	cmdImage.AddCommand(cmdPullImage)
	cmdImage.AddCommand(cmdTagImage)
	cmdImage.AddCommand(cmdShowImageHistory)
	cmdImage.AddCommand(cmdRemoveImages)
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
		arrName := strings.Split(name, ":")

		for _, tag := range tags {
			arrTag := strings.Split(tag, ":")
			if arrTag[0] == arrName[0] {
				if (len(arrName) < 2) || (arrTag[1] == arrName[1]) {
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

	registry, name, tag, err := client.ParseRepositoryName(args[0])
	if err != nil {
		log.Fatal(err)
	}

	repository := name + ":" + tag

	if boolAll {
		repository = name
	}

	if registry != "" {
		repository = registry + "/" + repository
	}

	config, err := client.LoadConfig(configPath)
	if err != nil {
		log.Fatal(err)
	}

	if registry == "" {
		registry = client.INDEX_SERVER
	}

	registryConfig, err := config.GetRegistry(registry)
	// Some custom registries may not be needed to login.
	//	if (err != nil) || (registryConfig.Auth == "") {
	//		log.Fatal("Please login prior to pulling an image.")
	//	}

	docker, err := client.NewDockerClient(configPath, hostName)
	if err != nil {
		log.Fatal(err)
	}

	err = docker.PullImage(repository, registryConfig.Auth)
	if err != nil {
		log.Fatal(err)
	}
}

func tagImage(ctx *cobra.Command, args []string) {
	if len(args) < 2 {
		fmt.Println("Needs two arguments <IMAGE-NAME[:TAG] or IMAGE-ID> <NEW-NAME[:TAG]>")
		ctx.Usage()
		return
	}

	registry, name, tag, err := client.ParseRepositoryName(args[1])
	if err != nil {
		log.Fatal(err)
	}

	if registry != "" {
		name = registry + "/" + name
	}

	docker, err := client.NewDockerClient(configPath, hostName)
	if err != nil {
		log.Fatal(err)
	}

	err = docker.TagImage(args[0], name, tag, boolForce)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Tagged %s as %s:%s\n", args[0], name, tag)
}

func showImageHistory(ctx *cobra.Command, args []string) {
	if len(args) < 1 {
		fmt.Println("Needs an argument <IMAGE-NAME[:TAG] or IMAGE-ID>")
		ctx.Usage()
		return
	}

	docker, err := client.NewDockerClient(configPath, hostName)
	if err != nil {
		log.Fatal(err)
	}

	history, err := docker.GetImageHistory(args[0])
	if err != nil {
		log.Fatal(err)
	}

	var items [][]string

	for i := len(history) - 1; i >= 0; i-- {
		image := history[i]
		re := regexp.MustCompile("\\s+")
		createdBy := re.ReplaceAllLiteralString(image.CreatedBy, " ")
		re = regexp.MustCompile("^/bin/sh -c #\\(nop\\) ")
		createdBy = re.ReplaceAllLiteralString(createdBy, "")
		re = regexp.MustCompile("^/bin/sh -c")
		createdBy = re.ReplaceAllLiteralString(createdBy, "RUN")
		tags := strings.Join(image.Tags, ", ")
		if !boolAll {
			createdBy = FormatNonBreakingString(Truncate(createdBy, 50))
			tags = FormatNonBreakingString(tags)
		}
		out := []string{
			Truncate(image.Id, 12),
			createdBy,
			tags,
			FormatDateTime(time.Unix(image.Created, 0)),
			FormatFloat(float64(image.Size) / 1000000),
		}
		items = append(items, out)
	}

	header := []string{
		"ID",
		"Created by",
		"Name:Tags",
		"Created at",
		"Size(MB)",
	}

	table := tablewriter.NewWriter(os.Stdout)
	table.SetColWidth(20)
	if !boolNoHeader {
		table.SetHeader(header)
	} else {
		table.SetBorder(false)
	}
	table.AppendBulk(items)
	table.Render()
}

func removeImages(ctx *cobra.Command, args []string) {
	if len(args) < 1 {
		fmt.Println("Needs an argument <NAME> at least to remove")
		ctx.Usage()
		return
	}

	docker, err := client.NewDockerClient(configPath, hostName)
	if err != nil {
		log.Fatal(err)
	}

	var lastError error
	for _, name := range args {
		err = docker.RemoveImage(name, boolForce, boolNoPrune)
		if err != nil {
			lastError = err
		}
	}
	if lastError != nil {
		log.Fatal(lastError)
	}
}
