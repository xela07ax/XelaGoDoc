// Copyright 2015 The Gorilla WebSocket Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build ignore

package main

import (
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

func sh(w http.ResponseWriter, r *http.Request) {
	b, err := ioutil.ReadAll(r.Body) // Считывание тела
	if err != nil {
		ertx := fmt.Sprintf("COM:Ошибка чтения тела: %s | ERTX:can't read body", err)
		fmt.Println(ertx)
		http.Error(w, ertx, http.StatusConflict) // 409
		return
	}
	fmt.Println("Пришла команда по HTTP 1.1")
	sendMessage(string(b)) // Отпаравляем вебсокату
}
func main() {
	http.HandleFunc("/sh", sh) // Будем принимать команды curl
	log.Fatal(http.ListenAndServe("localhost:7456", nil))
}

func sendMessage(msg string) {
	// Инициализируем перехватчик сигнала "Завершение программы", для реализации корректного завершения сессии Websocket
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)
	// Прописываем пути
	u := url.URL{Scheme: "ws", Host: "localhost:5578", Path: "/echo"}
	log.Printf("Начинаем подключение к %s", u.String())
	ws, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		log.Fatal("dial:", err)
	}
	defer ws.Close()
	log.Printf("Подключение к %s удалось", u.String())
	time.Sleep(200 * time.Millisecond)
	done := make(chan struct{}) // создаем канал |отправлять ничего| *Go, практика асинхронного взаимодействия
	// --
	// Реализация ридера
	// в фон функцию, передадим ей этот канал
	i := 0
	go func(i int) {
		defer close(done) // При умирании Горутины канал должен закрыться
		fmt.Println("-R-> Ридер активирован")
		for {             // начнем цикл
			log.Printf("-R-> ",i)
			time.Sleep(1*time.Second)
			i++
			log.Printf("-R-> ReadMessage - activate")
			time.Sleep(1*time.Second)
			_, message, err := ws.ReadMessage() // пришло какое то сообщение с другой стороны WS
			time.Sleep(1*time.Second)
			log.Printf("-R-> ReadMessage - is put")
			if err != nil {
				log.Println("-R-> ошибка при чтении сообщения по ws:", err)
				time.Sleep(1*time.Second)
				return
			}
			time.Sleep(1*time.Second)
			log.Printf("-R-> Cообщение для нас: %s", message)
		}
	}(i)
	time.Sleep(200 * time.Millisecond)
	// в фоне реализация ридера конец
	// --

	// Реализация Writer-а
	fmt.Println("--W-> начнем отправку")
	err = ws.WriteMessage(websocket.TextMessage, []byte(msg))
	if err != nil {
		log.Println("--W-> ошибка при отправке:", err)
		return
	}
	time.Sleep(200*time.Millisecond)
	log.Printf("--W-> отправили:%s\n", msg)
	// ---
	select {
	case <-done:
		log.Println("--W-> пришел \"done\"")
		//return
	case <-interrupt:
		log.Println("--W-> interrupt")

		// Cleanly close the connection by sending a close message and then
		// waiting (with timeout) for the server to close the connection.
		err := ws.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
		if err != nil {
			log.Println("write close:", err)
			return
		}
		select {
		case <-done:
		case <-time.After(time.Second):
		}
		return
	}
}
