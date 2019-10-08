package main

import (
	"math/rand"
	"time"
)

func init() {
	randomGenerator = rand.New(rand.NewSource(time.Now().UnixNano()))
}

func main() {
	var numOfWorkers int = 10
	/*
	var startYear int
	var endYear int
	*/
	var start int
	var end int
	var jobsChannel chan int
	//var jobs2Channel chan Date

	jobsChannel = make(chan int, 1)
	//jobs2Channel = make(chan Date, 1)
	start = 666
	end = 200
	/*
	startYear = 2018
	endYear = 2016
	*/
	// membuat pekerja
	for i := 0; i < numOfWorkers; i++ {
		go workerKompas(int8(i), jobsChannel)
	}

	/* mendistribusikan jobs ke para pekerja */
	for i := start; i >= end; i-- {
		jobsChannel <- i
	}
	/*
	for y := startYear; y >= endYear; y-- {
		for m := 1; m <= 12; m++ {
			jobs2Channel <- Date{
				Year:  strconv.Itoa(y),
				Month: months2[m],
				Day:   TotalDaysOfMonth(y, months2[m]),
			}
		}
	}
	*/
	/* pemberitahuan bahwa jobs sudah habis */
	for i := 0; i < numOfWorkers; i++ {
		jobsChannel <- -1
	}
	/*
	for i := 0; i < numOfWorkers; i++ {
		jobs2Channel <- Date{Day: -1}
	}
	*/
	close(jobsChannel)

	/* TEST */
	// var file *os.File
	// SpiderDetik("oto", "04", "10", "2019", file, "./detik/oto/")
	// SpiderKompas("entertainment", 10, file, "kompas/entertainment/")
	// SpiderAntara("hiburan", "04", "10", "2019", file, "antara/entertainment/")
	// SpiderTempo("travel", "04", "10", "2019", file, "tempo/entertainment/")
	// SpiderAntaraOto(885, hotfile, "antara/oto/")
}
