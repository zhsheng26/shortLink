package main

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/justinas/alice"
	"gopkg.in/validator.v2"
	"log"
	"net/http"
)

type App struct {
	Router     *mux.Router
	Middleware *Middleware
	Config     *Env
}

type shortenReq struct {
	Url                 string `json:"url" validate:"nonzero"`
	ExpirationInMinutes int64  `json:"expiration_in_minutes" validate:"min=0"`
}

type shortLinkResp struct {
	ShortLink string `json:"short_link"`
}

func (a *App) Initialize(env *Env) {
	//log flag 的含义
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	a.Config = env
	a.Router = mux.NewRouter()
	a.Middleware = &Middleware{}
	a.initializeRoutes()
}

func (a *App) Run(addr string) {
	log.Fatal(http.ListenAndServe(addr, a.Router))
}

func (a *App) initializeRoutes() {
	//a.Router.HandleFunc("/api/shorten", a.createShortLink).Methods("POST")
	//a.Router.HandleFunc("/api/info", a.getShortLinkInfo).Methods("GET")
	//a.Router.HandleFunc("/{shortLink:[a-zA-z0-9]{1,11}", a.redirect).Methods("GET")
	m := alice.New(a.Middleware.LoggingHandler, a.Middleware.RecoverHandler)
	a.Router.Handle("/api/shorten", m.ThenFunc(a.createShortLink)).Methods("POST")
	a.Router.Handle("/api/info", m.ThenFunc(a.getShortLinkInfo)).Methods("GET")
	a.Router.Handle("/{shortLink:[a-zA-z0-9]{1,11}", m.ThenFunc(a.redirect)).Methods("GET")
}

// 生成短地址
func (a *App) createShortLink(w http.ResponseWriter, r *http.Request) {
	var req shortenReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		responseWithError(w, NewBadReqErr(fmt.Errorf("%s", "happen err when parsing json from body")), nil)
		return
	}
	if err := validator.Validate(req); err != nil {
		responseWithError(w, NewBadReqErr(fmt.Errorf("validate parameters failed : %+v", req)), nil)
		return
	}
	defer r.Body.Close()
	shorten, err := a.Config.S.Shorten(req.Url, req.ExpirationInMinutes)
	if err != nil {
		responseWithError(w, err, nil)
	} else {
		responseWithJson(w, http.StatusCreated, shortLinkResp{ShortLink: shorten})
	}
}

// 短地址解析
func (a *App) getShortLinkInfo(w http.ResponseWriter, r *http.Request) {
	values := r.URL.Query()
	sl := values.Get("shortLink")
	info, err := a.Config.S.ShortLinkInfo(sl)
	if err != nil {
		responseWithError(w, err, nil)
	} else {
		responseWithJson(w, http.StatusOK, info)
	}
}

//访问，重定向 302
func (a *App) redirect(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	url, err := a.Config.S.UnShorten(vars["shortLink"])
	if err != nil {
		responseWithError(w, err, nil)
	} else {
		http.Redirect(w, r, url, http.StatusTemporaryRedirect)
	}
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
