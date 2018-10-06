package jarhandler

import (
	"net/http"
	"github.com/PuerkitoBio/goquery"
)

type JarHandler struct {
}

func (jarHandler *JarHandler) Handler(jarname string) string {

	return searchJar(jarname)
}

func searchJar(jarname string) string {
	search, err := http.Get("https://mvnrepository.com/search?q=" + jarname)
	if err != nil {
		return "search http request error"
	}

	if search.StatusCode != 200 {
		return "status code not 200"
	}

	doc, err := goquery.NewDocumentFromReader(search.Body)
	defer search.Body.Close()
	if err != nil {
		return "newdocumentfromreader fail"
	}

	doc.Find("div.im").Each(func(i int, selection *goquery.Selection) {
		selection.Find("a").Each(func(i int, selection *goquery.Selection) {
			href, exist := selection.Attr("href")
			if !exist {
				println("href 不存在")
			}
			println(href)
		})
	})

	return ""
}
