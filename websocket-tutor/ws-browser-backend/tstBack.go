// Copyright 2020 The Xela07ax WebSocket Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build ignore

package main

import (
	"flag"
	"fmt"
	"github.com/gorilla/websocket"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"os/signal"
	"time"
)

var addr = flag.String("addr", "localhost:3450", "http service address")

var upgrader = websocket.Upgrader{} // use default options

type Controller struct { // Будем хранить данные о сессии
	unit int
	wsx *websocket.Conn
	done chan struct{} // При закрытии канала, сработает сообщение
}
func main() {
	flag.Parse()
	fmt.Println(*addr)
	cn := &Controller{}
	go cn.WaitFromSignal()
	http.HandleFunc("/conn", cn.Conn)
	http.HandleFunc("/send", cn.Send)

	log.Fatal(http.ListenAndServe(*addr, nil))
}
func (cn *Controller)Conn(w http.ResponseWriter, r *http.Request) {
	var err error
	// Достанем сообщение
	address, err := ioutil.ReadAll(r.Body)
	if err != nil {
		ertx := fmt.Sprintf("COM:Ошибка чтения тела: %s | ERTX:can't read body", err)
		log.Printf("recv: %s\n", ertx)
		http.Error(w, ertx, http.StatusConflict) // 409
		return
	}
	ad := fmt.Sprintf("%s",address)

	if cn.wsx != nil {
		fmt.Println("Подключение уже существует")
		http.Error(w,"Connect already exists!", http.StatusConflict)
		return
	}
	// Подключаемся
	u := url.URL{Scheme: "ws", Host: ad, Path: "/echo"}
	log.Printf("Начинаем подключение к %s", u.String())
	cn.wsx, _, err = websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		ertx := fmt.Sprintf("can't connect dial to address: %s| %v", address, err)
		log.Printf("%s\n", ertx)
		http.Error(w, ertx, http.StatusConflict) // 409
		return
	}
	//defer ws.Close()
	cn.done = make(chan struct{})

	go func() {
		defer close(cn.done)
		for {
			_, message, err := cn.wsx.ReadMessage()
			if err != nil {
				log.Println("read:", err)
				return
			}
			log.Printf("recv: %s", message)
		}
	}()
	tx := fmt.Sprintf("Connect to %s succesfuly", u.String())
	fmt.Fprintf(w, tx)
	log.Println(tx)
}
func (cn *Controller)Send(w http.ResponseWriter, r *http.Request) {
	if cn.wsx == nil {
		tx := "Connect not exists!"
		fmt.Println(tx)
		http.Error(w,tx, http.StatusNotFound)
		return
	}
	// Достанем сообщение
	msg, err := ioutil.ReadAll(r.Body)
	if err != nil {
		ertx := fmt.Sprintf("COM:Ошибка чтения тела: %s | ERTX:can't read body", err)
		log.Printf("recv: %s\n", ertx)
		http.Error(w, ertx, http.StatusConflict) // 409
		return
	}
	// отправим по ws
	fmt.Printf("отправим по ws:%s\n",msg)
	err = cn.wsx.WriteMessage(websocket.TextMessage, msg)
	if err != nil {
		log.Println("write:", err)
		return
	}
}

func (cn *Controller)WaitFromSignal()  {
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)
	<-interrupt
	log.Println("interrupt")
	if cn.wsx == nil {
		log.Println("connect ws not opened")
		os.Exit(0)
	}
	// Cleanly close the connection by sending a close message and then
	err := cn.wsx.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, "interrupt close"))
	if err != nil {
		log.Println("write close:", err)
	}
	select {
	case <-cn.done:
		log.Println("interrupt close - done")
	case <-time.After(time.Second):
		log.Println("interrupt close - timeout")
	}
	os.Exit(0)
}

func wsController(b []byte) {
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)

	u := url.URL{Scheme: "ws", Host: *addr, Path: "/echo"}
	log.Printf("connecting to %s", u.String())

	c, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		log.Fatal("dial:", err)
	}
	defer c.Close()

	done := make(chan struct{})

	go func() {
		defer close(done)
		for {
			_, message, err := c.ReadMessage()
			if err != nil {
				log.Println("read:", err)
				return
			}
			log.Printf("recv: %s", message)
		}
	}()

	for {
		select {
		case <-done:
			return
		case t := <-make(chan string):
			err := c.WriteMessage(websocket.TextMessage, []byte(t))
			if err != nil {
				log.Println("write:", err)
				return
			}
		case <-interrupt:
			log.Println("interrupt")

			// Cleanly close the connection by sending a close message and then
			// waiting (with timeout) for the server to close the connection.
			err := c.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
			if err != nil {
				log.Println("write close:", err)
				return
			}

		}
	}
}