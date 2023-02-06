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
	fmt.Printf("t1: %T\n", body)
	writeFile(body)
}

func writeFile(s []byte) string {
	f, err := os.Create("online.csv")
	if err != nil {
		log.Fatalln(err)
	}
	defer f.Close()

	f.Write(s)
	return ""
}
