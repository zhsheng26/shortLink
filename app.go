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
	a.Router.HandleFunc("api/shorten", a.createShortLink).Methods("post")
	a.Router.HandleFunc("api/info", a.getShortLinkInfo).Methods("GET")
	a.Router.HandleFunc("/{shortLink:[a-zA-z0-9]{1,11}", a.redirect).Methods("get")
}

func (a *App) createShortLink(w http.ResponseWriter, r *http.Request) {
	var req shortenReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return
	}
	if err := validator.Validate(req); err != nil {
		return
	}
	defer r.Body.Close()
	fmt.Printf("%v\n", req)
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
	log.Fatal(http.ListenAndServe(addr, a.Router))

}
