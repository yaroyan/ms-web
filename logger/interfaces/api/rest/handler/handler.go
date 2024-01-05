package handler

import (
	"net/http"

	ml "github.com/yaroyan/ms/logger/domain/model/log"
	"github.com/yaroyan/ms/logger/infrastructure/persistence"
	"github.com/yaroyan/ms/logger/interfaces/api/rest/response"
)

type JSONPayload struct {
	Name string `json:"name"`
	Data string `json:"data"`
}

type LogHandler struct {
	Repository persistence.LogRepository
}

func (app *LogHandler) WriteLog(w http.ResponseWriter, r *http.Request) {
	var requestPayload JSONPayload
	_ = response.ReadJSON(w, r, &requestPayload)

	event := ml.Log{
		Name: requestPayload.Name,
		Data: requestPayload.Data,
	}

	err := app.Repository.Insert(event)

	if err != nil {
		response.ErrorJSON(w, err)
		return
	}

	resp := response.JsonResponse{
		Error:   false,
		Message: "logged.",
	}

	response.WriteJSON(w, http.StatusAccepted, resp)
}
