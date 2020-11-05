package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"

	"github.com/hashicorp/go-retryablehttp"
)

type Coupon struct {
	Code string
}

type Coupons struct {
	Coupon []Coupon
}

func (c Coupons) Check(code string) string {
	for _, item := range c.Coupon {
		if code == item.Code {
			return "valid"
		}
	}
	return "invalid"
}

type Result struct {
	Status string
}

var coupons Coupons

func main() {
	coupon := Coupon{
		Code: "abc",
	}

	coupons.Coupon = append(coupons.Coupon, coupon)

	http.HandleFunc("/", process)
	http.ListenAndServe(":9092", nil)

}

func process(w http.ResponseWriter, r *http.Request) {
	cccvv := r.PostFormValue("cccvv")

	resultcodigoSeguranca := makeHttpCall("http://localhost:9093", cccvv)

	result := Result{Status: "decline"}

	if resultcodigoSeguranca.Status == "invalid" {
		result.Status = resultcodigoSeguranca.Status
		fmt.Println("deu ruim")
	}

	fmt.Println(resultcodigoSeguranca.Status)

	jsonData, err := json.Marshal(result)

	if err != nil {
		log.Fatal("Erro processing json")
	}

	fmt.Fprintf(w, string(jsonData))
}

func makeHttpCall(urlMicrosercice string, codigoSeguranca string) Result {
	values := url.Values{}
	values.Add("cccvv", codigoSeguranca)

	retryClient := retryablehttp.NewClient()
	retryClient.RetryMax = 5

	res, err := retryClient.PostForm(urlMicrosercice, values)

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

	fmt.Println(string(data))

	return result
}
