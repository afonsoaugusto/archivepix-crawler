package main

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"math/rand"
	"net/http"
	"os"
	"os/exec"
	"reflect"
	"strings"
	"time"
)

type Image struct {
	name string
	url  string
}

func launchError(text string, funcName string) {
	handlerError(errors.New(text), funcName)
}
func handlerError(err error, funcName string) {
	if err != nil {
		fmt.Println(err, funcName)
		os.Exit(1)
	}
}

func main() {
	urlPrin := "https://apod.nasa.gov/apod/archivepix.html"
	pagesUrls := getContentUrl(urlPrin)
	urls := getLinks(pagesUrls)
	rand.Seed(time.Now().UnixNano())
	choicePage := extractPageName(urls[rand.Intn(len(urls))])
	fmt.Println(choicePage)
	baseUrl := "https://apod.nasa.gov/apod/"
	url := baseUrl + choicePage
	fmt.Println(url)

	content := getContentUrl(url)
	image := getImage(content)
	fmt.Printf("%+v\n", image)

	fileUrl := baseUrl + image.url
	fmt.Println(fileUrl)
	downloadFile(fileUrl)
	changeWallPaper()
}

func getIndex(s string) int {
	return strings.Index(s, "\"")
}
func extractPageName(text string) string {
	text = text[getIndex(text)+1:]
	text = text[:getIndex(text)]
	fmt.Println(text)
	return text
}
func getLinks(page string) []string {
	var urls []string
	//"<a href=\"ap"
	scanner := bufio.NewScanner(strings.NewReader(page))
	for scanner.Scan() {
		text := scanner.Text()
		if checkILineLinkToPageWithImage(text) {
			urls = append(urls, removeCharacters(text))
		}
	}
	return urls
}
func getContentUrl(url string) string {
	var bodyString string
	resp, err := http.Get(url)
	handlerError(err, "getContentUrl")
	fmt.Println(reflect.TypeOf(resp))
	bodyString = extactContentFromResponse(resp)
	return bodyString
}

func extactContentFromResponse(resp *http.Response) string {
	if resp.StatusCode != http.StatusOK {
		launchError("resp.StatusCode != http.StatusOK", "extactContentFromResponse")
	}
	bodyBytes, err := ioutil.ReadAll(resp.Body)
	handlerError(err, "extactContentFromResponse")
	return string(bodyBytes)
}

func getImage(page string) Image {
	var urls []string
	var names []string
	scanner := bufio.NewScanner(strings.NewReader(page))
	for scanner.Scan() {
		text := scanner.Text()
		if checkIfUrlImageInText(text) {
			urls = append(urls, removeCharacters(text))
		}
		if checkNameInText(text) {
			names = append(names, removeCharacters(text))
		}
	}
	return Image{names[0], urls[0]}
}

func checkNameInText(text string) bool {
	return checkText(text, "<title>")
}

func checkILineLinkToPageWithImage(text string) bool {
	return checkText(text, "<a href=\"ap")
}

func checkIfUrlImageInText(text string) bool {
	return checkText(text, "<a href=\"image")
}

func checkText(text string, condition string) bool {
	if strings.Contains(text, condition) {
		return true
	}
	return false
}

func removeCharacters(text string) string {
	var stringsForReplace []string
	stringsForReplace = append(stringsForReplace, "<IMG SRC=")
	stringsForReplace = append(stringsForReplace, "<a href=")
	stringsForReplace = append(stringsForReplace, "<title>")
	stringsForReplace = append(stringsForReplace, "/a>")
	stringsForReplace = append(stringsForReplace, "<br")
	stringsForReplace = append(stringsForReplace, ">")
	stringsForReplace = append(stringsForReplace, "<")

	for _, strFPc := range stringsForReplace {
		text = strings.Replace(text, strFPc, "", 1)
	}
	text = removeInvertedComma(text)
	return text
}

func removeInvertedComma(s string) string {
	if strings.Contains(s[0:1], "\"") {
		return s[1 : len(s)-1]
	}
	return s
}

func downloadFile(url string) {
	response, e := http.Get(url)
	handlerError(e, "downloadFile-http.Get")
	defer response.Body.Close()

	//open a file for writing
	file, err := os.Create("/tmp/atual.jpg")
	handlerError(err, "downloadFile-os.Create")
	defer file.Close()

	_, err = io.Copy(file, response.Body)
	handlerError(err, "downloadFile-io.Copy")
	fmt.Println("Success!")
}

func changeWallPaper() {
	cmd := "gsettings set org.gnome.desktop.background picture-uri file:///tmp/atual.jpg"
	out, err := exec.Command("/bin/bash", "-c", cmd).Output()
	handlerError(err, "changeWallPaper-exec.Command")
	fmt.Printf("%s", out)
}
