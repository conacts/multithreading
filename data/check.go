package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

func main() {
	resp, err := http.Get("http://localhost:8000/file1.csv")
	if err != nil {
		log.Fatalln(err)
	}
	//We Read the response body on the line below.
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln(err)
	}
	//Convert the body to type string
	sb := string(body)
	fmt.Printf("t1: %T\n", body)
}

func create() {
	f, err := os.Create("data.txt")

	if err != nil {
		log.Fatal(err)
	}

	defer f.Close()

	_, err2 := f.WriteString("old falcon\n")

	if err2 != nil {
		log.Fatal(err2)
	}

	fmt.Println("done")
}
