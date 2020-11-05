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
	Code  string
	Email string
}

type Coupons struct {
	Coupon []Coupon
}

func (c Coupons) Check(code string) (string, string) {
	for _, item := range c.Coupon {
		if code == item.Code {
			return "valid", item.Email
		}
	}
	return "", "invalid"
}

type Result struct {
	Status string
}

var coupons Coupons

func main() {
	coupon := Coupon{
		Code:  "abc",
		Email: "schmidt_tech@outlook.com",
	}

	coupons.Coupon = append(coupons.Coupon, coupon)

	http.HandleFunc("/", home)
	http.ListenAndServe(":9092", nil)
}

func home(w http.ResponseWriter, r *http.Request) {
	coupon := r.PostFormValue("coupon")
	valid, email := coupons.Check(coupon)

	if valid == "valid" {
		resultMail := makeHttpCall("http://localhost:9093", email, coupon)
		if resultMail.Status != "mail send" {
			fmt.Println("Email não enviado!")
			fmt.Println(resultMail.Status)
		} else {
			fmt.Println("Email enviado")
		}
	}

	result := Result{Status: valid}

	jsonResult, err := json.Marshal(result)

	if err != nil {
		log.Fatal("Error converting json")
	}

	fmt.Fprintf(w, string(jsonResult))
}

func makeHttpCall(urlMicrosservice string, email string, coupon string) Result {

	values := url.Values{}
	values.Add("email", email)
	values.Add("coupon", coupon)

	retryClient := retryablehttp.NewClient()

	retryClient.RetryMax = 5

	res, err := retryClient.PostForm(urlMicrosservice, values)

	if err != nil {
		result := Result{Status: "Servidor fora do ar!"}
		return result
	}

	defer res.Body.Close()

	data, err := ioutil.ReadAll(res.Body)

	if err != nil {
		log.Fatal("Error processing result")
	}

	result := Result{}
	json.Unmarshal(data, &result)

	return result
}
