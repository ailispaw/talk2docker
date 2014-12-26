package commands

import (
	"fmt"
	"log"
	"regexp"
	"sort"
	"strings"
	"time"

	"github.com/spf13/cobra"
	"github.com/yungsang/tablewriter"
	"github.com/yungsang/talk2docker/api"
	"github.com/yungsang/talk2docker/client"
)

var (
	boolForce, boolNoPrune, boolStar bool
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
	Long:    APP_NAME + " image history - Show the histry of an image",
	Run:     showImageHistory,
}

var cmdInspectImage = &cobra.Command{
	Use:     "inspect <NAME[:TAG]|ID>",
	Aliases: []string{"ins", "info"},
	Short:   "Inspect an image",
	Long:    APP_NAME + " image inspect - Inspect an image",
	Run:     inspectImage,
}

var cmdPushImage = &cobra.Command{
	Use:   "push <NAME[:TAG]>",
	Short: "Push an image",
	Long:  APP_NAME + " image push - Push an image",
	Run:   pushImage,
}

var cmdRemoveImages = &cobra.Command{
	Use:     "remove <NAME[:TAG]|ID>...",
	Aliases: []string{"rm"},
	Short:   "Remove images",
	Long:    APP_NAME + " image remove - Remove images",
	Run:     removeImages,
}

var cmdSearchImages = &cobra.Command{
	Use:   "search <TERM>",
	Short: "Search images",
	Long:  APP_NAME + " image search - Search images",
	Run:   searchImages,
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

	cmdSearchImages.Flags().BoolVarP(&boolStar, "star", "s", false, "Sort by star")
	cmdSearchImages.Flags().BoolVarP(&boolQuiet, "quiet", "q", false, "Only display names")
	cmdSearchImages.Flags().BoolVarP(&boolNoHeader, "no-header", "n", false, "Omit the header")

	cmdImage.AddCommand(cmdListImages)
	cmdImage.AddCommand(cmdPullImage)
	cmdImage.AddCommand(cmdTagImage)
	cmdImage.AddCommand(cmdShowImageHistory)
	cmdImage.AddCommand(cmdInspectImage)
	cmdImage.AddCommand(cmdPushImage)
	cmdImage.AddCommand(cmdRemoveImages)
	cmdImage.AddCommand(cmdSearchImages)
}

