package main

import (
	"io/ioutil"
	"net/http"
	"regexp"
	"strings"
	"os"

	"github.com/kr/pretty"
)

func main() {
	resp, err := http.Get("http://pl.bab.la/slownik/polski-angielski/" + os.Args[1])
	ifError(err)

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	ifError(err)

	dictionary := retreiveFromResponse(body)

	for _, line := range dictionary {
		pretty.Println(line)
	}
}

func ifError(err error) {
	if err != nil {
		pretty.Println(err)
	}
}

func retreiveFromResponse(body []byte) []string {
	re, _ := regexp.Compile("<a .*?result-link\".*>(.*?)</span>")
	resultLinks := re.FindAllString(string(body[:]), -1)
	removeHTMLTag := strings.Join(resultLinks, "=")
	replace, _ := regexp.Compile("<[^>]*>")
	removeHTMLTag = replace.ReplaceAllString(removeHTMLTag, "")

	result := strings.Split(removeHTMLTag, "=")
	result = append(result[:0], result[1:]...)

	var dictionary []string
	for i := 1; i <= len(result); i += 2 {
		dictionary = append(dictionary, result[i-1]+" = "+result[i])
	}

	return dictionary
}
