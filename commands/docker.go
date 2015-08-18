package commands

import (
	"os"
	"os/exec"
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

	args = append([]string{"--host", host.URL}, args...)
	if host.TLS {
		args = append([]string{"--tls", "true"}, args...)
		args = append([]string{"--tlscacert", host.TLSCaCert}, args...)
		args = append([]string{"--tlscert", host.TLSCert}, args...)
		args = append([]string{"--tlskey", host.TLSKey}, args...)
		args = append([]string{"--tlsverify", FormatBool(host.TLSVerify, "true", "false")}, args...)
	}

	cmd := exec.Command("docker", args...)
	cmd.Env = env
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		log.Fatal(err)
	}
}
