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
type Foo struct { p1, p2 string }
type Bar struct { i1, i2 int }

func (f Foo) Receive() error {
	return print(strings.Join([]string{f.p1, f.p2}, ", "))
}

func (b Bar) Receive() error {
	return print(strconv.Itoa(b.i1 + b.i2))
}

func NewFoo(data map[string]interface{}) (Foo, bool) {
	p1, ok1 := data["p1"].(string);
	p2, ok2 := data["p2"].(string);
	if ok1 && ok2 {
		return Foo{p1, p2}, true
	}
	return Foo{"",""}, false
}

func NewBar(data map[string]interface{}) (Bar, bool) {
	i1, err1 := strconv.Atoi(data["i1"].(string))
	i2, err2 := strconv.Atoi(data["i2"].(string))
	if (err1 == nil && err2 == nil) {
		return Bar{i1, i2}, true
	}
	return Bar{0, 0}, false
}

func handler(w http.ResponseWriter, r *http.Request) {
	var msg_json map[string]interface{}
	if err := json.Unmarshal([]byte(r.URL.Query().Get("msg")), &msg_json); err == nil {
		if data, ok := msg_json["data"].(map[string]interface{}); ok {
			var msg Message
			var msg_ok bool
			if msg_json["type"] == "Foo" {
				msg, msg_ok = NewFoo(data)
			} else if msg_json["type"] == "Bar" {
				msg, msg_ok = NewBar(data)
			}
			
			if msg_ok {
				msg.Receive()
			}
		}
	}
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
