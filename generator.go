package main

import (
	"fmt"
	"os"
	"encoding/json"
	"io/ioutil"
	"io"
	"strings"
	"archive/zip"
)

type Rules struct {
	Directories []string `json:"directories"`
	ExcludedFiles []string `json:"excludedFiles"`
	Destination string `json:"destination"` 
}

func main() {
	fmt.Println("Hello, world.")
	buildData := getRules()
	excludedFiles := make(map[string]bool)
	for _, num := range buildData.ExcludedFiles {
    excludedFiles[num] = true
	}
	fmt.Println(excludedFiles)
	for i:=0; i<len(buildData.Directories); i++ {
		fmt.Println(buildData.Directories[i])
		filesToCompres:=listDirectoryFiles(buildData.Directories[i], excludedFiles)
		fmt.Println(filesToCompres)
		//For the MVP loop again the listed files and compress them
		//For the Version 2, compress them in parallel as they are listed
		fmt.Println("creating zip archive...")
		fmt.Println(strings.TrimLeft(buildData.Directories[i], "./")+".zip")
    archive, err := os.Create(strings.TrimLeft(buildData.Directories[i], "./")+".zip")
    if err != nil {
        panic(err)
    }
		defer archive.Close()
		zipWriter := zip.NewWriter(archive)
		//Write the files on the corresponding zip file
		for j:=0; j<len(filesToCompres); j++ {
			fmt.Println(filesToCompres[j])
			fmt.Println("Opening file: " + filesToCompres[j])
			f, err := os.Open(filesToCompres[j])
			if err != nil {
				fmt.Println(err)
			}
			defer f.Close()

			fmt.Println("Writing file in the zip archive")
			w, err := zipWriter.Create(filesToCompres[j])
			if err != nil {
				fmt.Println(err)
			}
			if _, err := io.Copy(w, f); err != nil {
				fmt.Println(err)
				panic(err)
			}
		}
		fmt.Println("Close Zip Archive...")
		zipWriter.Close()

	}
}

func getRules() Rules {
	jsonFile, err := os.Open("rules.json")
	// if we os.Open returns an error then handle it
	if err != nil {
			fmt.Println(err)
	}
	fmt.Println("Successfully Opened rules.json")
	// defer the closing of our jsonFile so that we can parse it later on
	defer jsonFile.Close()

	fileInBytes, _ := ioutil.ReadAll(jsonFile)
	var buildData Rules
	json.Unmarshal(fileInBytes, &buildData)
	return buildData
}

func listDirectoryFiles(directory string, excludedFiles map[string]bool) []string{
	files, err := ioutil.ReadDir(directory)
	if err != nil {
			fmt.Println(err)
	}
	var res []string

	for _, file := range files {
		if _, ok := excludedFiles[file.Name()]; ok {
			fmt.Println("Excluded file: " + file.Name())
		} else {
			if file.IsDir() {
				innerDir := listDirectoryFiles(directory + "/" + file.Name(), excludedFiles)
				res = append(res, innerDir...)
			} else {
				res = append(res, directory + "/" + file.Name())
			}
		}
		
	}
	return res
}