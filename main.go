package main

import (
	"errors"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"path"
	"time"

	"github.com/rwcarlsen/goexif/exif"
)

const version = "v0.0.2_2024-01-26"
const layout = "2006-01-02T15:04:05.000Z"

func main() {
	fmt.Println("version:", version)

	var versionFlag bool
	var inputDirFlag, outputDirFlag string
	flag.BoolVar(&versionFlag, "version", false, "print current version")
	flag.BoolVar(&versionFlag, "v", false, "print current version")
	flag.StringVar(&inputDirFlag, "input", "./", "input directory")
	flag.StringVar(&inputDirFlag, "i", "./", "input directory")
	flag.StringVar(&outputDirFlag, "output", "output", "output directory")
	flag.StringVar(&outputDirFlag, "o", "output", "output directory")
	flag.Parse()

	if versionFlag {
		os.Exit(0)
	}

	if err := os.Mkdir(outputDirFlag, os.ModePerm); err != nil {
		log.Fatal("output directory already exists")
	}

	entries, err := os.ReadDir(inputDirFlag)
	if err != nil {
		log.Fatal(err)
	}

	for _, e := range entries {
		// skip subdirectories
		if e.IsDir() {
			continue
		}

		fileName := inputDirFlag + "/" + e.Name()
		fileExtension := path.Ext(fileName)
		// get file
		fmt.Printf("---FileName: %s,", fileName)
		f, err := os.Open(fileName)
		if err != nil {
			fmt.Println("\ne", err)
		}
		defer f.Close()

		// get camera date
		date, err := getExifDate(f)
		if err != nil && (err == io.EOF || err.Error() == "exif: failed to find exif intro marker") {
			fmt.Println("  not an img file, skipping file")
			continue
		}
		if err != nil {
			fmt.Println("err", err)
			continue
		}
		fmt.Println("\n  DATE", date)

		newFileName := getValidFileName(outputDirFlag+"/"+date.String(), fileExtension)
		newFileName = newFileName + fileExtension

		fmt.Println("  storing the new img as", newFileName)

		// duplicate the original file into newFileName
		f, err = os.Open(fileName)
		if err != nil {
			fmt.Println("\ne", err)
		}
		defer f.Close()
		fo, err := os.Create(newFileName)
		if err != nil {
			panic(err)
		}
		defer fo.Close()
		buf := make([]byte, 1024)
		for {
			n, err := f.Read(buf)
			if err != nil && err != io.EOF {
				panic(err)
			}
			if n == 0 {
				break
			}

			if _, err := fo.Write(buf[:n]); err != nil {
				panic(err)
			}
		}
		// set the img date into the file date (file's access and modification times)
		os.Chtimes(newFileName, *date, *date)

	}

}

func getExifDate(f *os.File) (*time.Time, error) {
	x, err := exif.Decode(f)
	if err != nil {
		return nil, err
	}

	d, err := x.DateTime()
	if err != nil {
		return nil, err
	}

	return &d, nil
}

func getValidFileName(fileName, fileExtension string) string {
	_, err := os.Open(fileName + fileExtension)
	if !errors.Is(err, os.ErrNotExist) {
		return getValidFileName(fileName+"_1", fileExtension)
	}
	return fileName
}
