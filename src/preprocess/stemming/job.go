package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	// "github.com/RadhiFadlillah/go-sastrawi"

	"../../../go-sastrawi"
)

func TODO(id int, jobChannel chan Job) {
	var job Job
	var inputCurrFile, outputCurrFile *os.File
	var e error
	var b []byte
	var stemmer sastrawi.Stemmer
	var strTemp string

	stemmer = sastrawi.NewStemmer(sastrawi.DefaultDictionary)

	for {
		job = <-jobChannel
		if job.END {
			fmt.Printf("worker %d selesai\n", id)
			break
		}
		fmt.Println(id, "processing:", job.inputCurrFile)
		inputCurrFile, e = os.Open(job.inputCurrFile)
		if e != nil {
			panic(e)
		}
		outputCurrFile, e = os.Create(job.OutputCurrFile)
		if e != nil {
			panic(e)
		}
		b, e = ioutil.ReadAll(inputCurrFile)
		if e != nil {
			panic(e)
		}

		/* 3. Per Word, stem */
		for _, word := range strings.Split(string(b), " ") {
			strTemp += stemmer.Stem(word) + " "
		}

		/* 4. Save to new file in ALL_2 dir. */
		outputCurrFile.WriteString(strTemp)
		inputCurrFile.Close()
		outputCurrFile.Close()
		b = nil
		strTemp = ""
	}
}
