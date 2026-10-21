package student

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/asthap08/Golang_notes/internal/storage"
	"github.com/asthap08/Golang_notes/internal/types"
	"github.com/asthap08/Golang_notes/internal/utils/response"

	"github.com/go-playground/locales/id"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	idtranslations "github.com/go-playground/validator/v10/translations/id"
)

var (
	validate   = validator.New()
	translator ut.Translator
)

func init() {
	idLocale := id.New()
	uni := ut.New(idLocale, idLocale)

	var found bool
	translator, found = uni.GetTranslator("id")
	if !found {
		panic("Indonesian translator not found")
	}

	if err := idtranslations.RegisterDefaultTranslations(validate, translator); err != nil {
		panic(err)
	}
}

func New(storage storage.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		slog.Info("Creating a student")

		var student types.Student

		err := json.NewDecoder(r.Body).Decode(&student)
		if errors.Is(err, io.EOF) {
			response.WriteJSON(
				w,
				http.StatusBadRequest,
				response.GeneralError(fmt.Errorf("empty body")),
			)
			return
		}

		if err != nil {
			response.WriteJSON(
				w,
				http.StatusBadRequest,
				response.GeneralError(err),
			)
			return
		}

		// Validate request
		if err := validate.Struct(student); err != nil {
			var validationErrors validator.ValidationErrors

			if errors.As(err, &validationErrors) {
				response.WriteJSON(
					w,
					http.StatusBadRequest,
					response.ValidatorError(validationErrors),
				)
				return
			}

			response.WriteJSON(
				w,
				http.StatusBadRequest,
				response.GeneralError(err),
			)
			return
		}

		lastID, err := storage.CreateStudent(
			student.Name,
			student.Email,
			student.Age,
		)
		if err != nil {
			slog.Error("error creating student", slog.String("error", err.Error()))
			response.WriteJSON(
				w,
				http.StatusInternalServerError,
				response.GeneralError(err),
			)
			return
		}

		slog.Info(
			"student created successfully",
			slog.String("userId", fmt.Sprint(lastID)),
		)

		response.WriteJSON(
			w,
			http.StatusCreated,
			map[string]int64{"id": lastID},
		)
	}
}

func GetById(storage storage.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := r.PathValue("id")

		slog.Info("getting a student", slog.String("id", id))

		intID, err := strconv.ParseInt(id, 10, 64)
		if err != nil {
			response.WriteJSON(
				w,
				http.StatusBadRequest,
				response.GeneralError(err),
			)
			return
		}

		student, err := storage.GetStudentByID(intID)
		if err != nil {
			slog.Error(
				"error getting student",
				slog.String("id", id),
				slog.String("error", err.Error()),
			)

			response.WriteJSON(
				w,
				http.StatusInternalServerError,
				response.GeneralError(err),
			)
			return
		}

		response.WriteJSON(w, http.StatusOK, student)
	}
}

func GetList(storage storage.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		slog.Info("getting all students")

		students, err := storage.GetStudents()
		if err != nil {
			slog.Error("error getting students", slog.String("error", err.Error()))
			response.WriteJSON(
				w,
				http.StatusInternalServerError,
				response.GeneralError(err),
			)
			return
		}

		response.WriteJSON(w, http.StatusOK, students)
	}
}
