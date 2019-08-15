package main

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"gopkg.in/validator.v2"
	"log"
	"net/http"
)

type App struct {
	Router *mux.Router
}

type shortenReq struct {
	Url                 string `json:"url" validate:"nonzero"`
	ExpirationInMinutes int64  `json:"expiration_in_minutes" validate:"min=0"`
}

type shortLinkResp struct {
	ShortLink string `json:"short_link"`
}

func (a *App) Initialize() {
	//log flag 的含义
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	a.Router = mux.NewRouter()
}

func (a *App) initializeRoutes() {
	a.Router.HandleFunc("/api/shorten", a.createShortLink).Methods("POST")
	a.Router.HandleFunc("/api/info", a.getShortLinkInfo).Methods("GET")
	a.Router.HandleFunc("/{shortLink:[a-zA-z0-9]{1,11}", a.redirect).Methods("GET")
}

func (a *App) createShortLink(w http.ResponseWriter, r *http.Request) {
	var req shortenReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		responseWithError(w, ReqError{http.StatusBadRequest, fmt.Errorf("%s", "happen err when parsing json from body")}, nil)
		return
	}
	if err := validator.Validate(req); err != nil {
		responseWithError(w, ReqError{http.StatusBadRequest, fmt.Errorf("validate parameters failed : %+v", req)}, nil)
		return
	}
	defer r.Body.Close()
	responseWithJson(w, 200, req)
}

func (a *App) getShortLinkInfo(w http.ResponseWriter, r *http.Request) {
	values := r.URL.Query()
	sl := values.Get("shortLink")
	fmt.Printf("%s\n", sl)
}
func (a *App) redirect(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	fmt.Printf("%s\n", vars["shortLink"])
}

func (a *App) Run(addr string) {
	a.Initialize()
	a.initializeRoutes()
	log.Fatal(http.ListenAndServe(addr, a.Router))

}

func responseWithError(w http.ResponseWriter, err error, payload interface{}) {
	switch e := err.(type) {
	case MiError:
		log.Printf("http %d - %s", e.Status(), e)
		resp, _ := json.Marshal(Response{Code: e.Status(), Message: e.Error(), Content: payload})
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write(resp)
	default:
		responseWithJson(w, http.StatusInternalServerError, payload)
	}
}

func responseWithJson(w http.ResponseWriter, status int, payload interface{}) {
	resp, _ := json.Marshal(Response{Code: status, Message: http.StatusText(status), Content: payload})
	w.Header().Set("Content-Type", "application/json")
	_, _ = w.Write(resp)
}

type Response struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Content interface{} `json:"content"`
}
