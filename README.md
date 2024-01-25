# credit-card-validateMe

## Description
Credit card validator as a REST API service, which validates ad-hoc a credit card number given in the GET request body based on the Luhn algorhitm. It will also output the major credit card network the number is associated with (Visa, MasterCard or AMEX). 

## WHY 
Very frequently I purchase groceries online and each time I have to enter a card number to submit the online purchase. I always found it interesting that upon typing your card number, it already identifies the card issuer, even without fully entering the number. Sometimes I make typos and then I wonder why the payment failed, since I was sure I made correct input, but that wasn't the case. So these scenarios gave me an idea to try to implement smething similar in my project, at least the verification of the card number part. 

## Quick Start
I'm using mostly standard Go libraries, but I did import also these couple of external libraries: 
```
"github.com/alexflint/go-arg"
"github.com/gorilla/mux"
```
so make sure that the libraries are installed on your machine. As always, you can rely on the  **go mod tidy**  command before running or compilng the app. 

Then a REST API client is needed to make the GET request and input the credit card number in a JSON format in the request body. 
The JSON format to be used in the body is of type: **{"ccNum":"card_number"}**
As for the REST API client, I was using the Visual Studio Code **Thunder** extension to test, but you can use any other HTTP based client, i.e. Postman and the like. 

## Usage
The application starts an HTTP server, on which you can make **ONLY** a **GET** request. Typically for a GET request you wouldn't need to enter a request body, however in this case, since we want to only check the credit card number validity ad-hoc and not store the data anywhere for a regular retrieval, you would have to fill in the request body. 

You can run the application with the command
```
go run ./ccvalidator
```
You should then notice a message that the REST API server started and is listening for any requests. 
Since we were using the **go-arg** external package, you can add to the go run command an IP Address argument if you want to run this on a specific IP Address, otherwise the server will run on localhost on port 8080 as default. 
In the HTTP client of your choice, make sure you choose as request type  **GET**,  **http://127.0.0.1:8080/creditcard/get** in the URL if it runs on localhost and in the request body enter as JSON content the following format:  **{"ccNum":"card_number"}** . Lastly, hit the Send button on the http client to send the request. As output, you should see a following type of response back: 
```
{
  "cardNum": "the_card_number_you_made_the_request_with",
  "isValid": "true",   // or can be false if the number is not valid
  "cardIssuer": "Visa" // other choices are MasterCard, AMEX or invalid  if the card is not identified as being from one of those 3 issuers/types
}
```

The card number to check must be either 15 or 16 digits, otherwise the request will error out with bad request as HTTP status in response, along with the message to notify you to enter a 15 or 16 digit card number. 
