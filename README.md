# exif-dates
Simple script to help a friend to set the image date to the file date for all the images of a directory.

#### Usage
- Get the binary from https://github.com/arnaucube/exif-dates/releases
- Execute the binary `> ./exif-dates`

```
> ./exif-dates -h
version: v0.0.2_2024-01-26
Usage of ./exif-dates:
  -input, -i string
        input directory (default "./")
  -output, -o string
        output directory (default "output")
  -version, -v
        print current version
```

This will generate a new directory `output` which contains all the images with the image date as the file date (both file's access and modification times), and also with the image date as the filename.
