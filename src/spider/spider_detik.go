package main

import (
	"fmt"
	"net/http"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/antchfx/htmlquery"
	"golang.org/x/net/html"
)

type WorkerDetik struct {
	Id int8
}

func SpiderDetik(workerId int8, site string, day string, month string, year string, fileOutput *os.File, dir string) {
	var resp, content *http.Response
	var e error
	var root, linkNode, dateNode, contentRoot, div *html.Node
	var targetNodes, contentTag []*html.Node
	var linkString string
	var stringTemp string
	var pattern *regexp.Regexp
	var targetXMLPath string
	var randomNumber int
	var duration time.Duration

	pattern = regexp.MustCompile(`\s+`)
	resp, e = http.Get("https://" + site + ".detik.com/indeks?date=" + month + "%2F" + day + "%2F" + year)
	if e != nil {
		panic(e)
	}
	defer resp.Body.Close()

	root, e = htmlquery.Parse(resp.Body)
	if e != nil {
		panic(e)
	}

	if site == "hot" {
		targetXMLPath = `//div[@class="lf_content boxwhite mt10 w850"]/ul/li`
	} else if site == "sport" {
		targetXMLPath = `//div[@class="lf_content boxlr w868 fr ml10"]/ul/li`
	} else if site == "oto" {
		targetXMLPath = `//div[@class="lf_content fl w870"]/ul/li`
	}

	targetNodes = htmlquery.Find(root, targetXMLPath)

	for i := 0; i < len(targetNodes); i++ {
		div = htmlquery.FindOne(targetNodes[i], `./article/div`)

		if htmlquery.FindOne(div, `./span[@class="sub_judul"]`) != nil {
			continue
		}

		stringTemp = ""
		fileOutput, e = os.Create(dir + fmt.Sprintf("%s%s%s_%d", day, month, year, i))
		if e != nil {
			panic(e)
		}
		defer fileOutput.Close()
		linkNode = htmlquery.FindOne(div, `./a`)
		linkString = htmlquery.SelectAttr(linkNode, "href")
		dateNode = htmlquery.FindOne(div, `./span`)

		// fmt.Println(linkString, strings.TrimSpace(InnerText(linkNode)))

		fmt.Fprintf(fileOutput, "link: %s\ntitle: %s\ndate: %s\n", linkString, strings.TrimSpace(InnerText(linkNode)), dateNode.FirstChild.Data)

		// download content
		content, e = http.Get(linkString)
		if e != nil {
			panic(e)
		}
		defer content.Body.Close()

		contentRoot, e = htmlquery.Parse(content.Body)
		if e != nil {
			panic(e)
		}

		contentTag = htmlquery.Find(contentRoot, `//div[@class="itp_bodycontent detail_text"]`)
		for i := 0; i < len(contentTag); i++ {
			stringTemp += pattern.ReplaceAllString(InnerText(contentTag[i]), " ")
		}

		fmt.Fprintf(fileOutput, "content: %s", stringTemp)

		content.Body.Close()
		fileOutput.Close()
		randomNumber = randomGenerator.Intn(5) + 5
		duration, e = time.ParseDuration(strconv.Itoa(randomNumber) + "s")
		if e != nil {
			panic(e)
		}
		fmt.Println(workerId, "sleep for", duration)
		time.Sleep(duration)
	}
}

func workerDetik(id int8, job chan Date) {
	var worker WorkerDetik
	var todo Date
	var fileOutput *os.File
	var e error
	var randomNumber int
	var duration time.Duration
	var workerPath string

	worker = WorkerDetik{Id: id}
	workerPath = "detik/oto/" + strconv.Itoa(int(worker.Id))

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
			SpiderDetik(worker.Id, "oto", days[i], todo.Month, todo.Year, fileOutput, workerPath+"/")
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
