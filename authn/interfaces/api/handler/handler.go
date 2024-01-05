package handler

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"

	repository "github.com/yaroyan/ms/authn/infrastructure/persistence/repository/postgres"
	"github.com/yaroyan/ms/authn/interfaces/api/response"
)

type AuthenticationHandler struct {
	Repository *repository.UserRepository
	Client     *http.Client
}

const LoggerURI = "http://logger:80"

// Authenticate user
func (h *AuthenticationHandler) Authenticate(w http.ResponseWriter, r *http.Request) {
	var payload struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	err := response.ReadJSON(w, r, &payload)
	if err != nil {
		response.ErrorJSON(w, err, http.StatusBadRequest)
		return
	}

	user, err := h.Repository.FindByEmail(payload.Email)
	if err != nil {
		response.ErrorJSON(w, errors.New("invalid credentials"), http.StatusUnauthorized)
		return
	}

	valid, err := user.Password.Match(payload.Password)
	isValidPassword := err == nil && valid
	if !isValidPassword {
		response.ErrorJSON(w, errors.New("invalid credentials"), http.StatusUnauthorized)
		return
	}

	// log authentication
	err = h.log("authentication", fmt.Sprintf("%s logged in", user.Email))
	if err != nil {
		response.ErrorJSON(w, err)
		return
	}

	res := response.JsonResponse{
		Error:   false,
		Message: fmt.Sprintf("Logged in user %s", user.Email),
		Data:    response.NewResponse(user),
	}

	response.WriteJSON(w, http.StatusAccepted, res)
}

func (h *AuthenticationHandler) log(name, data string) error {
	var entry struct {
		Name string `json:"name"`
		Data string `json:"data"`
	}
	entry.Name = name
	entry.Data = data

	jsonData, _ := json.MarshalIndent(entry, "", "\t")
	request, err := http.NewRequest("POST", fmt.Sprintf("%s/log", LoggerURI), bytes.NewBuffer(jsonData))
	if err != nil {
		return err
	}
	res, err := h.Client.Do(request)
	if err != nil {
		return err
	}
	defer func() {
		io.Copy(io.Discard, res.Body)
		res.Body.Close()
	}()
	return nil
}
