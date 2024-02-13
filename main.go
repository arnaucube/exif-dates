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

	"github.com/evanoberholster/imagemeta"
)

const version = "v0.0.3_2024-02-13"
const layout = "2006-01-02T15:04:05.000Z"

func main() {
	fmt.Printf("version: %s\n  (get latest version at https://github.com/arnaucube/exif-dates/releases )\n\n", version)

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

	nImgsDetected := 0
	nImgsConverted := 0
	var unconvertedFileNames []string
	for _, e := range entries {
		// skip subdirectories
		if e.IsDir() {
			continue
		}
		nImgsDetected++

		fileName := inputDirFlag + "/" + e.Name()
		fileExtension := path.Ext(fileName)
		// get file
		fmt.Printf("--> FileName: %s,", fileName)
		f, err := os.Open(fileName)
		if err != nil {
			fmt.Println("\ne", err)
		}
		defer f.Close()

		// get camera date
		date, err := getExifDate(f)
		if err != nil {
			fmt.Printf(" Error: %s\n\n", err)
			unconvertedFileNames = append(unconvertedFileNames, fileName)
			continue
		}
		fmt.Println("\n  DATE", date)
		dateString := date.String()
		if date.IsZero() {
			// if date is not set (=0), do not use the zero date as
			// name, and reuse the original name of the file
			fn := e.Name()
			dateString = fn[0 : len(fn)-len(fileExtension)]
		}

		newFileName := getValidFileName(outputDirFlag+"/"+dateString, fileExtension)
		newFileName = newFileName + fileExtension

		fmt.Printf("  storing the new img as %s\n\n", newFileName)

		// duplicate the original file into newFileName
		f, err = os.Open(fileName)
		if err != nil {
			fmt.Printf("\nerror (os.Open): %s\n\n", err)
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
		// Notice that when exif.time=0, this will be the date being
		// set here. In a future version might change it to just reuse
		// the original file date.
		os.Chtimes(newFileName, *date, *date)
		nImgsConverted++
	}
	fmt.Printf("converted %d images out of %d\n", nImgsConverted, nImgsDetected)
	fmt.Println("Unconverted images:")
	for i := 0; i < len(unconvertedFileNames); i++ {
		fmt.Println("    ", unconvertedFileNames[i])
	}
}

func getExifDate(f *os.File) (*time.Time, error) {
	x, err := imagemeta.Decode(f)
	if err != nil {
		return nil, err
	}
	t := x.DateTimeOriginal()
	return &t, nil
}

func getValidFileName(fileName, fileExtension string) string {
	_, err := os.Open(fileName + fileExtension)
	if !errors.Is(err, os.ErrNotExist) {
		return getValidFileName(fileName+"_1", fileExtension)
	}
	return fileName
}
