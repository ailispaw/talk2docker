package commands

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"

	"github.com/ailispaw/talk2docker/client"
)

var cmdDocker = &cobra.Command{
	Use:     "docker [OPTIONS] COMMAND [arg...]",
	Aliases: []string{"dock"},
	Short:   "Execute the original docker cli",
	Long:    APP_NAME + " docker - Execute the original docker cli",
	Run:     docker,
}

func init() {
	cmdDocker.SetUsageFunc(func(ctx *cobra.Command) error {
		docker(ctx, []string{"--help"})
		return nil
	})
}

func docker(ctx *cobra.Command, args []string) {
	if _, err := exec.LookPath("docker"); err != nil {
		log.Fatal(err)
	}

	cmd := exec.Command("docker", args...)
	cmd.Env = getEnv()
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		log.Fatal(err)
	}
}

func getEnv() []string {
	config, err := client.LoadConfig(configPath)
	if err != nil {
		log.Fatal(err)
	}

	host, err := config.GetHost(hostName)
	if err != nil {
		log.Fatal(err)
	}

	var env []string
	for _, value := range os.Environ() {
		switch {
		case strings.HasPrefix(value, "DOCKER_HOST="):
		case strings.HasPrefix(value, "DOCKER_CERT_PATH="):
		case strings.HasPrefix(value, "DOCKER_TLS_VERIFY="):
		default:
			env = append(env, value)
		}
	}

	env = append(env, "DOCKER_HOST="+host.URL)
	if host.TLS {
		env = append(env, "DOCKER_CERT_PATH="+filepath.Dir(host.TLSCert))
		env = append(env, "DOCKER_TLS_VERIFY="+FormatBool(host.TLSVerify, "true", "false"))
	}

	return env
}
