# 2pln

This is simple currency converter from EUR, USD, GBP, CHF to PLN.
It's writen in golang, and uses official NBP currency data as source of data.

## Usage

	$ 2pln 10.0 EUR
	$ 2022-01-28; 1 EUR = 4.5697 PLN; 10.0000 EUR = 45.6970 PLN

* The first value is the date of the table.
* The second value is actual exchange rate for chosen currency
* The third argument is calculated by 2pln and represents value in PLN for given ealier amount.

## Compiling:
To compile this app you should have working instalation of golang.
You can download it from https://go.dev/

After that you should download my repository: https://github.com/mczechyra/2pln.git
You can use git for this:

	git clone https://github.com/mczechyra/2pln.git
	go build

## Extendig:
Main concept for this simple app is ApiProvider interface.

	type ApiProvider interface {
		GetCurrentRate(context.Context, UserRequest) (ApiResponce, error)
	}

By implement it for your own types you could add more data providers.