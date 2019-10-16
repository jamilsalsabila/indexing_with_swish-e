package main

import (
	"fmt"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/antchfx/htmlquery"
	"golang.org/x/net/html"
)

type WorkerKompas struct {
	Id int8
}

func SpiderKompas(site string, page int, fileOutput *os.File, dir string) {
	var resp, content *http.Response
	var e error
	var root, linkNode, tagNode, dateNode, contentRoot *html.Node
	var targetNodes, contentTag []*html.Node
	var linkString string
	var stringTemp string

	for {
		resp, e = http.Get("https://indeks.kompas.com/?site=" + site + "&page=" + strconv.Itoa(page))
		if e == nil {
			break
		}
		time.Sleep(30 * time.Second)
	}

	defer resp.Body.Close()

	root, e = htmlquery.Parse(resp.Body)
	if e != nil {
		panic(e)
	}

	targetNodes = htmlquery.Find(root, `//div[@class="article__list clearfix"]`)

	for i := 0; i < len(targetNodes); i++ {
		stringTemp = ""
		fileOutput, e = os.Create(dir + strconv.Itoa(page) + "_" + strconv.Itoa(i))
		if e != nil {
			panic(e)
		}
		defer fileOutput.Close()
		
		/* perlu try-catch */
		linkNode = htmlquery.FindOne(targetNodes[i], `./div[@class="article__list__title"]/h3/a`)
		linkString = htmlquery.SelectAttr(linkNode, "href")
		tagNode = htmlquery.FindOne(targetNodes[i], `./div[@class="article__list__info"]/div[@class="article__subtitle article__subtitle--inline"]`)
		dateNode = htmlquery.FindOne(targetNodes[i], `./div[@class="article__list__info"]/div[@class="article__date"]`)

		fmt.Fprintf(fileOutput, "link: %s\ntitle: %s\ntag: %s\ndate: %s\n", linkString, linkNode.FirstChild.Data, tagNode.FirstChild.Data, dateNode.FirstChild.Data)
		/* ---------------- */
		
		// download content
		for {
			content, e = http.Get(linkString)
			if e == nil {
				break
			}
			time.Sleep(30 * time.Second)
		}

		defer content.Body.Close()

		contentRoot, e = htmlquery.Parse(content.Body)
		if e != nil {
			panic(e)
		}

		contentTag = htmlquery.Find(contentRoot, `//div[@class="col-bs9-7"]/div[@class="read__content"]`)
		for i := 0; i < len(contentTag); i++ {
			stringTemp += strings.TrimSpace(InnerText(contentTag[i]))
		}

		fmt.Fprintf(fileOutput, "content: %s", stringTemp)

		content.Body.Close()
		fileOutput.Close()
	}

}

func workerKompas(id int8, job chan int) {
	var worker WorkerKompas
	var todo int
	var fileOutput *os.File
	var e error
	var randomNumber int
	var duration time.Duration
	var pathWorker string

	worker = WorkerKompas{Id: id}
	// randomGenerator = rand.New(rand.NewSource(time.Now().UnixNano()))
	pathWorker = strconv.Itoa(int(worker.Id))
	e = os.MkdirAll("kompas/tekno/"+pathWorker, 0755)
	if e != nil {
		panic(e)
	}

	for {
		todo = <-job
		fmt.Println(worker.Id, todo)
		if todo == -1 {
			break
		}
		SpiderKompas("tekno", todo, fileOutput, "kompas/tekno/"+strconv.Itoa(int(worker.Id))+"/")

		randomNumber = randomGenerator.Intn(60) + 10
		duration, e = time.ParseDuration(strconv.Itoa(randomNumber) + "s")
		if e != nil {
			panic(e)
		}
		fmt.Println(worker.Id, "sleep for", duration)
		time.Sleep(duration)
	}

}
