package ping

import (
	"fmt"
	"log"
	"net/http"
)

type Server struct {
	Port      string
	OnPing    func(id string, msg string) (idout string, msgout string, err error)
	OnMessage func(id, msg string) (msgout string, err error)
	ChannelManager
	Responder
}

var (
	channels       = make(map[string]([]chan string))
	DEFAULT_ONPING = func(id, msg string) (idout, msgout string, err error) {
		return id, msg, nil
	}
	DEFAULT_ONMESSAGE = func(id, msg string) (msgout string, err error) {
		return msg, nil
	}
)

func (ps *Server) Serve() error {
	http.HandleFunc("/listen", ps.listen)
	http.HandleFunc("/ping", ps.ping)
	return http.ListenAndServe(":"+ps.Port, nil)
}

func NewServer(port string) *Server {
	return &Server{port, DEFAULT_ONPING, DEFAULT_ONMESSAGE, NewChannelManager(), 0}
}

func (ps *Server) listen(rw http.ResponseWriter, req *http.Request) {
	f, cn := checkStreamable(rw)
	if cn == nil {
		fmt.Fprintf(rw, `{"success":0}`)
		return
	}

	setHeaders(rw)

	id, ok := parseListen(req)
	if !ok {
		fmt.Fprintf(rw, `{"success":0}`)
		return
	}

	ps.CreateChannel(id)

	for {
		c := make(chan string)
		ps.AddListener(id, c)
		select {
		case <-cn.CloseNotify():
			log.Println("done: closed connection")
			ps.RespondOK(rw)
			// fmt.Fprintf(rw, `{"success":1}`)
			return
		case msg := <-c:
			msg, err := ps.OnMessage(id, msg)
			if err != nil {
				fmt.Fprintf(rw, `{"success":0}`)
				return
			}
			ps.WriteData(rw, msg)
			f.Flush()
		}
	}
}

func (ps *Server) ping(rw http.ResponseWriter, req *http.Request) {
	id, msg, ok := parsePing(req)
	if !ok {
		fmt.Fprintf(rw, `{"success":0}`)
		return
	}

	id, msg, err := ps.OnPing(id, msg)
	if err != nil {
		fmt.Fprintf(rw, `{"success":0}`)
		return
	}

	if ps.PingChannel(id, msg) != 0 {
		// fmt.Fprintf(rw, `{"success":1}`)
		ps.RespondOK(rw)
	} else {
		ps.RespondFail(rw)
		// fmt.Fprintf(rw, `{"success":0}`)
	}
}

func checkStreamable(rw http.ResponseWriter) (http.Flusher, http.CloseNotifier) {
	f, ok := rw.(http.Flusher)
	if !ok {
		http.Error(rw, "cannot stream", http.StatusInternalServerError)
		return nil, nil
	}

	cn, ok := rw.(http.CloseNotifier)
	if !ok {
		http.Error(rw, "cannot stream", http.StatusInternalServerError)
		return nil, nil
	}
	return f, cn
}

func setHeaders(rw http.ResponseWriter) {
	rw.Header().Set("Content-Type", "text/event-stream")
	rw.Header().Set("Cache-Control", "no-cache")
	rw.Header().Set("Connection", "keep-alive")
	rw.Header().Set("Access-Control-Allow-Origin", "*")
}

func parsePing(req *http.Request) (id, msg string, ok bool) {
	err := req.ParseForm()
	if err != nil {
		return "", "", false
	}

	m := req.PostForm
	id_arr, ok1 := m["id"]
	msg_arr, ok2 := m["msg"]

	if !ok1 || !ok2 {
		return "", "", false
	}
	return id_arr[0], msg_arr[0], true
}

func parseListen(req *http.Request) (id string, ok bool) {
	err := req.ParseForm()
	if err != nil {
		return "", false
	}

	m := req.PostForm
	if err != nil {
		return "", false
	}
	if id_arr, ok := m["id"]; ok {
		return id_arr[0], true
	}
	return "", false
}
