// Copyright 2013 The Gorilla WebSocket Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"time"
)

var addr = flag.String("addr", ":8180", "http service address")

func serveHome(w http.ResponseWriter, r *http.Request) {
	log.Println(r.URL)
	if r.URL.Path != "/" {
		http.Error(w, "Not found", http.StatusNotFound)
		return
	}
	if r.Method != "GET" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	http.ServeFile(w, r, "home.html")
}

func main() {
	flag.Parse()
	Logx("-main->start[newHub]")
	hub := newHub()
	Logx("-main->end[newHub]")
	Logx("-main->start[hub.run]")
	go hub.run()
	Logx("-main->end[hub.run]")
	http.HandleFunc("/", serveHome)
	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		serveWs(hub, w, r)
	})
	err := http.ListenAndServe(*addr, nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
	Logx("-main->end")
}

func Logx(txt string)  {
	fmt.Printf("%s|%v\n",txt,Getime())
	time.Sleep(1*time.Second)
}
func Getime()string  {
	return time.Now().Format("2006-01-02 15:04:05")
}