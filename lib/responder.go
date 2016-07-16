package ping

import (
	"encoding/json"
	"io"
)

// TODO Ugly as hell
type Responder int

type Response struct {
	Success      int    `json:"success"`
	ErrorCode    int    `json:"error_code,omitempty"`
	ErrorMessage string `json:"error_message,omitempty"`
}

func (r *Response) Write(w io.Writer) error {
	bytes, err := json.MarshalIndent(r, "", "  ")
	if err != nil {
		return err
	}
	w.Write(bytes)
	return nil
}

func (r *Responder) RespondOK(w io.Writer) error {
	res := new(Response)
	res.Success = 1
	return res.Write(w)
}

func (r *Responder) RespondFail(w io.Writer) error {
	res := new(Response)
	res.Success = 0
	return res.Write(w)
}

func (r *Responder) FromJson(s string) (*Response, error) {
	res := new(Response)
	err := json.Unmarshal([]byte(s), res)
	return res, err
}
