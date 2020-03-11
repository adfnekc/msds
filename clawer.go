package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"golang.org/x/text/encoding/simplifiedchinese"
	"golang.org/x/text/transform"
)

func getCAS() []string {
	CasList := []string{}
	for i := 1; i < 1000; i++ {
		CASURL := fmt.Sprintf("http://www.cnreagent.com/msds/list_%d.html", i)
		casList, done := getCasList(CASURL)
		if done {
			break
		}
		CasList = append(CasList, casList...)
	}
	return CasList
}

func getCasList(url string) ([]string, bool) {
	doc := getDOM(url)

	casList := []string{}

	(doc.Find("tbody > tr[align!=center] > td[align=center]")).Each(
		func(index int, doc *goquery.Selection) {
			text := doc.Find("a").Text()
			casList = append(casList, text)
		},
	)
	if len(casList) == 0 {
		return nil, true
	}
	return casList, false
}

func handleErr(err error, errMsg string) {
	if err != nil {
		log.Fatal(errMsg)
	}
}

type msdsDATA struct {
	json        string
	instruction string
}

func getMsdsByCas(cas string) string {
	msdsURL := fmt.Sprintf("http://www.cnreagent.com/msds/cas_%s.html", cas)

	doc := getDOM(msdsURL)
	dic := map[string]string{}
	(doc.Find("table[bgcolor=\"#CCDDEE\"] > tbody > tr")).Each(
		func(index int, doc *goquery.Selection) {
			key := strings.ReplaceAll(doc.Find("td.msdsbt").Text(), "：", "")
			html, err := doc.Find("td.msdsnr").Html()
			handleErr(err, "doc.Find td.msdsnr err")
			value := strings.ReplaceAll(html, "\u003cbr/\u003e", "\n")
			dic[key] = value
		},
	)
	json, err := mapToJSON(dic)
	handleErr(err, fmt.Sprint("map to string err &v", err))
	return json
}

func getDOM(url string) *goquery.Document {
	resp, err := http.Get(url)
	handleErr(err, "http get error url:"+url)
	defer resp.Body.Close()

	reader := transform.NewReader(resp.Body, simplifiedchinese.GB18030.NewDecoder())

	doc, err := goquery.NewDocumentFromReader(reader)
	handleErr(err, "parseHtml err")
	return doc
}

func mapToJSON(m map[string]string) (string, error) {
	jsonByte, err := json.Marshal(m)
	if err != nil {
		fmt.Printf("Marshal with error: %+v\n", err)
		return "", err
	}

	return string(jsonByte), nil
}
