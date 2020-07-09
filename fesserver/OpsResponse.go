package fesserver

import (
	"log"
	"net/http"

	"github.com/bitly/go-simplejson"
)

type OperationResponse struct {
	Result  bool   `json:"result"`
	Message string `json:"message"`
}

func (opsres *OperationResponse) SetToSimpleJSON() *simplejson.Json {
	json := simplejson.New()
	json.Set("result", opsres.Result)
	json.Set("message", opsres.Message)

	return json
}

func WriteToResponseStream(w http.ResponseWriter, jsonData *simplejson.Json, contentType string) {
	payload, err := jsonData.MarshalJSON()
	if err != nil {
		log.Println(err)
	}

	w.Header().Set("Content-Type", contentType)
	w.Write(payload)
}
