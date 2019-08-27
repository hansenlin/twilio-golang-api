package main

import (
	"log"
	"fmt"
	"strings"
	"net/http"
	"net/url"
	"encoding/json"

	"github.com/gorilla/mux"
)

func sendRequest(to string, body string) {
	// set account keys & information
	accountSid := "ACXXXX"
	authToken := "XXXXXX"
	urlStr := "https://api.twilio.com/2010-04-01/Accounts/" + accountSid + "/Messages.json"

	msgData := url.Values{}
	msgData.Set("From", "NUMBER_FROM")
	msgData.Set("To", to)
	msgData.Set("Body", body)
	msgDataReader := *strings.NewReader(msgData.Encode())

	// create http request client
	client := &http.Client{}
	req, _ := http.NewRequest("POST", urlStr, &msgDataReader)
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	req.SetBasicAuth(accountSid, authToken)

	// make http post request and return message SID
	resp, _ := client.Do(req)
	if (resp.StatusCode >= 200 && resp.StatusCode < 300) {
		var data map[string]interface{}
		decoder := json.NewDecoder(resp.Body)
		err := decoder.Decode(&data)
		if (err == nil) {
			fmt.Println(data["sid"])
		}
	} else {
		fmt.Println(resp.StatusCode)
	}
}

func params(w http.ResponseWriter, r *http.Request) {
	msgBody := r.FormValue("msgBody")
	pathParams := mux.Vars(r)
	w.Header().Set("Content-Type", "application/json")

	telNum, ok := pathParams["telNum"]
	if !ok {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{"message": "need a number"}`))
		return
	}

	sendRequest(telNum, msgBody)

	w.Write([]byte(fmt.Sprintf(`{"phone_number": %s, "message_sent": %s }`, telNum, msgBody)))
}

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/{telNum}", params).Methods(http.MethodPost)
	log.Fatal(http.ListenAndServe(":8080", r))
}
