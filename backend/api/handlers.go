package api

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
)

type Handler interface {
	turnOn(w http.ResponseWriter, r *http.Request)
	turnOff(w http.ResponseWriter, r *http.Request)
	setAlarm(w http.ResponseWriter, r *http.Request)
	cancelAlarm(w http.ResponseWriter, r *http.Request)
}

type handler struct {
	service Service
}

func NewRouter(h Handler) http.Handler {
	router := mux.NewRouter()
	router.Path("/api/on").Methods("GET").HandlerFunc(h.turnOn)
	router.Path("/api/off").Methods("GET").HandlerFunc(h.turnOff)
	router.Path("/api/set").Methods("POST").HandlerFunc(h.setAlarm)
	router.Path("/api/clear").Methods("DELETE").HandlerFunc(h.cancelAlarm)
	return router
}

func NewHandler(s Service) Handler {
	return &handler{
		service: s,
	}
}

func (h *handler) turnOn(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")

	err := h.service.turnOn()
	if err != nil {
		w.WriteHeader(500)
		return
	}
	w.WriteHeader(200)
}

func (h *handler) turnOff(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")

	err := h.service.turnOff()
	if err != nil {
		w.WriteHeader(500)
		return
	}
	w.WriteHeader(200)
}

type setReq struct {
	Hour *int `json:"hour,omitempty"`
	Min  *int `json:"min,omitempty"`
}

func (h *handler) setAlarm(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")

	req := setReq{}
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		w.WriteHeader(400)
		return
	}
	if req.Hour == nil || req.Min == nil {
		w.WriteHeader(400)
		return
	}
	err = h.service.setAlarm(*req.Hour, *req.Min)
	if err != nil {
		w.WriteHeader(500)
		return
	}
	w.WriteHeader(200)
}

func (h *handler) cancelAlarm(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")

	err := h.service.cancelAlarm()
	if err != nil {
		w.WriteHeader(500)
	}
	w.WriteHeader(200)
}
