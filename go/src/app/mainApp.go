package main

import (
	"fmt"
	"net/http"
	"os"
	"golang.org/x/net/html"
	"path/filepath"
	"bufio"
)


//Extract useful data from a link
func extractType1Data(pUrl string, pChFinish chan bool) {
	resp, error := http.Get(pUrl)

	defer func() {
		// Notify that routine is done
		pChFinish <- true
	} ()

	if error != nil {
		fmt.Println("Error getting url "+ pUrl)
		return
	}

	body := resp.Body
	defer body.Close() // Close body when function returns

	tokenizer := html.NewTokenizer(body)

	for {
		tItem := tokenizer.Next()

		switch {
		case tItem  == html.ErrorToken:
			// End of the document
			fmt.Println("End Of doc")
			return
		case tItem == html.StartTagToken:
			//Token found
			//TODO: Crear funcion para obtener los atributos del tag que contenga la info.
			//TODO: Una vez obtenida la info, almacenarla en bbdd 
			token := tokenizer.Token()
			isSearchedTag := token.Data == "h4"
			if isSearchedTag {
				fmt.Println("We found a h4!!!!!")
			}
		}
	}
}

func main() {

	absPath, error := filepath.Abs("../data/urlFile.txt")
	if error != nil {
		panic(error)
	}

	file, error := os.Open(absPath);
	if error != nil {
		panic(error)
	}
	defer file.Close()

	channelFinishedGoroutine := make(chan bool)

	fileScanner := bufio.NewScanner(file)
	numberGoroutines := 0
	for fileScanner.Scan() {
		url := fileScanner.Text()
		fmt.Println(url)
		if numberGoroutines == 3 {
			select {
			case <-channelFinishedGoroutine: // Wait until another finish
				numberGoroutines = numberGoroutines - 1
			}
		}
		go extractType1Data(url, channelFinishedGoroutine)
		numberGoroutines = numberGoroutines + 1
	}

	for ;numberGoroutines > 0; { //Wait until finish all
		select {
			case <-channelFinishedGoroutine:
				numberGoroutines = numberGoroutines - 1
			}
	}
	
}
