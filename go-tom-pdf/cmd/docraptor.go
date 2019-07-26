package cmd

import (
	"fmt"
	"github.com/jansorg/tom/go-tom-pdf/converter/docraptor"
	"github.com/spf13/cobra"
	"io"
	"io/ioutil"
	"log"
	"os"
	"strings"
)

func newDocraptorCommand() *cobra.Command {
	var apiKey string
	var mediaType string
	var documentType string
	var testMode bool
	var name string

	cmd := &cobra.Command{
		Use:   "docraptor [input file]",
		Short: "Converts HTML to PDF using docraptor.com",
		Args:  cobra.MaximumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			outFile, err := cmd.Flags().GetString("out")
			if err != nil {
				log.Fatal(err)
			}

			media := docraptor.PRINT
			switch mediaType {
			case string(docraptor.PRINT):
				media = docraptor.PRINT
			case string(docraptor.SCREEN):
				media = docraptor.SCREEN
			default:
				log.Fatal("unknown engine type", mediaType)
			}

			docType := docraptor.PDF
			switch documentType {
			case string(docraptor.PDF):
				docType = docraptor.PDF
			case string(docraptor.XLSX):
				docType = docraptor.XLSX
			case string(docraptor.XLS):
				docType = docraptor.XLSX
			default:
				log.Fatalf("unknown document type %s", documentType)
			}

			request := docraptor.Request{
				Test:       testMode,
				Name:       name,
				Type:       docType,
				JavaScript: true,
				PrinceOptions: docraptor.PrinceOptions{
					Media: media,
				},
			}

			var input io.ReadCloser
			if len(args) == 1 && (strings.HasPrefix(args[0], "http://") || strings.HasPrefix(args[0], "https://")) {
				request.DocumentURL = args[0]
			} else if len(args) == 1 {
				html, err := os.OpenFile(args[0], os.O_RDONLY, 0600)
				if err != nil {
					log.Fatal(err)
				}
				input = html
			} else {
				input = os.Stdin
			}

			if input != nil {
				defer input.Close()
				if htmlBytes, err := ioutil.ReadAll(input); err != nil {
					log.Fatal(err)
				} else {
					request.DocumentContent = string(htmlBytes)
				}
			}

			file, err := os.OpenFile(outFile, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0600)
			if err != nil {
				log.Fatal(err)
			}
			defer file.Close()

			conv := docraptor.NewConverter(apiKey, request)
			if err = conv.ConvertHTML(input, file); err != nil {
				_ = os.Remove(outFile)
				log.Fatalf("error converting HTML to PDF: %v", err.Error())
			}

			fmt.Println("successfully converted HTML to PDF")
		},
	}

	cmd.Flags().StringVarP(&apiKey, "key", "", "", "The API key to send requests to api2pdf.com")
	cmd.Flags().BoolVarP(&testMode, "test", "d", false, "Generate a test document")
	cmd.Flags().StringVarP(&mediaType, "media", "m", string(docraptor.PRINT), "The media type. Values: print | screen. Default: print")
	cmd.Flags().StringVarP(&documentType, "type", "t", string(docraptor.PDF), "The document type. Values: pdf | xlsx | xls. Default: pdf")
	cmd.Flags().StringVarP(&name, "name", "", "", "The name of the document, used by docraptor.com")
	cmd.MarkFlagRequired("key")

	return cmd
}
