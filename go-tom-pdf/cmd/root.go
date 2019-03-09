package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "tom-pdf",
	Short: "Converts HTML to PDF using one of several converters",
}

func init() {
	rootCmd.PersistentFlags().StringP("out", "o", "", "The PDF will be written into this file.")
	rootCmd.MarkPersistentFlagRequired("out")
	rootCmd.MarkPersistentFlagFilename("out", "pdf")

	rootCmd.AddCommand(newAPI2PDFCommand())
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Printf("Error: %v", err)
		os.Exit(1)
	}
}
