package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"os"
)

type Address struct {
	Suggestions []struct {
		Data              struct {
			Address     struct {
				Value             string `json:"value"`
			} `json:"address"`
		} `json:"data"`
	} `json:"suggestions"`
}

type address struct {
	Address string `json:"Address"`
}

type inn struct {
	Inn string `json:"Inn"`
}

func (inn *inn) setInn(body []byte){
	err := json.Unmarshal(body, inn)
	if err != nil {
		panic("Could not parse json from request")
	}
}

func (inn inn) getInn() string{
	return inn.Inn
}

func ip (w http.ResponseWriter, r *http.Request){
	host, _ := os.Hostname()
	var ip string
	addrs, _ := net.LookupIP(host)
	for _, addr := range addrs {
		if ipv4 := addr.To4(); ipv4 != nil && !ipv4.Equal(net.IPv4(127,0,0,1)){
			ip = ipv4.String()
		}
	}
	_, _ = fmt.Fprintf(w, "Host ip is: %s\n", ip)
	fmt.Println("Host ip is: " + ip)
}

func service (w http.ResponseWriter, r *http.Request){
	_, err := fmt.Fprintf(w, "You have reached inn service!")
	if err != nil{
		panic(err)
	}
	fmt.Println("You have reached inn service!")
}

func org (w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		panic("Could not read body from request")
	}
	defer r.Body.Close()
	var inn inn
	inn.setInn(body)

	url := "https://suggestions.dadata.ru/suggestions/api/4_1/rs/findById/party/"
	req, err := http.NewRequest("POST", url, bytes.NewBuffer([]byte(`{ "query": "`+inn.getInn()+`" }`)))
	if err != nil {
		panic("Bad request")
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Authorization", "Token 0397f1e984b6077488c77050f856f97f6cd63e28")
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic("Could not execute request")
	}
	defer resp.Body.Close()
	body, _ = ioutil.ReadAll(resp.Body)
	var addr Address
	err = json.Unmarshal(body, &addr)
	if err != nil {
		panic("Could not unmarshal")
	}
	if len(addr.Suggestions) == 0 {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{ "Error": "wrong inn" }`))
		fmt.Println("Wrong inn")
	} else {
		data := &address{Address: addr.Suggestions[0].Data.Address.Value}
		result, _ := json.Marshal(data)
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write(result)
		fmt.Println(string(result))
	}
}

func main() {
	fmt.Println("Starting server on http://localhost:8080\nTo check whether service is working run http://localhost:8080/service")
	http.HandleFunc("/ip", ip)
	http.HandleFunc("/service", service)
	http.HandleFunc("/org", org)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
