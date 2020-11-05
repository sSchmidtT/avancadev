package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/smtp"
)

type Result struct {
	Status string
}

func main() {
	http.HandleFunc("/", home)
	http.ListenAndServe(":9093", nil)
}

func home(w http.ResponseWriter, r *http.Request) {
	email := r.PostFormValue("email")
	coupon := r.PostFormValue("coupon")
	result := sendMail(email, coupon)

	jsonData, err := json.Marshal(result)

	if err != nil {
		log.Fatal(err)
	}

	fmt.Fprintf(w, string(jsonData))
}

func sendMail(email string, coupon string) Result {
	from := "...@gmail.com"
	pass := "..."

	msg := "From: " + from + "\n" +
		"To: " + email + "\n" +
		"Subject: Cupom de desconto ultilizado \n" +
		"body: O cupom de desconto : " + coupon + " foi utilizado com sucesso!"

	err := smtp.SendMail("smtp.gmail.com:587",
		smtp.PlainAuth("", from, pass, "smtp.gmail.com"),
		from, []string{email}, []byte(msg))

	if err != nil {
		result := Result{Status: err.Error()}
		return result
	}

	result := Result{
		Status: "mail send",
	}

	return result
}
