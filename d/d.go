package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

type Result struct {
	Status string
}

func main() {
	http.HandleFunc("/", home)
	http.ListenAndServe(":9093", nil)
}

func home(w http.ResponseWriter, r *http.Request) {

	codigoSeguranca := r.PostFormValue("cccvv")

	result := Result{Status: "decline"}

	if codigoSeguranca == "456" {
		result.Status = "ok ok ok ratinho"
	}

	jsonData, err := json.Marshal(result.Status)

	if err != nil {
		log.Fatal("Erro processing json")
	}

	fmt.Fprintf(w, string(jsonData))
}
