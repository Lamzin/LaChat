package main

import (
	"log/syslog"
	"fmt"
	"strconv"
	"strings"
	"encoding/json"
	"net/http"
)

type Message interface { Receive() error }

type Foo struct { 
	P1 string `json:"p1"`
	P2 string `json:"p2"`
}

type FooMessage struct { 
	Type string `json:"type"`
	Data Foo `json:"data"`
}

type Bar struct { 
	I1 int `json:"i1"`
	I2 int `json:"i2"`
}

type BarMessage struct { 
	Type string `json:"type"`
	Data Bar `json:"data"`
}

type Msg struct {
	Type string `json:"type"`
	Data interface{} `json:"data"`
}

func (f Foo) Receive() error {
	return print(strings.Join([]string{f.P1, f.P2}, ", "))
}

func (b Bar) Receive() error {
	return print(strconv.Itoa(b.I1 + b.I2))
}

func handler(w http.ResponseWriter, r *http.Request) {
	request_message := []byte(r.URL.Query().Get("msg"))

	var msg_json map[string]interface{}
	if err := json.Unmarshal(request_message, &msg_json); err != nil {
		return
	}
		
	var msg Message
	if msg_json["type"] == "Foo" {
		var foomsg FooMessage
		if err := json.Unmarshal(request_message, &foomsg); err != nil {
			return
		}
		msg = foomsg.Data
	} else if msg_json["type"] == "Bar" {
		var barmsg BarMessage
		if err := json.Unmarshal(request_message, &barmsg); err != nil {
			return
		}
		msg = barmsg.Data
	}
	
	msg.Receive()
}

var logger, _ = syslog.New(syslog.LOG_INFO, "GoServer")
func print(param string) error {
	logger.Info(param)
	fmt.Println(param)
	return nil
}

func main() {
    http.HandleFunc("/in", handler)
    http.ListenAndServe(":8080", nil)
	return
}

// localhost:8080/in?msg={"type":"Foo", "data":{"p1":"Foo", "p2":"message"}}
// localhost:8080/in?msg={"type":"Bar", "data":{"i1":"33", "i2":"66"}}
