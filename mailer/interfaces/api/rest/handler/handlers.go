package handler

import (
	"net/http"

	"github.com/yaroyan/gms/mail/application/usecase"
	"github.com/yaroyan/gms/mail/domain/model"
	"github.com/yaroyan/gms/mail/interfaces/api/rest/response"
)

type MailHandler struct {
	Usecase *usecase.Usecase
}

func (h *MailHandler) SendMail(w http.ResponseWriter, r *http.Request) {
	type mailMessage struct {
		From    string `json:"from"`
		To      string `json:"to"`
		Subject string `json:"subject"`
		Message string `json:"message"`
	}

	var requestPayload mailMessage

	err := response.ReadJSON(w, r, &requestPayload)
	if err != nil {
		response.ErrorJSON(w, err)
		return
	}

	msg := model.Message{
		From:    requestPayload.From,
		To:      requestPayload.To,
		Subject: requestPayload.Subject,
		Data:    requestPayload.Message,
	}

	err = h.Usecase.SendSMTPMessage(msg)
	if err != nil {
		response.ErrorJSON(w, err)
		return
	}

	payload := response.JsonResponse{
		Error:   false,
		Message: "send to " + requestPayload.To,
	}

	response.WriteJSON(w, http.StatusAccepted, payload)
}
