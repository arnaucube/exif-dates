# exif-dates
Simple script to help a friend to set the image date to the file date for all the images of a directory.

#### Usage
- Get the binary from https://github.com/arnaucube/exif-dates/releases
- Execute the binary
```
> ./exif-dates
```

Alternatively can specify the input directory: `> ./exif-dates inputDir`

This will generate a new directory `output` which contains all the images with the image date as the file date (both file's access and modification times), and also with the image date as the filename.
