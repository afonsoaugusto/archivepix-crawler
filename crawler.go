package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"reflect"
)

type Image struct {
	Data        int32
	Url         string
	Description string
}

func handlerError(err error) {
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func main() {
	// https://apod.nasa.gov/apod/astropix.html
	// "https://apod.nasa.gov/apod/archivepix.html"
	// https://apod.nasa.gov/apod/ap190427.html
	resp, err := http.Get("https://apod.nasa.gov/apod/ap190426.html")
	handlerError(err)
	fmt.Println(reflect.TypeOf(resp))
	if resp.StatusCode == http.StatusOK {
		bodyBytes, err := ioutil.ReadAll(resp.Body)
		handlerError(err)
		bodyString := string(bodyBytes)
		fmt.Println(bodyString)
	}
}
