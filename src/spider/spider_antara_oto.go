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

type WorkerAntaraOto struct {
	Id int8
}

func SpiderAntaraOto(workerId int8, todo int, fileOutput *os.File, dir string) {
	var resp *http.Response
	var e error
	var root *html.Node
	var targetNodes []*html.Node
	var targetXMLPath string

	for {
		resp, e = http.Get(fmt.Sprintf("https://otomotif.antaranews.com/umum/%d", todo))
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

	targetXMLPath = `//div[@class="container"]/div/div/div[@class="row"]/div[@class="col-xl-12 col-lg-6 col-md-6 col-sm-12"]`
	targetNodes = htmlquery.Find(root, targetXMLPath)

	DownloadOto(workerId, targetNodes, fileOutput, dir, todo)
}

func DownloadOto(workerId int8, targetNodes []*html.Node, fileOutput *os.File, dir string, todo int) {
	var pattern *regexp.Regexp
	var e error
	var stringTemp, linkString string
	var div, linkNode, dateNode, contentRoot *html.Node
	var content *http.Response
	var contentTag []*html.Node

	pattern = regexp.MustCompile(`\s+`)
	for i := 0; i < len(targetNodes); i++ {
		div = htmlquery.FindOne(targetNodes[i], `./div/div[2]`)

		stringTemp = ""
		fileOutput, e = os.Create(dir + fmt.Sprintf("%d_%d", todo, i))
		if e != nil {
			panic(e)
		}
		defer fileOutput.Close()
		
		/* perlu try-catch */
		linkNode = htmlquery.FindOne(div, `./h3/a`)
		linkString = htmlquery.SelectAttr(linkNode, "href")
		dateNode = htmlquery.FindOne(div, `./div`)

		fmt.Fprintf(fileOutput, "link: %s\ntitle: %s\ndate: %s\n", linkString, strings.TrimSpace(InnerText(linkNode)), strings.TrimSpace(InnerText(dateNode)))
		/* ---------------- */
		
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
			stringTemp += pattern.ReplaceAllString(InnerText(contentTag[i]), " ")
		}

		fmt.Fprintf(fileOutput, "content: %s", stringTemp)

		content.Body.Close()
		fileOutput.Close()
	}
}

func WorkerAntaraOtomotif(id int8, job chan int) {
	var worker WorkerAntaraOto
	var todo int
	var fileOutput *os.File
	var e error
	var randomNumber int
	var duration time.Duration
	var pathWorker string

	worker = WorkerAntaraOto{Id: id}
	// randomGenerator = rand.New(rand.NewSource(time.Now().UnixNano()))
	pathWorker = strconv.Itoa(int(worker.Id))
	e = os.Mkdir("antara/oto/"+pathWorker, 0755)
	if e != nil {
		panic(e)
	}

	for {
		todo = <-job
		fmt.Println(worker.Id, todo)
		if todo == -1 {
			break
		}
		SpiderAntaraOto(id, todo, fileOutput, "antara/oto/"+pathWorker+"/")

		randomNumber = randomGenerator.Intn(60) + 10
		duration, e = time.ParseDuration(strconv.Itoa(randomNumber) + "s")
		if e != nil {
			panic(e)
		}
		fmt.Println(worker.Id, "sleep for", duration)
		time.Sleep(duration)
	}

}
