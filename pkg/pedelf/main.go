package pedelf

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"strconv"

	"github.com/disintegration/imaging"
	"github.com/pdfcpu/pdfcpu/pkg/api"
	"github.com/pdfcpu/pdfcpu/pkg/pdfcpu"
	"github.com/pdfcpu/pdfcpu/pkg/pdfcpu/model"
)

type Context struct {
	*model.Context
	TotalPages int
	Src        io.ReadSeeker
	Images     []map[int]model.Image
}

func GetCtxFromInput(src io.ReadSeeker) (*Context, error) {
	// ctx, err := api.ReadContext(src, nil)
	ctx, err := api.ReadAndValidate(src, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to ReadContext, err: %s", err.Error())
	}

	if err := pdfcpu.OptimizeXRefTable(ctx); err != nil {
		return nil, fmt.Errorf("failed to OptimizeXRefTable, err: %s", err.Error())
	}
	if err := api.OptimizeContext(ctx); err != nil {
		return nil, fmt.Errorf("failed to OptimizeContext, err: %s", err.Error())
	}
	if err := ctx.EnsurePageCount(); err != nil {
		return nil, fmt.Errorf("failed to EnsurePageCount, err: %s", err.Error())
	}

	pages := ctx.PageCount
	PedelfCtx := &Context{
		Context:    ctx,
		TotalPages: pages,
		Src:        src,
	}

	return PedelfCtx, nil
}

func (ctx *Context) SetRawImagesFromGivenPages() error {
	log.Printf("total pages: %+v\n", ctx.TotalPages)

	pagesStr := []string{}
	for page := range ctx.TotalPages {
		pagesStr = append(pagesStr, strconv.Itoa(page))
	}

	rawImages, err := api.ExtractImagesRaw(ctx.Src, pagesStr, nil)
	if err != nil {
		return fmt.Errorf("failed to ExtractImagesRaw, err: %+v", err.Error())
	}

	ctx.Images = rawImages

	return nil
}

func (ctx *Context) CompressImages(compressLvl int) error {
	if compressLvl <= 1 {
		return fmt.Errorf("compressLvl must be larger than 1")
	}
	jpegQuality := 70

	for _, rawImage := range ctx.Images {
		for _, img := range rawImage {
			// log.Printf("img: name-%s, type-%s, page-%d, %dx%d\n", img.Name, img.FileType, img.PageNr, img.Width, img.Height)
			im, err := imaging.Decode(img)
			if err != nil {
				// unsupported image format may happen here, skip handling it
				fmt.Printf("image.Decode err: %+v\n", err.Error())
				continue
			}

			// resizedImg := imaging.Fit(im, im.Bounds().Dx()/compressLvl, im.Bounds().Dy()/compressLvl, imaging.Lanczos)
			resizedImg := imaging.Resize(im, im.Bounds().Dx()/compressLvl, 0, imaging.Lanczos)

			smallerBuf := new(bytes.Buffer)
			err = imaging.Encode(smallerBuf, resizedImg, imaging.JPEG, imaging.JPEGQuality(jpegQuality))
			if err != nil {
				return fmt.Errorf("jpeg.Encode err: %+v", err.Error())
			}

			// Update the original image with the resized image
			sd2, _, _, _ := model.CreateImageStreamDict(ctx.XRefTable, smallerBuf, false, false)
			ctx.XRefTable.Table[img.ObjNr].Object = *sd2
		}
	}

	return nil
}

func (ctx *Context) GetPdfBuffer() (*bytes.Buffer, error) {
	ctx.EnsureVersionForWriting()
	wr := new(bytes.Buffer)
	if err := api.WriteContext(ctx.Context, wr); err != nil {
		return nil, fmt.Errorf("failed to WriteContext, err: %s", err.Error())
	}

	return wr, nil
}
