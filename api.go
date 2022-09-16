package api

import (
	"encoding/json"
	"github.com/moeshin/go-errs"
	"io"
	"log"
	"net/http"
)

const MimeJson = "application/json"

type Api struct {
	JOk        bool   `json:"ok"`
	JMsg       string `json:"msg"`
	JData      any    `json:"data"`
	Code       int    `json:"-"`
	IsBeautify bool   `json:"-"`
	request    *http.Request
}

func (a *Api) JsonMarshal(v any) ([]byte, error) {
	if a.IsBeautify {
		return json.MarshalIndent(v, "", "  ")
	}
	return json.Marshal(v)
}

func (a *Api) ParseJson(v any) error {
	data, err := io.ReadAll(a.request.Body)
	if err != nil {
		return err
	}
	return json.Unmarshal(data, v)
}

func (a *Api) Response(w http.ResponseWriter) {
	w.Header().Set("Content-Type", MimeJson)
	if a.Code != http.StatusOK {
		w.WriteHeader(a.Code)
		return
	}
	data, err := a.JsonMarshal(a)
	if errs.Print(err) {
		return
	}
	_, err = w.Write(data)
	errs.Print(err)
}

func (a *Api) Ok(data any) {
	a.JData = data
	a.JOk = true
}

func (a *Api) Msg(msg string) {
	a.JOk = false
	a.JMsg = msg
	log.Println(a.request.Method, a.request.URL.String(), msg)
}

func (a *Api) Err(err error) bool {
	b := err != nil
	if b {
		a.Msg(err.Error())
		errs.Print(err)
	}
	return b
}

func (a *Api) AddMsg(msg string) {
	if a.JMsg != "" {
		a.JMsg += "\n"
	}
	a.JMsg += msg
}

func (a *Api) Request() *http.Request {
	return a.request
}

func New(r *http.Request) *Api {
	return &Api{
		JOk:        false,
		JMsg:       "",
		JData:      nil,
		request:    r,
		Code:       http.StatusOK,
		IsBeautify: true,
	}
}
