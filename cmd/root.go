// Package cmd contains all of the commands that may be executed in the cli
package cmd

import (
	"fmt"
	"io"
	"log"
	"os"
	"strings"

	"github.com/devops-kung-fu/common/github"
	"github.com/gookit/color"
	"github.com/spf13/cobra"
)

var (
	version = "0.0.1"
	//Verbose determines if the execution of hing should output verbose information
	debug   bool
	rootCmd = &cobra.Command{
		Use:     "kissbom [flags] file",
		Example: "  kissbom convert test.cyclonedx.json",
		Short:   "Converts a CycloneDX file to a KISSBOM.",
		Version: version,
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			if !debug {
				log.SetOutput(io.Discard)
			}

			log.Println("Start")
			fmt.Println()
			color.Style{color.FgWhite, color.OpBold}.Println("█▄▀ █ █▀ █▀ ██▄ █▀█ █▀▄▀█")
			color.Style{color.FgWhite, color.OpBold}.Println("█ █ █ ▄█ ▄█ █▄█ █▄█ █ ▀ █")
			fmt.Println()
			fmt.Println("DKFM - DevOps Kung Fu Mafia")
			fmt.Println("https://github.com/devops-kung-fu/kissbom")
			fmt.Printf("Version: %s\n", version)
			fmt.Println()
			latestVersion, _ := github.LatestReleaseTag("devops-kung-fu", "kissbom")
			if !strings.Contains(latestVersion, version) {
				color.Yellow.Printf("A newer version of kissbom is available (%s)\n\n", latestVersion)
			}

		},
	}
)

// Execute creates the command tree and handles any error condition returned
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().BoolVar(&debug, "debug", false, "displays debug level log messages.")
}
