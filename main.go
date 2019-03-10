package main

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os/exec"
	"strings"
	"sync"

	"golang.org/x/net/html"
)

func main() {
	doc, _ := html.Parse(strings.NewReader(makeRequest("https://mednafen.github.io/")))
	tag, err := getItem(doc, "title")
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	item := renderNode(tag)
	fmt.Printf("Item : %s", item)
}

func makeRequest(URL string) string {
	resp, err := http.Get(URL /*"https://mednafen.github.io/"*/)
	if err != nil {
		log.Fatalln(err)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln(err)
	}

	return string(body)
}
func getItem(doc *html.Node, tagName string) (*html.Node, error) {
	var b *html.Node
	var f func(*html.Node)
	f = func(n *html.Node) {
		if n.Type == html.ElementNode && n.Data == tagName {
			b = n
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			f(c)
		}
	}
	f(doc)
	if b != nil {
		return b, nil
	}
	return nil, errors.New(fmt.Errorf("Missing item : '%s' in the node tree", tagName).Error())
}

func renderNode(n *html.Node) string {
	var buf bytes.Buffer
	w := io.Writer(&buf)
	html.Render(w, n)
	return buf.String()
}

func printName(jString string, wg *sync.WaitGroup) {
	cmd := exec.Command(fmt.Sprintf("{ sleep 4; echo %s;}", jString))
	cmdReader, err := cmd.StdoutPipe()
	if err != nil {
		log.Fatal(err)
	}
	scanner := bufio.NewScanner(cmdReader)
	go func() {
		for scanner.Scan() {
			fmt.Println(scanner.Text())
		}
	}()
	if err := cmd.Start(); err != nil {
		log.Fatal(err)
	}
	if err := cmd.Wait(); err != nil {
		log.Fatal(err)
	}
	wg.Done()
}
