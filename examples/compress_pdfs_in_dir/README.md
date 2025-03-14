# Compress PDFs in a folder

## Desc.

Example for processing the PDFs in a folder, try to compress and then update them in each PDF, then output to the destination folder.

-   src folder `./examples/compress_pdfs_in_dir/sample-pdf`

    You can put some PDFs in this src folder to see the results.

-   dst folder `./examples/compress_pdfs_in_dir/processed-pdf`

    After the example executed, you can compare the file's size of the original one in src folder.

## Try it

Feel free to adjust the param of the `CompressImages(2)`, larger number for scaling down the image more. 
Also can adjust the `jpegQuality` for jpeg quality.

```
go run main.go
```
