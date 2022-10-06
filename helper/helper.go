package helper

import "github.com/go-playground/validator/v10"

type Meta struct {
	Message string `json:"message"`
	Code    int    `json:"code"`
	Status  string `json:"status"`
}

type Response struct {
	Meta Meta        `json:"meta"`
	Data interface{} `json:"data"`
}

func APIResponse(message string, code int, status string, data interface{}) Response {
	//	create meta
	meta := Meta{
		Message: message,
		Code:    code,
		Status:  status,
	}

	//	create response
	jsonResponse := Response{
		Meta: meta,
		Data: data,
	}

	return jsonResponse
}

func FormatValidationError(err error) []string {
	//formatting errors validation into array of errors
	var errors []string

	//iterate for each errors and cast it to validator.ValidationErrors
	for _, e := range err.(validator.ValidationErrors) {
		errors = append(errors, e.Error())
	}

	return errors
}
