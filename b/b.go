package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
)

type Result struct {
	Status string
}

func main() {
	http.HandleFunc("/", home)
	http.ListenAndServe(":9091", nil)
}

func home(w http.ResponseWriter, r *http.Request) {
	coupon := r.PostFormValue("coupon")
	ccNumber := r.PostFormValue("ccNumber")
	cccvv := r.PostFormValue("cccvv")

	resultCoupon := makeHttpCall("http://localhost:9092", coupon, cccvv)

	result := Result{Status: "decline"}

	if ccNumber == "1" {
		result.Status = "approved"
	}

	if resultCoupon.Status == "invalid" {
		result.Status = resultCoupon.Status
	}

	jsonData, err := json.Marshal(result)

	if err != nil {
		log.Fatal("Erro processing json")
	}

	fmt.Fprintf(w, string(jsonData))
}

func makeHttpCall(urlMicrosercice string, coupon string, cccvv string) Result {
	values := url.Values{}
	values.Add("counpon", coupon)
	values.Add("cccvv", cccvv)
	res, err := http.PostForm(urlMicrosercice, values)

	if err != nil {
		result := Result{Status: "Servidor fora do ar!"}
		return result
	}

	defer res.Body.Close()

	data, err := ioutil.ReadAll(res.Body)

	if err != nil {
		log.Fatal("Erro processing result")
	}

	result := Result{}

	json.Unmarshal(data, &result)

	return result
}
