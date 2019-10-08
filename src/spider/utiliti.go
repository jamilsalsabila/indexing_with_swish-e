package main

import (
	"bytes"
	"math/rand"

	"golang.org/x/net/html"
)

type Date struct {
	Year  string
	Month string
	Day   int
}

var randomGenerator *rand.Rand

func InnerText(node *html.Node) string {
	var output func(buf *bytes.Buffer, n *html.Node)
	var buf bytes.Buffer

	output = func(buf *bytes.Buffer, n *html.Node) {
		for child := n.FirstChild; child != nil; child = child.NextSibling {
			// fmt.Println(child.Data, child.Type)
			if child.Data == "div" || child.Data == "br" || child.Data == "table" || child.Data == "script" || child.Data == "ins" {
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

var days [31]string = [31]string{
	"01", "02", "03", "04", "05",
	"06", "07", "08", "09", "10",
	"11", "12", "13", "14", "15",
	"16", "17", "18", "19", "20",
	"21", "22", "23", "24", "25",
	"26", "27", "28", "29", "30",
	"31"}
var months map[string]int = map[string]int{
	"01": 31,
	"02": 28,
	"03": 31,
	"04": 30,
	"05": 31,
	"06": 30,
	"07": 31,
	"08": 31,
	"09": 30,
	"10": 31,
	"11": 30,
	"12": 31,
}

var months2 map[int]string = map[int]string{
	1:  "01",
	2:  "02",
	3:  "03",
	4:  "04",
	5:  "05",
	6:  "06",
	7:  "07",
	8:  "08",
	9:  "09",
	10: "10",
	11: "11",
	12: "12",
}

func isLeapYear(year int) bool {
	if year%4 != 0 {
		return false
	} else if year%100 != 0 {
		return true
	} else if year%400 != 0 {
		return false
	} else {
		return true
	}
}

func TotalDaysOfMonth(year int, month string) int {
	if month == "02" {
		if isLeapYear(year) {
			return months[month] + 1
		}
	}
	return months[month]
}
