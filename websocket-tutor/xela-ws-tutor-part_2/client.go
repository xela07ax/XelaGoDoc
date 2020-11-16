// Copyright 2015 The Gorilla WebSocket Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.


package main

import (
	"bytes"
	"fmt"
	"github.com/gorilla/websocket"
	"github.com/xela07ax/toolsXela/tp"
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
	sendMessageOld(string(b)) // Отпаравляем вебсокату
}
func readBodySimple(w http.ResponseWriter, r *http.Request) []byte {
	fmt.Println("Пришла команда по HTTP 1.1")
	b, err := ioutil.ReadAll(r.Body) // Считывание тело, ожидаем адрес сервера "localhost:5578"
	if err != nil {
		ertx := fmt.Sprintf("COM:Ошибка чтения тела: %s | ERTX:can't read body", err)
		fmt.Println(ertx)
		http.Error(w, ertx, http.StatusConflict) // 409
		return []byte{}
	}
	return b
}
func main() {
	wsxt := NewWsx()
	http.HandleFunc("/sh", sh) // Будем принимать команды curl
	http.HandleFunc("/wsx/connect", wsxt.NewConnWsx)
	http.HandleFunc("/wsx/sendMsg", wsxt.SendMsg)
	http.HandleFunc("/wsx/status", wsxt.ConnStatus)
	log.Fatal(http.ListenAndServe("localhost:7456", nil))
}

// Для взаимодействия нужен объект
type Wsx struct {
	Conn *websocket.Conn
}

func NewWsx() *Wsx {
	return &Wsx{}
}
type Notify struct {
	FuncName string
	Text string
	Status int
	Show bool
	UpdNum int
}

func resp (w http.ResponseWriter, r *http.Request, funcName string,text string, status int, show bool) {
	Notify := Notify {
		FuncName: funcName,
		Text:     text,
		Status:   status,
		Show:     show,
	}
	if err := tp.Httpjson(w, r, Notify); err != nil {
		log.Fatalf("Критическая ошибка, не удалось отправить сообщение в UI: %s| %v", err,Notify)
	}

}
func (wsx *Wsx) ConnStatus(w http.ResponseWriter, r *http.Request) {
	fn := "ConnStatus"
	// Подразумеваем, что подключение у нас создалось другой командой
	// отправим тестовое сообщение


	ertx := "Вэбсокет коннект успешен"
	resp(w,r,fn, ertx, 0, true)
}
func (wsx *Wsx) SendMsg(w http.ResponseWriter, r *http.Request) {
	fn := "ConnStatus"
	// Подразумеваем, что подключение у нас создалось другой командой
	// отправим тестовое сообщение
	msgRaw := readBodySimple(w,r)
	if len(msgRaw) == 0 {
		return
	}
	// Реализация Writer-а
	fmt.Println("--W-process> начнем отправку")
	err := wsx.Conn.WriteMessage(websocket.TextMessage, msgRaw)
	if err != nil {
		log.Println("--W-false> errtx при отправке:", err)
		return
	}
	time.Sleep(200*time.Millisecond)
	log.Printf("--W-true> отправили:%s\n", msgRaw)
	// ---
	ertx := "Вэбсокет коннект успешен"
	resp(w,r,fn, ertx, 0, true)
}
func (wsx *Wsx) NewConnWsx(w http.ResponseWriter, r *http.Request) {
	fn := "NewConnWsx"
	// пока не проверяем есть оно уже или нет, и так знаем, что ничего пока нет!
	// так же не проверяем правильность адреса!
	msgRaw := readBodySimple(w,r)
	if len(msgRaw) == 0 {
		return
	}
	// Прописываем пути
	u := url.URL{Scheme: "ws", Host: string(msgRaw), Path: "/ws"} // todo: надо будет поправить маршрут сервера
	log.Printf("Начинаем подключение к %s", u.String())

	ws, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		log.Fatal("dial:", err)
	}
	wsx.Conn = ws
	// вотсюда ртдер
	go wsx.readPump()
	// отправим какое-нибудь сообщение
	ertx := "[true] -Вэбсокет коннект успешен"
	resp(w,r,fn, ertx, 0, true)
}
const (
	// Time allowed to write a message to the peer.
	writeWait = 10 * time.Second

	// Time allowed to read the next pong message from the peer.
	pongWait = 60 * time.Second

	// Send pings to peer with this period. Must be less than pongWait.
	pingPeriod = (pongWait * 9) / 10

	// Maximum message size allowed from peer.
	maxMessageSize = 512
)
var (
	newline = []byte{'\n'}
	space   = []byte{' '}
)
// readPump pumps messages from the websocket connection to the hub.
//
// The application runs readPump in a per-connection goroutine. The application
// ensures that there is at most one reader on a connection by executing all
// reads from this goroutine.
func (c *Wsx) readPump() {
	//defer func() {
	//	c.hub.unregister <- c
	//	c.Conn.Close()
	//}()
	
	c.Conn.SetReadLimit(maxMessageSize)
	c.Conn.SetReadDeadline(time.Now().Add(pongWait))
	c.Conn.SetPongHandler(func(string) error { c.Conn.SetReadDeadline(time.Now().Add(pongWait)); return nil })
	for {
		_, message, err := c.Conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("error: %v", err)
			}
			break
		}
		message = bytes.TrimSpace(bytes.Replace(message, newline, space, -1))
		fmt.Printf("--rw->%s\n",message)
	}
}

func sendMessageOld(msg string) {
	// Инициализируем перехватчик сигнала "Завершение программы", для реализации корректного завершения сессии Websocket
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)
	// Прописываем пути
	u := url.URL{Scheme: "ws", Host: "localhost:8180", Path: "/echo"}
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
