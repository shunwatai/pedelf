# Pedelf

Simple program for trying to "compress" the size of PDF by reducing the images in the file.

Used packages:

-   [pdfcpu](https://github.com/pdfcpu/pdfcpu/tree/master) for extract & update the images
-   [imaging](https://github.com/disintegration/imaging) for reduce the image size

Inpired by [pdfcomprezzor](https://github.com/henrixapp/pdfcomprezzor/blob/master/main.go)

## How it works

1. Extract images from pdf
2. Scale down the images and convert the images to jpg with lower quality
3. Update the images back in pdf
4. Export the results

## Example

```
cd examples/compress_pdfs_in_dir
go run main.go
```

## Contributing

- Please create issue if encounter issues
- Feature requests
