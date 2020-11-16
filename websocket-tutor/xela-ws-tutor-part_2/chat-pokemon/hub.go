// Copyright 2013 The Gorilla WebSocket Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import "fmt"

type Hub struct {
	// Зарегистрированный клиент.
	clients map[*Client]bool
	Input chan []byte
	// Входящие сообщения от клиентов.
	broadcast chan []byte

	// Регистрируйте запросы от клиентов.
	register chan *Client

	// Отмените регистрацию запросов от клиентов.
	unregister chan *Client
}

func newHub() *Hub {
	return &Hub{
		broadcast:  make(chan []byte),
		Input:      make(chan []byte,100),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		clients:    make(map[*Client]bool),
	}
}

func (h *Hub) run() {
	Logx("-hub.run->init")
	for {
		Logx("-hub.run->for[circle-start]")
		select {
		case client := <-h.register:
			Logx("-hub.run->select[h.register]")
			h.clients[client] = true
		case client := <-h.unregister:
			Logx("-hub.run->select[h.unregister]")
			if _, ok := h.clients[client]; ok {
				Logx("-hub.run->select[h.unregister]-ok")
				delete(h.clients, client)
				close(client.send)
			}
		case message := <-h.Input:
			Logx("-hub.run->select[h.Input]new")
			for client := range h.clients {
				Logx("-hub.run->select[h.Input]for")
				select {
				case client.send <- message:
					Logx(fmt.Sprintf("-hub.run->select[h.Input]for-select[%s]",message))
				default:
					Logx("-hub.run->select[h.Input]for-select-default")
					close(client.send)
					delete(h.clients, client)
				}
			}
		case message := <-h.broadcast:
			// Входящие сообщения от клиентов.
			// Собственно работать будем тут
			Logx("-hub.run->select[h.broadcast]")
			for client := range h.clients {
				Logx("-hub.run->select[h.broadcast]for")
				select {
				case client.send <- message:
					Logx(fmt.Sprintf("-hub.run->select[h.broadcast]for-select[%s]",message))

				default:
					Logx("-hub.run->select[h.broadcast]for-select-default")
					close(client.send)
					delete(h.clients, client)
				}
			}
		}
	}
}
