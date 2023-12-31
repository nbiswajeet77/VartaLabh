package model

import (
	"fmt"
	"net/http"

	jsoniter "github.com/json-iterator/go"
)

var jsonIter = jsoniter.ConfigDefault

func WriteOutput(w http.ResponseWriter, data interface{}, status int, err error) {
	response := Response{
		StatusCode: status,
		Data:       data,
	}
	w.WriteHeader(status)
	if err != nil {
		byteData, _ := jsonIter.Marshal(err)
		w.Write(byteData)
		return
	}
	byteData, marshalErr := jsonIter.Marshal(response)
	if marshalErr != nil {
		fmt.Println(marshalErr)
	}
	w.Write(byteData)
}
