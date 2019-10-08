package main

import (
	"fmt"
	"io/ioutil"
	"os"
	//"github.com/RadhiFadlillah/go-sastrawi"
)

type Job struct {
	inputCurrFile  string
	OutputCurrFile string
	END            bool
}

func main() {
	var dirList []string
	var files []os.FileInfo
	var inputDir, outputDir string
	var e error
	var numOfWOrkers int
	var jobChannel chan Job

	dirList = []string{"kompas/", "tempo/", "antara/"}
	inputDir = "../../../ALL/"
	outputDir = "../../../ALL_2/"
	numOfWOrkers = 10
	jobChannel = make(chan Job, 1)

	/* Initiate Workers */
	for i := 0; i < numOfWOrkers; i++ {
		go TODO(i, jobChannel)
	}

	/* 1. Read files in dir. ALL Concurrently using 10 workers */
	for _, dir := range dirList {
		files, e = ioutil.ReadDir(inputDir + dir)
		if e != nil {
			panic(e)
		}

		/* 2. Per File, read the content */
		for _, file := range files {
			fmt.Println(inputDir + dir + file.Name())
			jobChannel <- Job{
				inputCurrFile:  inputDir + dir + file.Name(),
				OutputCurrFile: outputDir + dir + file.Name(),
				END:            false,
			}

		}

	}

	/* Tell workers there is no job */
	for i := 0; i < numOfWOrkers; i++ {
		jobChannel <- Job{
			END: true,
		}
	}

	close(jobChannel)
}
