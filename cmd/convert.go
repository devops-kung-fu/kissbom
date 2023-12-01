package cmd

import (
	"fmt"
	"log"
	"os"

	"github.com/devops-kung-fu/common/util"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"

	"github.com/devops-kung-fu/kissbom/lib"
)

var (
	outputFormats = []string{"json", "yaml", "csv", "minimal", "compatible"}

	selectedFormat string
	convertCmd     = &cobra.Command{
		Use:   "convert",
		Short: "Converts a provided CycloneDX file to a KISSBOM format",
		PreRun: func(cmd *cobra.Command, args []string) {
			if len(args) < 1 {
				util.PrintErr(errors.New("Please specify a file to convert"))
				os.Exit(1)
			}
		},
		Run: func(cmd *cobra.Command, args []string) {
			converter := lib.NewConverter()
			converter.OutputFormat = selectedFormat

			log.Println("starting conversion")
			err := converter.Convert(args[0])
			if err != nil {
				util.PrintErr(err)
				os.Exit(1)
			}

			log.Println("finished")
			util.PrintInfof("Saved KISSBOM as: %v\n", converter.OutputFileName)
			util.PrintSuccess("DONE!")
			os.Exit(0)
		},
	}
)

func init() {
	rootCmd.AddCommand(convertCmd)
	convertCmd.Flags().StringVarP(&selectedFormat, "format", "f", "json", fmt.Sprintf("select one of the valid options: %s", outputFormats))
	_ = rootCmd.Flags().SetAnnotation("format", cobra.BashCompOneRequiredFlag, []string{"true"})

}
