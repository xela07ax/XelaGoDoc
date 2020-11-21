package main

import (
	"flag"
	"fmt"
	"github.com/xela07ax/toolsXela/chLogger"
	"log"
	"net/http"
	"time"
)

var addr = flag.String("addr", ":8180", "http service address")

func main() {
	logEr := chLogger.NewChLoger(&chLogger.Config{
		IntervalMs:     300,
		ConsolFilterFn: map[string]int{"Front Http Server":  0},
		ConsolFilterUn: map[string]int{"Pooling": 1},
		Mode:           0,
		Dir:            "chloger",
	})
	logEr.RunMinion()
	time.Sleep(1*time.Second)
	logEr.ChInLog <- [4]string{"Welcome","nil",fmt.Sprintf("Вас приветствует \"%s\"\n",*addr)}
	k := &Kod{}
	// Сделаем сначала управление по HTTP
	http.HandleFunc("/kod/stopContainer", k.StopContainer)
	http.HandleFunc("/kod/runContainer", k.RunContainer)
	err := http.ListenAndServe(*addr, nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
