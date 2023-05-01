package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"log"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/mczechyra/2pln/api"
)

func main() {
	flagTimeout := flag.Int("timeout", 5, "timeout for request")

	userInput := strings.Join(os.Args[1:], " ")
	if len(userInput) == 0 {
		fmt.Println(info)
		return
	}
	userReq, err := decodeUserInput(userInput)
	if err != nil {
		fmt.Printf("%s\nbłąd dekodowania danych wejściowych: %v\n", info, err)
		return
	}

	// Użyj danych z NBP:
	var apiProvider = api.NbpApiProvider{}

	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(*flagTimeout)*time.Second)
	defer cancel()

	resp, err := apiProvider.GetCurrentRate(ctx, userReq)
	if err != nil {
		fmt.Println(info)
		log.Println(err)
		return
	}

	fmt.Printf("%s; 1 %s = %.4f PLN; %.4f %s = %.4f PLN\n",
		resp.Date,
		resp.CurrCode,
		resp.Value,
		userReq.Value,
		userReq.Curr,
		userReq.Value*resp.Value,
	)
}

// decodeUserInput convert user input into userRequest.
// Default currency is EUR.
func decodeUserInput(str string) (api.UserRequest, error) {
	str = strings.ReplaceAll(str, ",", ".")
	str = strings.ToUpper(strings.TrimSpace(str))
	str = strings.ReplaceAll(str, " ", "")

	rxValue := regexp.MustCompile(`^\d*\.?\,?\d*`)
	list := rxValue.FindAllString(str, 1)
	if list == nil {
		return api.UserRequest{}, errors.New("nie znaleziono wartości, lista == null")
	}
	if len(list) == 0 {
		return api.UserRequest{}, errors.New("nie znaleziono wartości, pusta lista")
	}
	// change ',' to '.';
	val := list[0]
	f, err := strconv.ParseFloat(val, 64)
	if err != nil {
		return api.UserRequest{}, err
	}

	rxCurr := regexp.MustCompile(`EUR|USD|GBP|CHF`)
	strCurr := strings.ToUpper(strings.TrimSpace(rxCurr.FindString(str)))

	var curr api.Currency
	switch strCurr {
	case "EUR":
		return api.UserRequest{Value: f, Curr: api.EUR}, nil
	case "USD":
		return api.UserRequest{Value: f, Curr: api.USD}, nil
	case "GBP":
		return api.UserRequest{Value: f, Curr: api.GBP}, nil
	case "CHF":
		return api.UserRequest{Value: f, Curr: api.CHF}, nil
	}
	return api.UserRequest{Value: f, Curr: curr}, nil
}

const info = `
      _______  ________  ___       ________      
     /  ___  \|\   __  \|\  \     |\   ___  \    
    /__/|_/  /\ \  \|\  \ \  \    \ \  \\ \  \   
    |__|//  / /\ \   ____\ \  \    \ \  \\ \  \  
        /  /_/__\ \  \___|\ \  \____\ \  \\ \  \ 
       |\________\ \__\    \ \_______\ \__\\ \__\
        \|_______|\|__|     \|_______|\|__| \|__|
=========================================================
2pln converts given currency to Polish zloty [PLN].
Default source of data is NBP (National Polish Bank)
Example of use:
.\2pln.exe {AMOUNT} {CURRENCY} (TIMEOUT)
where:
  AMOUNT are digits with or without coma or dot.
  valid AMOUNT options are: 
    - 10
    - 10,0
    - 10.0
  
  CURRENCY is one option from list below:
    - EUR (default option)
    - USD
    - GBP
    - CHF

  TIMEOUT is optional paramter for set request timeout.
  Default value for timeout is set as 5.
  Timeout have to be set as last parameter!
  To set your own timeout write: -timeout=10
  
Example of use:
.\2pln.exe 10 EUR
.\2pln.exe 12 USD -timeout 1

Sample output:
2022-01-28; 1 EUR = 4.5697 PLN; 10.0000 EUR = 45.6970 PLN`
