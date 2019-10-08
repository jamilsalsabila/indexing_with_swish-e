package main

import (
	"fmt"
	"net/http"
	"os"
	"regexp"
	"strconv"
	"time"

	"github.com/antchfx/htmlquery"
	"golang.org/x/net/html"
)

type WorkerTempo struct {
	Id int8
}

func SpiderTempo(workerId int8, site string, day string, month string, year string, fileOutput *os.File, dir string) {
	var resp, content *http.Response
	var randomNumber int
	var duration time.Duration
	var e error
	var root, linkNode, dateNode, contentRoot, div *html.Node
	var targetNodes, contentTag []*html.Node
	var linkString string
	var stringTemp string
	var pattern *regexp.Regexp
	var targetXMLPath, titleString string

	pattern = regexp.MustCompile(`\s+`)

	for {
		resp, e = http.Get(fmt.Sprintf("https://www.tempo.co/indeks/%s/%s/%s/%s", year, month, day, site))
		if e == nil {
			break
		}
		fmt.Printf("%d: %s\n", workerId, "retry...")
		time.Sleep(30 * time.Second)
	}
	defer resp.Body.Close()

	root, e = htmlquery.Parse(resp.Body)
	if e != nil {
		panic(e)
	}

	targetXMLPath = `//section[@class="list list-type-1"][1]/ul/li`
	targetNodes = htmlquery.Find(root, targetXMLPath)

	for i := 0; i < len(targetNodes); i++ {
		div = htmlquery.FindOne(targetNodes[i], `./div/div`)

		stringTemp = ""
		fileOutput, e = os.Create(dir + fmt.Sprintf("%s%s%s_%s_%d", day, month, year, site, i))
		if e != nil {
			panic(e)
		}
		defer fileOutput.Close()
		linkNode = htmlquery.FindOne(div, `./a[2]`)
		linkString = htmlquery.SelectAttr(linkNode, "href")
		dateNode = htmlquery.FindOne(linkNode, `./span`)
		titleString = htmlquery.FindOne(linkNode, `./h2`).FirstChild.Data

		fmt.Fprintf(fileOutput, "link: %s\ntitle: %s\ndate: %s\n", linkString, titleString, dateNode.FirstChild.Data)

		// download content
		for {
			content, e = http.Get(linkString)
			if e == nil {
				break
			}
			fmt.Printf("%d: %s\n", workerId, "retry...")
			time.Sleep(30 * time.Second)
		}

		defer content.Body.Close()

		contentRoot, e = htmlquery.Parse(content.Body)
		if e != nil {
			panic(e)
		}

		contentTag = htmlquery.Find(contentRoot, `//div[@id="isi"]`)
		for i := 0; i < len(contentTag); i++ {
			stringTemp += pattern.ReplaceAllString(InnerText(contentTag[i]), " ")
		}

		fmt.Fprintf(fileOutput, "content: %s", stringTemp)

		content.Body.Close()
		fileOutput.Close()
		randomNumber = randomGenerator.Intn(30) + 30
		duration, e = time.ParseDuration(strconv.Itoa(randomNumber) + "s")
		if e != nil {
			panic(e)
		}
		fmt.Println(workerId, "sleep for", duration)
		time.Sleep(duration)
	}
}

func workerTempo(id int8, job chan Date) {
	var worker WorkerTempo
	var todo Date
	var fileOutput *os.File
	var e error
	var randomNumber int
	var duration time.Duration
	var workerPath string

	worker = WorkerTempo{Id: id}
	workerPath = "tempo/entertainment/" + strconv.Itoa(int(worker.Id))

	e = os.Mkdir(workerPath, 0755)
	if e != nil {
		panic(e)
	}

	for {
		todo = <-job
		if todo.Day == -1 {
			break
		}
		for i := 0; i < todo.Day; i++ {
			fmt.Printf("%d: %s/%s/%s\n", worker.Id, todo.Year, todo.Month, days[i])
			SpiderTempo(worker.Id, "seleb", days[i], todo.Month, todo.Year, fileOutput, workerPath+"/")
		}
		randomNumber = randomGenerator.Intn(60) + 10
		duration, e = time.ParseDuration(strconv.Itoa(randomNumber) + "s")
		if e != nil {
			panic(e)
		}
		fmt.Println(worker.Id, "sleep for", duration)
		time.Sleep(duration)

	}
}
