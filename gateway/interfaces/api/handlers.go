package api

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"html/template"
	"io"
	"net/http"

	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/yaroyan/ms/gateway/infrastructure/messaging"
)

type Handlers struct {
	Client     *http.Client
	Connection *amqp.Connection
}

type Action string

const (
	Authn = Action("Authn")
	Log   = Action("Log")
	Mail  = Action("Mail")
)

const (
	AuthnURI  = "http://authn:80"
	LoggerURI = "http://logger:80"
	MailerURI = "http://mailer:80"
)

type RequestPayload struct {
	Action Action `json:"action"`
	Body   any    `json:"body,omitempty"`
}

func (h *Handlers) Broker(w http.ResponseWriter, r *http.Request) {
	payload := JsonResponse{
		Error:   false,
		Message: "Hit the broker",
	}

	_ = WriteJSON(w, http.StatusOK, payload)
}

func (h *Handlers) HandleSubmission(w http.ResponseWriter, r *http.Request) {
	var p RequestPayload

	err := ReadJSON(w, r, &p)
	if err != nil {
		ErrorJSON(w, err)
		return
	}

	switch p.Action {
	case Authn:
		h.authenticate(w, p.Body)
	case Log:
		// h.log(w, payload.Body)
		h.logViaAmqp(w, p.Body)
	case Mail:
		h.sendMail(w, p.Body)
	default:
		ErrorJSON(w, errors.New("unknown action"))
	}
}

func (h *Handlers) authenticate(w http.ResponseWriter, p any) {
	jsonData, _ := json.MarshalIndent(p, "", "\t")

	request, err := http.NewRequest("POST", fmt.Sprintf("%s/authenticate", AuthnURI), bytes.NewBuffer(jsonData))
	if err != nil {
		ErrorJSON(w, err)
		return
	}

	responseFromService, err := h.Client.Do(request)
	if err != nil {
		ErrorJSON(w, err)
		return
	}
	defer responseFromService.Body.Close()

	var jsonFromService JsonResponse
	err = json.NewDecoder(responseFromService.Body).Decode(&jsonFromService)
	if err != nil {
		ErrorJSON(w, err)
		return
	}
	switch responseFromService.StatusCode {
	case http.StatusUnauthorized:
		ErrorJSON(w, errors.New(jsonFromService.Message))
		return
	case http.StatusAccepted:
		// OK
	default:
		ErrorJSON(w, errors.New("error calling authn service"))
		return
	}
	if jsonFromService.Error {
		ErrorJSON(w, err, http.StatusUnauthorized)
		return
	}

	WriteJSON(w, http.StatusAccepted, JsonResponse{
		Error:   false,
		Message: "Authenticated.",
		Data:    jsonFromService.Data,
	})
}

func (h *Handlers) log(w http.ResponseWriter, b any) {
	jsonData, _ := json.MarshalIndent(b, "", "\t")

	request, err := http.NewRequest(
		"POST",
		fmt.Sprintf("%s/log", LoggerURI),
		bytes.NewBuffer(jsonData))

	if err != nil {
		ErrorJSON(w, err)
		return
	}

	request.Header.Set("Content-Type", "application/json")

	responseFromService, err := h.Client.Do(request)
	if err != nil {
		ErrorJSON(w, err)
		return
	}
	defer func() {
		io.Copy(io.Discard, responseFromService.Body)
		responseFromService.Body.Close()
	}()

	if responseFromService.StatusCode != http.StatusAccepted {
		ErrorJSON(w, errors.New("error calling logger service"))
		return
	}

	WriteJSON(w, http.StatusAccepted, JsonResponse{
		Error:   false,
		Message: "logged",
	})
}

func (h *Handlers) sendMail(w http.ResponseWriter, b any) {
	jsonData, _ := json.MarshalIndent(b, "", "\t")

	request, err := http.NewRequest(
		"POST",
		fmt.Sprintf("%s/send", MailerURI),
		bytes.NewBuffer(jsonData))

	if err != nil {
		ErrorJSON(w, err)
		return
	}

	request.Header.Set("Content-Type", "application/json")

	responseFromService, err := h.Client.Do(request)
	if err != nil {
		ErrorJSON(w, err)
		return
	}
	defer func() {
		io.Copy(io.Discard, responseFromService.Body)
		responseFromService.Body.Close()
	}()

	if responseFromService.StatusCode != http.StatusAccepted {
		ErrorJSON(w, errors.New("error calling mail service"))
		return
	}

	WriteJSON(w, http.StatusAccepted, JsonResponse{
		Error:   false,
		Message: "message sent",
	})
}

func (h *Handlers) logViaAmqp(w http.ResponseWriter, b any) {
	err := h.pushToQueue(b)
	if err != nil {
		ErrorJSON(w, err)
		return
	}

	WriteJSON(w, http.StatusAccepted, JsonResponse{
		Error:   false,
		Message: "logged via Amqp",
	})
}

func (h *Handlers) pushToQueue(b any) error {
	emitter, err := messaging.NewEventEmitter(h.Connection)
	if err != nil {
		return err
	}
	j, _ := json.MarshalIndent(&b, "", "\t")
	err = emitter.Push(string(j), "log.INFO")
	if err != nil {
		return err
	}
	return nil
}

func (h *Handlers) Render(w http.ResponseWriter, r *http.Request) {

	partials := []string{
		"./web/templates/base.layout.gohtml",
		"./web/templates/header.partial.gohtml",
		"./web/templates/footer.partial.gohtml",
	}

	var templateSlice []string
	templateSlice = append(templateSlice, fmt.Sprintf("./web/templates/%s", "test.page.gohtml"))

	for _, x := range partials {
		templateSlice = append(templateSlice, x)
	}

	tmpl, err := template.ParseFiles(templateSlice...)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if err := tmpl.Execute(w, nil); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
