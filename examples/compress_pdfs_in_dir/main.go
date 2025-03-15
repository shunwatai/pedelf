package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/dustin/go-humanize"
	"github.com/shunwatai/pedelf/pkg/pedelf"
)

func main() {
	fmt.Printf("test pdf \n")
	inDir := "./sample-pdf"
	outDir := "./processed-pdf"

	fileNames := []string{}
	files, err := os.ReadDir(inDir)
	if err != nil {
		log.Fatal(err)
	}
	for _, file := range files {
		fileNames = append(fileNames, file.Name())
	}

	// Loop pdfs in ./sample-pdf
	for _, fileName := range fileNames {
		inFile := filepath.Join(inDir, fileName)
		outFile := filepath.Join(outDir, fileName)

		f, err := os.Open(inFile)
		if err != nil {
			log.Fatalf("failed to open file: %+v, err: %+v", inFile, err.Error())
		}

		// Open the pdf by pdfcpu and get its context
		ctx, err := pedelf.GetCtxFromInput(f)
		// ctx, err := api.ReadContext(f, nil)
		if err != nil {
			log.Fatalf("failed to ReadContext, err: %+v", err.Error())
		}
		log.Printf("processing: %s", inFile)

		// Get all images from pdf and set them in ctx
		ctx.SetRawImagesFromGivenPages()

		// Reduce images size
		if err := ctx.CompressImages(2); err != nil {
			log.Fatalf("failed to CompressImages, err: %+v", err.Error())
		}

		// Get the updated pdf buffer
		wr, err := ctx.GetPdfBuffer()
		if err != nil {
			log.Fatalf("failed to GetPdfBuffer, err: %+v", err.Error())
		}

		oriStat, _ := f.Stat()
		log.Printf("finished: %s", outFile)
		log.Printf("original size: %+v, resized size: %+v", humanize.Bytes(uint64(oriStat.Size())), humanize.Bytes(uint64(wr.Len())))
		fmt.Println()

		os.WriteFile(outFile, wr.Bytes(), 0644)
	}
}