func listImages(ctx *cobra.Command, args []string) {
	docker, err := client.NewDockerClient(configPath, hostName, ctx.Out())
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
				ctx.Println(Truncate(image.Id, 12))
			}
		}
		return
	}

	if boolYAML || boolJSON {
		items := []api.Image{}
		for _, image := range images {
			if (matchName == "") || matchImageByName(image.RepoTags, matchName) {
				items = append(items, image)
			}
		}
		if err := FormatPrint(ctx.Out(), items); err != nil {
			log.Fatal(err)
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

	PrintInTable(ctx.Out(), header, items, 0, tablewriter.ALIGN_DEFAULT)
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
		ctx.Println("Needs an argument <NAME[:TAG]> to pull")
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

	docker, err := client.NewDockerClient(configPath, hostName, ctx.Out())
	if err != nil {
		log.Fatal(err)
	}

	if err := docker.PullImage(repository); err != nil {
		log.Fatal(err)
	}
}

func tagImage(ctx *cobra.Command, args []string) {
	if len(args) < 2 {
		ctx.Println("Needs two arguments <IMAGE-NAME[:TAG] or IMAGE-ID> <NEW-NAME[:TAG]>")
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

	docker, err := client.NewDockerClient(configPath, hostName, ctx.Out())
	if err != nil {
		log.Fatal(err)
	}

	if err := docker.TagImage(args[0], name, tag, boolForce); err != nil {
		log.Fatal(err)
	}

	ctx.Printf("Tagged %s as %s:%s\n", args[0], name, tag)
}

func showImageHistory(ctx *cobra.Command, args []string) {
	if len(args) < 1 {
		ctx.Println("Needs an argument <IMAGE-NAME[:TAG] or IMAGE-ID>")
		ctx.Usage()
		return
	}

	docker, err := client.NewDockerClient(configPath, hostName, ctx.Out())
	if err != nil {
		log.Fatal(err)
	}

	history, err := docker.GetImageHistory(args[0])
	if err != nil {
		log.Fatal(err)
	}

	sort.Sort(history)

	if boolYAML || boolJSON {
		if err := FormatPrint(ctx.Out(), history); err != nil {
			log.Fatal(err)
		}
		return
	}

	var items [][]string

	for _, image := range history {
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

	PrintInTable(ctx.Out(), header, items, 20, tablewriter.ALIGN_DEFAULT)
}

func inspectImage(ctx *cobra.Command, args []string) {
	if len(args) < 1 {
		ctx.Println("Needs an argument <IMAGE-NAME[:TAG] or IMAGE-ID>")
		ctx.Usage()
		return
	}

	docker, err := client.NewDockerClient(configPath, hostName, ctx.Out())
	if err != nil {
		log.Fatal(err)
	}

	imageInfo, err := docker.InspectImage(args[0])
	if err != nil {
		log.Fatal(err)
	}

	if err := FormatPrint(ctx.Out(), imageInfo); err != nil {
		log.Fatal(err)
	}
}

func pushImage(ctx *cobra.Command, args []string) {
	if len(args) < 1 {
		ctx.Println("Needs an argument <NAME[:TAG]> to push")
		ctx.Usage()
		return
	}

	registry, name, tag, err := client.ParseRepositoryName(args[0])
	if err != nil {
		log.Fatal(err)
	}

	if len(strings.SplitN(name, "/", 2)) == 1 {
		log.Fatalf("You cannot push a \"root\" repository. Please rename your repository in <yourname>/%s", name)
	}

	if registry != "" {
		name = registry + "/" + name
	}

	config, err := client.LoadConfig(configPath)
	if err != nil {
		log.Fatal(err)
	}

	if registry == "" {
		registry = client.INDEX_SERVER
	}

	registryConfig, err := config.GetRegistry(registry)
	if (err != nil) || (registryConfig.Auth == "") {
		log.Fatal("Please login prior to pushing an image.")
	}

	docker, err := client.NewDockerClient(configPath, hostName, ctx.Out())
	if err != nil {
		log.Fatal(err)
	}

	if err := docker.PushImage(name, tag, registryConfig.Auth); err != nil {
		log.Fatal(err)
	}
}

func removeImages(ctx *cobra.Command, args []string) {
	if len(args) < 1 {
		ctx.Println("Needs an argument <NAME[:TAG]> at least to remove")
		ctx.Usage()
		return
	}

	docker, err := client.NewDockerClient(configPath, hostName, ctx.Out())
	if err != nil {
		log.Fatal(err)
	}

	var lastError error
	for _, name := range args {
		if err := docker.RemoveImage(name, boolForce, boolNoPrune); err != nil {
			lastError = err
		}
	}
	if lastError != nil {
		log.Fatal(lastError)
	}
}

func searchImages(ctx *cobra.Command, args []string) {
	if len(args) < 1 {
		ctx.Println("Needs an argument <TERM> to search")
		ctx.Usage()
		return
	}

	docker, err := client.NewDockerClient(configPath, hostName, ctx.Out())
	if err != nil {
		log.Fatal(err)
	}

	images, err := docker.SearchImages(args[0])
	if err != nil {
		log.Fatal(err)
	}

	if boolStar {
		sort.Sort(sort.Reverse(api.SortImagesByStars{images}))
	} else {
		sort.Sort(api.SortImagesByName{images})
	}

	if boolQuiet {
		for _, image := range images {
			ctx.Println(image.Name)
		}
		return
	}

	if boolYAML || boolJSON {
		if err := FormatPrint(ctx.Out(), images); err != nil {
			log.Fatal(err)
		}
		return
	}

	var items [][]string

	for _, image := range images {
		out := []string{
			image.Name,
			image.Description,
			FormatInt(int64(image.Stars)),
			FormatBool(image.Official, "*", " "),
			FormatBool(image.Automated, "*", " "),
		}
		items = append(items, out)
	}

	header := []string{
		"Name",
		"Description",
		"Stars",
		"Official",
		"Automated",
	}

	PrintInTable(ctx.Out(), header, items, 50, tablewriter.ALIGN_DEFAULT)
}
