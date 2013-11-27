package main

import (
	"fmt"
	"net/http"
	"net/url"
	"encoding/json"
	"log"
	"strings"
)

func authenticate(w http.ResponseWriter, r *http.Request) {

	var code string
	code = r.URL.Path[6:] // github authed temp code
	log.Println(code)
	
	w.Header().Set("Access-Control-Allow-Origin", "*")
  	w.Header().Set("Access-Control-Allow-Methods", "GET, OPTIONS")
  	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
  	w.Header().Set("Content-Type", "application/json")

  	// create proxy request to Github.com
  	client := &http.Client{}

	v := url.Values{}
	v.Set("code", code)
	v.Add("client_id", "0975f5317506bd6da14b")
	v.Add("client_secret", "170e1c7d0a61b4d0479125a66a5badced1d2ed97")
	req, err := http.NewRequest("POST", "https://github.com/login/oauth/access_token", strings.NewReader(v.Encode()))
	if err != nil {
		http.Error(w, err.Error(), 500)
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded; param=value")
	req.Header.Add("Accept", "application/json")
	resp, err := client.Do(req)

	if err != nil {
    	http.Error(w, err.Error(), 500)
	}
	defer resp.Body.Close()
	var jsonMap interface{}
    json.NewDecoder(resp.Body).Decode(&jsonMap)

	log.Println(jsonMap)
    b, err := json.Marshal(jsonMap)
    w.Write(b)
}

func baseHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello %s", r.URL.Path[1:])
}

func main() {
	http.HandleFunc("/", baseHandler)
	http.HandleFunc("/auth/", authenticate)
	http.ListenAndServe(":8080", nil)
}