/*
Author: Adrian Potra
Version 1.0

Will implement a project that validates CC cards and identifies the card network. The validation will use the Luhn algorithm.
The project will implement a REST API and a GET method for the validation
*/

package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sync"

	"github.com/alexflint/go-arg"
	"github.com/gorilla/mux"
)

var args struct {
	BindRest string `arg:"env:CARDVALIDATOR_BIND_REST"`
}

type BankCard struct {
	CardNumber string `json:"cardNum"`
	IsValid    string `json:"isValid"`
	CardIssuer string `json:"cardIssuer"`
}

type CCNumber struct {
	CCnum string `json:"ccNum"`
}

// function takes a card number, including the check digit, as an array of integers and outputs true if the check digit is correct, false otherwise
func CardNumValidation(ccNum *CCNumber) *BankCard {

	// implementing Luhn's algorithm
	sum := 0
	parity := len(ccNum.CCnum) % 2
	//fmt.Println("length ccNum is ", len(ccNum.CCnum))
	for i, val := range ccNum.CCnum {
		d := val - 48

		//fmt.Println("input element", d)
		if i%2 != parity {
			sum += int(d)
			//fmt.Println(" i%2 != parity ", sum)
		} else if d > 4 {
			sum += 2*int(d) - 9
			//fmt.Println(" d> 4 ", sum)
		} else {
			sum += 2 * int(d)
			//fmt.Println(" else sum ", sum)
		}

	}
	//fmt.Println("bool check : ", len(ccNum) == (10-(sum%10)))
	//fmt.Println("bool check2 : ", sum%10 == 0)

	//return len(ccNum) == (10 - (sum % 10))
	ccIssuer := CardIssuer(ccNum.CCnum)
	if sum%10 == 0 {
		//fmt.Println("ccNum is ", ccNum)
		//fmt.Println("ccissuer is ", ccIssuer)
		return &BankCard{CardNumber: ccNum.CCnum, IsValid: "true", CardIssuer: ccIssuer}
	} else {
		return &BankCard{CardNumber: ccNum.CCnum, IsValid: "false", CardIssuer: ccIssuer}
	}
}

func CardIssuer(ccNum string) string {

	cardIssuer := ""

	for i, val := range ccNum {

		d := val - 48

		if i == 0 {

			switch d {
			case 3:
				//fmt.Println("Card issuer: American Express")
				cardIssuer = "American Express"
			case 4:
				//fmt.Println("Card issuer: Visa")
				cardIssuer = "Visa"
			case 5:
				//fmt.Println("Card issuer: Mastercard")
				cardIssuer = "Mastercard"
			default:
				//fmt.Println("Card issuer: Not known")
				cardIssuer = "invalid"
			}
		}

	}
	return cardIssuer
}

func GetCCValid(w http.ResponseWriter, req *http.Request) { // the w is what we will write responses into and the req is what the client sends and we can read from that request
	w.Header().Set("Content-Type", "application/json; charset= utf-8")
	//fmt.Println("request body is: ", req.Body)

	//we need to make sure that its a GET method
	if req.Method != "GET" {
		return
	}
	// creating empty struct - will be populated from the JSON data
	ccNum := CCNumber{}
	cc := BankCard{}
	err := json.NewDecoder(req.Body).Decode(&ccNum)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	fmt.Println("ccNum length is: ", len(ccNum.CCnum))

	if len(ccNum.CCnum) <= 14 || len(ccNum.CCnum) > 16 {
		http.Error(w, "Invalid input. Please enter a 15 or 16 digit card number to validate", http.StatusBadRequest)
		return
	}
	validation := CardNumValidation(&ccNum)
	cc = *validation

	// we'll return a response to get the information out of the JSON object
	json.NewEncoder(w).Encode(cc)

}

func main() {

	// parsing arguments
	arg.MustParse(&args)
	//setting defaults

	if args.BindRest == "" {
		args.BindRest = ":8080" // if we don't send ip address, it will default to localhost
	}

	router := mux.NewRouter()

	// Define the endpoints for CRUD operations and start the HTTP server
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		log.Printf("starting REST API server...\n")
		router.HandleFunc("/creditcard/get", GetCCValid).Methods("GET")
		log.Fatal(http.ListenAndServe(args.BindRest, router))
		wg.Done()
	}()

	wg.Wait()
}

