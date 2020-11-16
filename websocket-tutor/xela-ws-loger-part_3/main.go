// Copyright 2013 The Gorilla WebSocket Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	bcl "./broadcastLogger"
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
	fmt.Println(addr)
	flag.Parse()
	Logx("-main->start[newHub]")
	hub := newHub()
	go hub.run()
	time.Sleep(1*time.Second)
	logEr := bcl.NewChLoger(&bcl.Config{
		IntervalMs:     300,
		ConsolFilterFn: map[string]int{"Front Http Server":  0},
		ConsolFilterUn: map[string]int{"Pooling": 1},
		Mode:           0,
		Dir:            "x-loger",
		Broadcast: hub.Input,
	})
	logEr.RunMinion()
	Logx("-main->end[newHub]")
	Logx("-main->start[hub.run]-p2")
	time.Sleep(1*time.Second)
	logEr.ChInLog <- [4]string{"Welcome","nil",fmt.Sprintf("Вас приветствует \"Silika-FileManager Контроллер\" v1.1 (11112020) \n")}

	Logx("-main->end[hub.run]")
	// По вебсокетам у нас будет логер
	http.HandleFunc("/", serveHome)
	http.HandleFunc("/wsx/sendMsg", func(w http.ResponseWriter, r *http.Request) {
		sendMsg(logEr.ChInLog, w, r)
	})
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