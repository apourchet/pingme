package ping

import (
	"bufio"
	b64 "encoding/base64"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
)

type Client struct {
	Host string
	Port string
}

func (c *Client) HostName() string {
	return c.Host + ":" + c.Port
}

func (c *Client) Ping(id, msg string) error {
	v := url.Values{"id": {id}, "msg": {msg}}

	resp, err := http.PostForm("http://"+c.Host+":"+c.Port+"/ping", v)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	_, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	// TODO parse body and retrieve "success" field
	return nil
}

func (c *Client) Listen(id string, out func(string) bool) error {
	v := url.Values{"id": {id}}

	resp, err := http.PostForm("http://"+c.HostName()+"/listen", v)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	scanner := bufio.NewScanner(resp.Body)
	for scanner.Scan() {
		msg := parseEvent(scanner.Text())
		if len(msg) == 0 {
			continue
		}
		msgBytes, err := b64.StdEncoding.DecodeString(msg)
		if err != nil {
			continue
		}
		if !out(string(msgBytes)) {
			break
		}
	}
	return nil
}

func parseEvent(evt string) string {
	if !strings.Contains(evt, "data: ") {
		return ""
	}
	msg := strings.Split(evt, "data: ")[1]
	return strings.Trim(msg, "\n \n")
}
