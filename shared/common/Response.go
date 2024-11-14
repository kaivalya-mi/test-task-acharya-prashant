package common

import (
	"encoding/json"
	"net/http"
	_const "test-task/shared/utils/const"
)

const (
	CodeSuccess    = 200
	CodeBadRequest = 400
)

type BaseSuccessResponse struct {
	Data interface{} `json:"data,omitempty"  structs:"data"`
	Meta Meta        `json:"meta,omitempty"  structs:"meta"`
}

type BaseErrorResponse struct {
	Message string `json:"message" structs:"message"`
	Code    int    `json:"code" structs:"code"`
}

type Meta struct {
	Message string `json:"message" structs:"message"`
	Code    int    `json:"code" structs:"code"`
}

type ResponseData struct {
	Data interface{} `json:"data"`
	Meta interface{} `json:"meta"`
}

func Respond(w http.ResponseWriter, status int, data map[string]interface{}) {
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

func MessageWithCode(status int, message string) map[string]interface{} {
	return map[string]interface{}{"res_code": status, "message": message}
}

func MessageWithoutCode(message string) map[string]interface{} {
	return map[string]interface{}{"message": message}
}

func ResponseSuccessWithCode(message string, data ...interface{}) map[string]interface{} {
	response := map[string]interface{}{}
	response["meta"] = MessageWithCode(http.StatusOK, message)
	if data != nil {
		response["data"] = data
	}
	return response
}

func ResponseSuccessWithArray(message string, data ...interface{}) map[string]interface{} {
	response := map[string]interface{}{}
	response["meta"] = MessageWithCode(http.StatusOK, message)
	if data != nil {
		response["data"] = data
	}
	return response
}

func ResponseSuccessWithObj(message string, data interface{}) map[string]interface{} {
	response := map[string]interface{}{}
	response["meta"] = MessageWithCode(http.StatusOK, message)
	if data != nil {
		response["data"] = data
	}
	return response
}

func ResponseErrorWithCode(status int, message string) map[string]interface{} {
	return MessageWithCode(status, message)
}

func ResponseErrorWithoutCode(message string) map[string]interface{} {
	return MessageWithoutCode(message)
}

func GetError(field, tag string) string {
	if tag == "required" {
		return field + " is " + tag
	} else if tag == "email" {
		return field + " not valid"
	}
	return ""
}

func GetHTTPStatusCode(resCode interface{}) int {
	if resCode != nil {
		if resCode == _const.ResCodeError {
			return http.StatusBadRequest
		}
	}
	return http.StatusOK
}
