// Copyright 2013 The Gorilla WebSocket Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"flag"
	"fmt"
	"html/template"
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
func uploadHome(w http.ResponseWriter, r *http.Request) {
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

var homeTemplate = template.Must(template.New("").Parse(`
<html>
<head>
<meta http-equiv="Content-Type" content="text/html; charset=UTF-8">
<title>Upload Files</title>
</head>
<body>
    <h2>File Upload</h2>
    Select file
    <input type="file" id="filename" />
    <br>
    <input type="button" value="Connect" onclick="connectChatServer()" />
    <br>
    <input type="button" value="Upload" onclick="sendFile()" />
    <script>
        var ws;
        function connectChatServer() {
            ws = new WebSocket("ws://{{.}}/common");
            ws.binaryType = "arraybuffer";
            ws.onopen = function() {
                alert("Connected.")
            };
            ws.onmessage = function(evt) {
                alert(evt.msg);
            };
            ws.onclose = function() {
                alert("Connection is closed...");
            };
            ws.onerror = function(e) {
                alert(e.msg);
            }
        }
        function sendFile() {
            var file = document.getElementById('filename').files[0];
            var reader = new FileReader();
            var rawData = new ArrayBuffer();
            reader.loadend = function() {
            }
            reader.onload = function(e) {
                rawData = e.target.result;
                ws.send(rawData);
                alert("the File has been transferred.")
            }
            reader.readAsArrayBuffer(file);
        }
    </script>
</body>
</html>
`))

func main() {
	flag.Parse()
	Logx("-main->start[newHub]")
	hub := newHub()
	Logx("-main->end[newHub]")
	Logx("-main->start[hub.run]")
	go hub.run()
	Logx("-main->end[hub.run]")
	http.HandleFunc("/", serveHome)
	http.HandleFunc("/up", serveHome)
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