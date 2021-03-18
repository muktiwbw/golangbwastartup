package helpers

import (
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

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
	return Response{
		Meta: Meta{
			Message: message,
			Code:    code,
			Status:  status,
		},
		Data: data,
	}
}

func GetValidationErrors(e error) interface{} {
	var errors []string

	for _, err := range e.(validator.ValidationErrors) {
		errors = append(errors, err.Error())
	}

	return gin.H{"errors": errors}
}
