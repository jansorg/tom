package cmd

import (
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"

	"github.com/jansorg/gotime/go-tom-pdf/converter/api2pdf"
)

func newAPI2PDFCommand() *cobra.Command {
	var apiKey string
	var engine string

	cmd := &cobra.Command{
		Use:   "api2pdf",
		Short: "Converts HTML to PDF using api2pdf.com",
		Args:  cobra.MaximumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			outFile, err := cmd.Flags().GetString("out")
			if err != nil {
				log.Fatal(err)
			}

			var input io.Reader
			if len(args) == 1 {
				html, err := os.OpenFile(args[0], os.O_RDONLY, 0600)
				defer html.Close()
				if err != nil {
					log.Fatal(err)
				}

				input = html
			} else {
				input = os.Stdin
			}

			file, err := os.OpenFile(outFile, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0600)
			defer file.Close()
			if err != nil {
				log.Fatal(err)
			}

			requestType := api2pdf.RequestChrome
			switch engine {
			case "chrome":
				requestType = api2pdf.RequestChrome
			case "wkhtml":
				requestType = api2pdf.RequestWkHtml
			default:
				log.Fatal("unknown engine type", engine)
			}

			conv := api2pdf.NewConverter(apiKey, requestType, api2pdf.Request{
				FileName:  filepath.Base(outFile),
				InlinePDF: true,
			})
			if err = conv.ConvertHTML(input, file); err != nil {
				_ = os.Remove(outFile)
				log.Fatalf("error converting HTML to PDF: %v", err.Error())
			}

			fmt.Println("successfully converted HTML to PDF")
		},
	}

	cmd.Flags().StringVarP(&apiKey, "key", "", "", "The API key to send requests to api2pdf.com")
	cmd.Flags().StringVarP(&engine, "engine", "e", "chrome", "The HTML to PDF engine. Values: chrome | wkhtml. Default: chrome")
	cmd.MarkFlagRequired("key")

	return cmd
}
