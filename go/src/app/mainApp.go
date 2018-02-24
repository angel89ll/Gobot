package main

import "fmt"
import "net/http"
import "io/ioutil"

func main() {

	url := "https://golang.org/doc/"
	resp, _ := http.Get(url)
	bytes, _ := ioutil.ReadAll(resp.Body)

	fmt.Println("HTML:\n\n", string(bytes))
	resp.Body.Close()
}
