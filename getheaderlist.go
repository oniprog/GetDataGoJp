package main

import (
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"
	"os"
	"path"
	"time"
)

const MAX_PAGE = 618

func main() {

	os.MkdirAll("rare_html", 0777)

	for i := 1; i <= MAX_PAGE; i++ {
		fmt.Printf("Page %d     \r", i)
		url := fmt.Sprintf("http://www.data.go.jp/data/dataset?q=&sort=score+desc%%2C+metadata_modified+desc-20&page=%d", i)
		resp, err := http.Get(url)
		if err != nil {
			fmt.Printf("%s\n", err.Error())
			return
		}
		filepath := path.Join("rare_html", fmt.Sprintf("%d.html", i))
		body, err := ioutil.ReadAll(resp.Body)
		ioutil.WriteFile(filepath, body, 0777)
		resp.Body.Close()

		time.Sleep(time.Duration(10*rand.Float64()*1000+1000) * time.Millisecond)
	}
}
