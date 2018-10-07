package jarhandler

import (
	"net/http"
	"github.com/PuerkitoBio/goquery"
	"strings"
	"path/filepath"
	"github.com/gin-gonic/gin"
	"project/myGradle/src/model"
)

const MAVEN_REPO_HOST = "https://mvnrepository.com/"

type JarHandler struct {
}

func (jarHandler *JarHandler) Handler(jarname string, version string, context *gin.Context) {

	searchJar(jarname, version, context)
}

func searchJar(jarname string, version string, context *gin.Context) {
	search, err := http.Get("https://mvnrepository.com/search?q=" + jarname)
	if err != nil {
		context.JSON(http.StatusOK, model.ErrorResponse{500, "获取maven库结果失败"})
		return
	}

	if search.StatusCode != 200 {
		context.JSON(http.StatusOK, model.ErrorResponse{500, "获取maven库结果失败"})
		return
	}

	doc, err := goquery.NewDocumentFromReader(search.Body)
	defer search.Body.Close()
	if err != nil {
		context.JSON(http.StatusOK, model.ErrorResponse{500, "获取maven库结果失败"})
		return
	}

	isHasJar := false

	doc.Find("div.im").EachWithBreak(func(i int, selection *goquery.Selection) bool {
		selection.Find("a").EachWithBreak(func(i int, selection *goquery.Selection) bool{
			if i == 1 {
				href, exist := selection.Attr("href")
				if exist {
					if isThisJar(href, jarname) {
						jarDetails(href, version, context)
						isHasJar = true
						return false
					}
				}
			}
			return true
		})
		return false
	})


	if !isHasJar {
		context.JSON(http.StatusOK, model.ErrorResponse{500, "不存在这个库"})
	}
}

func isThisJar(href string, jarname string) bool {
	hrefSplit := strings.Split(href, "/")
	name := hrefSplit[len(hrefSplit)-1]
	return name == jarname
}

func jarDetails(href, version string, context *gin.Context) {

	jarDetailsUrl := MAVEN_REPO_HOST + href

	jarDetailsResponse, err := http.Get(jarDetailsUrl)

	if err != nil {
		context.JSON(http.StatusOK, model.ErrorResponse{500, "获取maven库结果失败"})
		return
	}

	if jarDetailsResponse.StatusCode != 200 {
		context.JSON(http.StatusOK, model.ErrorResponse{500, "获取maven库结果失败"})
		return
	}

	doc, err := goquery.NewDocumentFromReader(jarDetailsResponse.Body)
	if err != nil {
		context.JSON(http.StatusOK, model.ErrorResponse{500, "获取maven库结果失败"})
		return
	}

	hasVersion := false
	doc.Find("td").Each(func(i int, selection *goquery.Selection) {
		selection.Find("a").EachWithBreak(func(i int, selection *goquery.Selection) bool {
			if isVersionATag(i, selection, version) {
				hasVersion = true
				versionHref, _ := selection.Attr("href")
				realJarVersionDetailsUrl := getRealJarVersionDetailsUrl(versionHref, jarDetailsUrl)
				jarVersionDetailsUrl(realJarVersionDetailsUrl, context)
				return false
			}
			return true
		})
	})

	if !hasVersion {
		context.JSON(http.StatusOK, model.ErrorResponse{500, "未能在maven仓库中找到对应的版本"})
	}

}

// 是表示版本链接的 a标签
func isVersionATag(i int, selection *goquery.Selection, version string) bool {
	return i == 0 && selection.Text() == version
}

// 根据 jar详情的url 和 详细版本url的最后path 获取完整的版本详细信息的url
// 例如 versionHref 为 gson/2.8.3
//     jarDetailsUrl 为 https://mvnrepository.com/artifact/com.google.code.gson/gson
// 最终结果为 https://mvnrepository.com/artifact/com.google.code.gson/gson/2.8.3
func getRealJarVersionDetailsUrl(versionHref string, jarDetailsUrl string) string {
	pathSplit := strings.Split(versionHref, "/")
	return jarDetailsUrl + string(filepath.Separator) + pathSplit[len(pathSplit)-1]
}

// 爬取具体版本jar的详细信息页
func jarVersionDetailsUrl(realJarVersionDetailsUrl string, context *gin.Context) {
	jarVersionDetailsResponse, err := http.Get(realJarVersionDetailsUrl)
	if err != nil {
		context.JSON(http.StatusOK, model.ErrorResponse{500, "获取maven库结果失败"})
		return
	}

	if jarVersionDetailsResponse.StatusCode != 200 {
		context.JSON(http.StatusOK, model.ErrorResponse{500, "获取maven库结果失败"})
		return
	}

	doc, err := goquery.NewDocumentFromReader(jarVersionDetailsResponse.Body)

	if err != nil {
		context.JSON(http.StatusOK, model.ErrorResponse{500, "获取maven库结果失败"})
		return
	}

	hasVersion := false

	doc.Find("div#gradle-div").Each(func(i int, selection *goquery.Selection) {
		res := selection.Text()
		context.JSON(http.StatusOK, model.StringResponse{200, "获取gradle结果成功", res})
		hasVersion = true
	})

	if !hasVersion {
		context.JSON(http.StatusOK, model.ErrorResponse{500, "未能在maven仓库中找到对应的版本"})
	}
}
