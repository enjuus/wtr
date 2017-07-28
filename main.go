package main

import (
	"fmt"
	flag "github.com/ogier/pflag"
	"github.com/yhat/scrape"
	"golang.org/x/net/html"
	"io/ioutil"
	"net/http"
	"os"
)

const link string = "" // Add a wunderground link here
const weatherFile string = "/tmp/weatherfile"

var update bool
var help bool

func init() {
	flag.BoolVarP(&help, "help", "h", false, "Display this help message")
	flag.BoolVarP(&update, "update", "u", false, "Update the offline weather file")
	flag.Parse()
}

func main() {
	if help == true {
		PrintHelpMessage()
	}
	if update == true {
		fmt.Printf("%s", GetNewWeatherData())
	} else {
		fmt.Printf("%s", ReadWeatherData())
	}
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func WriteToFile(data string) bool {
	inputString := []byte(data)
	err := ioutil.WriteFile(weatherFile, inputString, 0644)
	check(err)

	return true

}

func ReadWeatherData() string {
	if _, err := os.Stat(weatherFile); os.IsNotExist(err) {
		return string(GetNewWeatherData())
	} else {
		data, err := ioutil.ReadFile(weatherFile)
		check(err)
		return string(data)
	}
}

func GetNewWeatherData() string {
	resp, err := http.Get(link)
	if err != nil {
		panic(err)
	}
	root, err := html.Parse(resp.Body)
	if err != nil {
		panic(err)
	}

	temp, _ := scrape.Find(root, scrape.ById("curTemp"))
	cond, _ := scrape.Find(root, scrape.ById("curCond"))

	data := fmt.Sprintf("%s - %s", scrape.Text(temp), scrape.Text(cond))
	WriteToFile(data)
	return data
}

func PrintHelpMessage() {
	fmt.Printf("Usage: %s [options]\n", os.Args[0])
	fmt.Printf("Options:\n")
	flag.PrintDefaults()
	os.Exit(1)
}
