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
	getAlarm(w http.ResponseWriter, r *http.Request)
}

type handler struct {
	service Service
}

func NewRouter(h Handler) http.Handler {
	router := mux.NewRouter()
	router.Path("/api/on").Methods("GET").HandlerFunc(h.turnOn)   // ideally a PATCH but problem with cors unsafe methods
	router.Path("/api/off").Methods("GET").HandlerFunc(h.turnOff) // ideally a PATCH but problem with cors unsafe methods
	router.Path("/api/set").Methods("POST").HandlerFunc(h.setAlarm)
	router.Path("/api/clear").Methods("GET").HandlerFunc(h.cancelAlarm) // ideally a DELETE but problem with cors unsafe methods
	router.Path("/api/alarm").Methods("GET").HandlerFunc(h.getAlarm)
	router.Use(CorsMiddleware)
	return router
}

func NewHandler(s Service) Handler {
	return &handler{
		service: s,
	}
}

func CorsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		next.ServeHTTP(w, r)
	})
}

func (h *handler) getAlarm(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	resp, err := h.service.getAlarm()
	if err != nil {
		w.WriteHeader(404)
		return
	}
	w.WriteHeader(200)
	json.NewEncoder(w).Encode(resp)
}

func (h *handler) turnOn(w http.ResponseWriter, r *http.Request) {

	err := h.service.turnOn()
	if err != nil {
		w.WriteHeader(500)
		return
	}
	w.WriteHeader(200)
}

func (h *handler) turnOff(w http.ResponseWriter, r *http.Request) {

	err := h.service.turnOff()
	if err != nil {
		w.WriteHeader(500)
		return
	}
	w.WriteHeader(200)
}

func (h *handler) setAlarm(w http.ResponseWriter, r *http.Request) {

	req := timeModel{}
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

	err := h.service.cancelAlarm()
	if err != nil {
		w.WriteHeader(500)
	}
	w.WriteHeader(200)
}
