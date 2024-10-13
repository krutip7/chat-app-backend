package utils

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func WriteJSONResponse(response http.ResponseWriter, payload interface{}) {

	jsonResponse, err := json.Marshal(payload)
	if err != nil {
		fmt.Println(err)
	}

	response.Header().Add("Content-Type", "application/json")

	response.Write(jsonResponse)
}


