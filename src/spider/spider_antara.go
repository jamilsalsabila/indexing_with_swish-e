package main

import (
	"bytes"
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

type WorkerAntara struct {
	Id int8
}

func SpiderAntara(workerId int8, site string, day string, month string, year string, fileOutput *os.File, dir string) {
	var resp *http.Response
	var e error
	var root *html.Node
	var targetNodes []*html.Node
	var pattern *regexp.Regexp
	var targetXMLPath string
	var numOfNewsString string
	var numOfNewsInt int

	for {
		resp, e = http.Get(fmt.Sprintf("https://www.antaranews.com/indeks/%s/%s-%s-%s", site, day, month, year))
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

	pattern = regexp.MustCompile(`\((\d*) berita\)`)
	numOfNewsString = pattern.FindStringSubmatch(htmlquery.FindOne(root, `//div[@class="col-sm-8"]/div/div/h1/span`).FirstChild.Data)[1]
	numOfNewsInt, e = strconv.Atoi(numOfNewsString)
	if e != nil {
		panic(e)
	}

	targetXMLPath = `//div[@class="col-sm-8"]/div[2]/div/article`
	targetNodes = htmlquery.Find(root, targetXMLPath)
	Download(workerId, targetNodes, fileOutput, dir, day, month, year, pattern, 0)

	for j := 1; j <= int(numOfNewsInt/10); j++ {
		for {
			resp, e = http.Get(fmt.Sprintf("https://www.antaranews.com/indeks/%s/%s-%s-%s/%d", site, day, month, year, j+1))
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
		targetNodes = htmlquery.Find(root, targetXMLPath)
		Download(workerId, targetNodes, fileOutput, dir, day, month, year, pattern, j)
		resp.Body.Close()
	}
}

func Download(workerId int8, targetNodes []*html.Node, fileOutput *os.File, dir string, day string, month string, year string, pattern *regexp.Regexp, pag int) {
	var e error
	var stringTemp, linkString, tagString string
	var header, linkNode, dateNode, contentRoot *html.Node
	var content *http.Response
	var contentTag []*html.Node
	var randomNumber int
	var duration time.Duration

	pattern = regexp.MustCompile(`\s+`)
	for i := 0; i < len(targetNodes); i++ {
		header = htmlquery.FindOne(targetNodes[i], `./header`)

		stringTemp = ""
		fileOutput, e = os.Create(dir + fmt.Sprintf("%s%s%s_%d_%d", day, month, year, pag, i))
		if e != nil {
			panic(e)
		}
		defer fileOutput.Close()
		linkNode = htmlquery.FindOne(header, `./h3/a`)
		linkString = htmlquery.SelectAttr(linkNode, "href")
		dateNode = htmlquery.FindOne(header, `./p/span`)
		tagString = htmlquery.FindOne(header, `./p/a`).FirstChild.Data

		fmt.Fprintf(fileOutput, "link: %s\ntitle: %s\ntag: %s\ndate: %s\n", linkString, strings.TrimSpace(InnerTextKhususAntara(linkNode)), tagString, strings.TrimSpace(InnerTextKhususAntara(dateNode)))

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

		contentTag = htmlquery.Find(contentRoot, `//div[@class="post-content clearfix"]`)
		for i := 0; i < len(contentTag); i++ {
			stringTemp += pattern.ReplaceAllString(InnerTextKhususAntara(contentTag[i]), " ")
		}

		fmt.Fprintf(fileOutput, "content: %s", stringTemp)

		content.Body.Close()
		fileOutput.Close()

		randomNumber = randomGenerator.Intn(40) + 30
		duration, e = time.ParseDuration(strconv.Itoa(randomNumber) + "s")
		if e != nil {
			panic(e)
		}
		fmt.Println(workerId, "sleep for", duration)
		time.Sleep(duration)
	}
}

func InnerTextKhususAntara(node *html.Node) string {
	var output func(buf *bytes.Buffer, n *html.Node)
	var buf bytes.Buffer

	output = func(buf *bytes.Buffer, n *html.Node) {
		for child := n.FirstChild; child != nil; child = child.NextSibling {
			// fmt.Println(child.Data, child.Type)
			if child.Data == "br" || child.Data == "table" || child.Data == "script" || child.Data == "ins" {
				buf.WriteString(" ")
				continue
			} else if child.Type == html.TextNode {
				buf.WriteString(child.Data + " ")
			} else {
				output(buf, child)
			}
		}
	}

	output(&buf, node)
	return buf.String()
}

func workerAntara(id int8, job chan Date) {
	var worker WorkerAntara
	var todo Date
	var fileOutput *os.File
	var e error
	var randomNumber int
	var duration time.Duration
	var workerPath string

	worker = WorkerAntara{Id: id}
	workerPath = "antara/sport/" + strconv.Itoa(int(worker.Id))

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
			SpiderAntara(worker.Id, "sport", days[i], todo.Month, todo.Year, fileOutput, workerPath+"/")
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
