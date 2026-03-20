package utils

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/go-playground/validator/v10"
)

var Validate = validator.New()

type ValidationError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

func DecodeAndValidate[T any](r *http.Request) (T, error) {
	var body T

	if r.Body == nil {
		return body, errors.New("o corpo da requisição está vazio")
	}
	defer r.Body.Close()

	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()

	if err := decoder.Decode(&body); err != nil {
		return body, fmt.Errorf("formato JSON inválido: %v", err)
	}

	if err := Validate.Struct(body); err != nil {
		return body, err
	}

	return body, nil
}

func FormatValidationErrors(err error) []ValidationError {
	var errs []ValidationError

	var validationErrors validator.ValidationErrors
	if errors.As(err, &validationErrors) {
		for _, e := range validationErrors {
			var msg string
			switch e.Tag() {
			case "required":
				msg = "Este campo é obrigatório"
			case "email":
				msg = "Formato de e-mail inválido"
			case "min":
				msg = fmt.Sprintf("O valor deve ser no mínimo %s", e.Param())
			case "max":
				msg = fmt.Sprintf("O valor deve ser no máximo %s", e.Param())
			default:
				msg = fmt.Sprintf("Falha na validação: %s", e.Tag())
			}

			errs = append(errs, ValidationError{
				Field:   e.Field(),
				Message: msg,
			})
		}
		return errs
	}

	return append(errs, ValidationError{
		Field:   "geral",
		Message: err.Error(),
	})
}
