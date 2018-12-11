package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
)

type Config struct {
	Webhook_server string `json:webhook_server`
	Webhook_port   string `json:webhook_port`
	TGToken        string `json:"tgToken"`
	TGGroup        string `json:"tgGroup"`
}

type Alerts struct {
	Version           string                 `json:"version"`
	GroupKey          string                 `json:"groupKey"`
	Status            string                 `json:"status"`
	Receiver          string                 `json:"receiver"`
	GroupLabels       map[string]interface{} `json:"groupLabels"`
	CommonLabels      map[string]interface{} `json:"commonLabels"`
	CommonAnnotations map[string]interface{} `json:"commonAnnotations"`
	ExternalURL       string                 `json:"externalURL"`
	Alerts            []Alert                `json:"alerts"`
}

type Alert struct {
	Status       string                 `json:"status"`
	Labels       map[string]interface{} `json:"labels"`
	Annotations  map[string]interface{} `json:"annotations"`
	StartsAt     string                 `json:"startsAt"`
	EndsAt       string                 `json:"endsAt"`
	GeneratorURL string                 `json:"generatorURL"`
}

var config Config

func init() {
	file := flag.String("file", "./config.json", "")
	flag.Parse()
	load(*file)
}

func load(file string) {
	cFile, err := os.Open(file)
	defer cFile.Close()
	if err != nil {
		panic(err)
	}
	decoder := json.NewDecoder(cFile)
	err = decoder.Decode(&config)
	if err != nil {
		panic(err)
	}
}

func alerts(w http.ResponseWriter, r *http.Request) {
	log.Println("alerts received")
	decoder := json.NewDecoder(r.Body)
	var msg Alerts
	err := decoder.Decode(&msg)
	if err != nil {
		panic(err)
	}
	w.WriteHeader(200)
	for i, v := range msg.Alerts {
		text := fmt.Sprintf("%s", v)
		log.Println("#", i, text)
		// simply forward alerts to tg
		// todo: better text format
		notify(text)
	}
}

func notify(text string) {
	var buf bytes.Buffer
	buf.WriteString("https://api.telegram.org/bot")
	buf.WriteString(config.TGToken)
	buf.WriteString("/sendMessage?chat_id=-")
	buf.WriteString(config.TGGroup)
	buf.WriteString("&text=")
	buf.WriteString(text)
	resp, err := http.Get(buf.String())
	if err != nil {
		log.Println(err)
	}
	defer resp.Body.Close()
	if resp.StatusCode == 200 {
		log.Println("alerts sent to telegram")
	}
}

func main() {
	http.HandleFunc("/alerts", alerts)
	addr := config.Webhook_server + ":" + config.Webhook_port
	panic(http.ListenAndServe(addr, nil))
}
