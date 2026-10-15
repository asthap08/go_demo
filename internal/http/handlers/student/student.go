package student

import (
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"
	"io"
	"fmt"
	

	"github.com/asthap08/Golang_notes/internal/types"
	"github.com/asthap08/Golang_notes/internal/utils/response"
	"github.com/go-playground/validator/v10"

)

func New() http.HandlerFunc{
	return func(w http.ResponseWriter, r *http.Request) {
		slog.Info("Creating a student")


		var student types.Student


	err := json.NewDecoder(r.Body).Decode(&student)
	if errors.Is(err, io.EOF) {
		response.WriteJSON(w, http.StatusBadRequest, response.GeneralError(fmt.Errorf("empty body")))
		return
	}

	if err != nil {
		response.WriteJSON(w, http.StatusBadRequest, response.GeneralError(err))
		return
	}

	//validate request
	if err := validator.New().Struct(student); err != nil {

		validateErrs := err.(validator.ValidationErrors)
		response.WriteJSON(w, http.StatusBadRequest, response.ValidatorError(validateErrs))
		return
	}

	

	response.WriteJSON(w, http.StatusCreated, map[string]string{"success":"OK"})
}	
}