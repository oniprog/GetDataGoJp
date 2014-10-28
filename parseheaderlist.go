package main

import (
	"bufio"
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"strings"
)

const LINK_HEAD = "http://www.data.go.jp"

func main() {

	files, err := ioutil.ReadDir("rare_html")
	if err != nil {
		fmt.Printf("%s\n", err.Error())
		return
	}

	//fw,err := os.Create("header.csv")
	fw, err := os.Create("header.html")
	if err != nil {
		fmt.Printf("%s\n", err.Error())
		return
	}

	wr := bufio.NewWriter(fw)

	wr.WriteString("<!DOCTYPE html><html><head></head><body>")

	for _, fi := range files {

		if fi.IsDir() {
			continue
		}
		fmt.Printf("%s    \r", fi.Name())

		path := path.Join("rare_html", fi.Name())
		buf, err := ioutil.ReadFile(path)
		if err != nil {
			fmt.Printf("%s\n", err.Error())
			return
		}

		r := bufio.NewReader(bytes.NewReader(buf))
		mode := 0
		for {
			line, err := r.ReadString('\n')
			if err != nil {
				break
			}
			if mode == 0 {
				if strings.Index(line, "class=\"dataset-item\"") >= 0 {
					mode = 1
				}
			} else if mode == 1 {

				if i1 := strings.Index(line, "a href=\""); i1 >= 0 {
					i2 := strings.LastIndex(line, "</a>")
					i3 := strings.LastIndex(line, "\">")
					if i2 < 0 {
						i2 = len(line)
					}
					if i3 < 0 {
						fmt.Printf("error %d,%d %s\n", i2, i3, line)
						continue
					}
					linkstr := line[i1+8 : i3]
					headstr := line[i3+2 : i2]
					//fmt.Printf("%s:%s\n", linkstr, headstr)
					//wr.WriteString( fmt.Sprintf("\"%s\", \"%s%s\" ", headstr, LINK_HEAD, linkstr ))
					wr.WriteString(fmt.Sprintf("<H2><A HREF=\"%s%s\">%s</A></H2>\n", LINK_HEAD, linkstr, headstr))
				} else if strings.Index(line, "</h3>") >= 0 {
					mode = 2
				}
			} else if mode == 2 {

				if i1 := strings.Index(line, "data-format=\""); i1 >= 0 {
					if i2 := strings.Index(line, "\">"); i2 >= 0 {
						datatypestr := line[i1+13 : i2]
						//fmt.Printf(" link : %s\n", datalinkstr)
						//wr.WriteString( fmt.Sprintf(",\"%s%s\" ", LINK_HEAD, datalinkstr))
						wr.WriteString(fmt.Sprintf("<q>%s</q>\n", datatypestr))
					}
				} else if strings.Index(line, "class=\"metadata\">") >= 0 {
					mode = 3
				}
			} else if mode == 3 {

				if strings.Index(line, "</li>") >= 0 {
					wr.WriteString("\r\n")
					mode = 0
				}
			}
		}
	}
	wr.WriteString("</body></html>\n")
	wr.Flush()
	fw.Close()
}
