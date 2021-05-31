package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"n_communication/controller"

	"n_communication/entity"

	"github.com/go-chi/chi/v5"
)

// HealthHandler represents interface for health of a service
type CommHandler interface {
	Send(w http.ResponseWriter, r *http.Request)
	NewCommRouter() http.Handler
}

type commHandler struct {
	Coordinator controller.Coordinator
}

// NewHealthHandler creates new object of HealthHandler
func NewCommHandler() CommHandler {
	c := controller.New()
	return &commHandler{Coordinator: c}
}

// NewHealthRouter constructs new router for health endpoints
func (ch *commHandler) NewCommRouter() http.Handler {
	r := chi.NewRouter()

	r.Post("/", ch.Send)
	return r
}

func (ch *commHandler) Send(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	defer r.Body.Close()

	var sendRequest entity.SendRequest
	err := decoder.Decode(&sendRequest)

	if err != nil {
		res, _ := entity.NewErrorJSON("invalid send request " + err.Error())
		w.Write(res)
		return
	}

	status, err := ch.Coordinator.Send(
		sendRequest.Channel,
		sendRequest.To,
		sendRequest.From,
		sendRequest.Payload,
		sendRequest.Title,
	)

	if err != nil {
		res, _ := entity.NewErrorJSON(err.Error())
		w.Write(res)
	}

	e := entity.SuccessResponse{Status: strconv.FormatBool(status)}
	res, _ := json.Marshal(e)
	w.Write(res)
}
