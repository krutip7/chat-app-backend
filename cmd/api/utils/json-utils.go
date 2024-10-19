package utils

import (
	"encoding/json"
	"errors"
	"io"
	"log"
	"net/http"
)

type JSONResponse struct {
	Success bool        `json:"success"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

func writeJSONResponse(response http.ResponseWriter, payload interface{}, status int) {

	jsonResponse, err := json.Marshal(payload)
	if err != nil {
		status = http.StatusInternalServerError
		jsonResponse = make([]byte, 0)
	}

	response.Header().Add("Content-Type", "application/json")
	response.WriteHeader(status)

	_, err = response.Write(jsonResponse)
	if err != nil {
		log.Println("Error occured while writing response: ", err)
	}
}

func WriteJSONResponse(response http.ResponseWriter, data interface{}, headers ...http.Header) {

	payload := JSONResponse{
		Success: true,
		Message: "success",
		Data:    data,
	}

	if len(headers) > 0 {
		for key, val := range headers[0] {
			response.Header()[key] = val
		}
	}

	writeJSONResponse(response, payload, http.StatusOK)
}

func ReadJSONRequest(w http.ResponseWriter, request *http.Request, data interface{}) error {
	request.Body = http.MaxBytesReader(w, request.Body, int64(1024*1024))

	decoder := json.NewDecoder(request.Body)

	decoder.DisallowUnknownFields()

	err := decoder.Decode(data)
	if err != nil {
		return err
	}

	err = decoder.Decode(&struct{}{})
	if err != io.EOF {
		return errors.New("request body must contain only single JSON entity")
	}

	return nil
}

func WriteJSONErrorResponse(response http.ResponseWriter, err error, status ...int) {
	statusCode := http.StatusBadRequest
	if len(status) > 0 {
		statusCode = status[0]
	}

	payload := JSONResponse{
		Success: false,
		Message: err.Error(),
	}

	writeJSONResponse(response, payload, statusCode)
}
