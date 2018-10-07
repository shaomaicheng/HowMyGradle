package jarhandler

import (
	"net/http"
	"github.com/PuerkitoBio/goquery"
	"strings"
	"path/filepath"
)

const MAVEN_REPO_HOST = "https://mvnrepository.com/"

type JarHandler struct {
}

func (jarHandler *JarHandler) Handler(jarname string, version string) string {

	return searchJar(jarname, version)
}

func searchJar(jarname string,version string) string {
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
			if i == 1 {
				href, exist := selection.Attr("href")
				if !exist {
					println("href 不存在")
				}
				jarVersionDetails(href, version)
			}
		})
	})

	return ""
}


func jarVersionDetails(href, version string) {

	jarDetailsUrl := MAVEN_REPO_HOST + href

	jarDetailsResponse, err := http.Get(jarDetailsUrl)

	if err != nil {
		return
	}

	if jarDetailsResponse.StatusCode != 200 {
		return
	}

	doc, err := goquery.NewDocumentFromReader(jarDetailsResponse.Body)
	if err != nil {
		return
	}

	doc.Find("td").Each(func(i int, selection *goquery.Selection) {
		selection.Find("a").Each(func(i int, selection *goquery.Selection) {
			if isVersionATag(i, selection, version) {
				versionHref, _ := selection.Attr("href")
				realJarVersionDetailsUrl := getRealJarVersionDetailsUrl(versionHref, jarDetailsUrl)
				println(realJarVersionDetailsUrl)
			}
		})
	})


}

// 是表示版本链接的 a标签
func isVersionATag(i int, selection * goquery.Selection, version string) bool {
	return i == 0 && selection.Text() == version
}


// 根据 jar详情的url 和 详细版本url的最后path 获取完整的版本详细信息的url
// 例如 versionHref 为 gson/2.8.3
//     jarDetailsUrl 为 https://mvnrepository.com/artifact/com.google.code.gson/gson
// 最终结果为 https://mvnrepository.com/artifact/com.google.code.gson/gson/2.8.3
func getRealJarVersionDetailsUrl(versionHref string, jarDetailsUrl string) string {
	pathSplit := strings.Split(versionHref, "/")
	return jarDetailsUrl + string(filepath.Separator) + pathSplit[len(pathSplit) - 1]
}