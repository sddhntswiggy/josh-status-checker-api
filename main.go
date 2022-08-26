package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"
)

var websitestatus = make(map[string]string)

func StatusCheck(weblink string, c chan string) {
	_, err := http.Get(weblink)
	if err != nil {
		fmt.Println(weblink, ": Invalid, Website not UP")
		websitestatus[weblink] = "DOWN"
		c <- weblink
	} else {
		fmt.Println(weblink, ": Status 200 OK")
		websitestatus[weblink] = "UP"
		c <- weblink
	}
}
func handler(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application.json")

	if r.Method == "GET" {
		params := r.URL.Query().Get("name")
		if params != "" {
			status, _ := websitestatus[params]
			fmt.Fprintln(w, fmt.Sprintf("Status of %s is %s", params, status))
		} else {
			for site, status := range websitestatus {
				fmt.Fprintln(w, site, "is", status)
			}
		}
	} else if r.Method == "POST" {
		var websites []string
		err := json.NewDecoder(r.Body).Decode(&websites)
		if err != nil {
			log.Fatal("error occured with decoding website :", err)
		}

		c := make(chan string)
		for _, links := range websites {
			go StatusCheck(links, c)
		}
		for url := range c {
			go func(s string) {
				time.Sleep(60 * time.Second)
				StatusCheck(s, c)
			}(url)
		}
	} else {
		log.Fatal("Invalid Request")
	}
}
func main() {
	fmt.Println("Server starting on port :8080")
	http.HandleFunc("/", handler)
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}
}
