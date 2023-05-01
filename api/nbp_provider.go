package api

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/pkg/errors"
)

// Main API adress
const adresNBPApi = "http://api.nbp.pl/api/exchangerates/rates"

// Chosen table code:
const tableCode = "a"

// Data format that api should return:
const responceFormat = "json"

// Example responce:
// {"table":"A","currency":"dolar amerykański","code":"USD","rates":[{"no":"219/A/NBP/2021","effectiveDate":"2021-11-12","mid":4.0559}]}
type nbpResp struct {
	Table    string `json:"table"`    // typ tabeli np: "A"
	CurrName string `json:"currency"` // nazwa waluty np: "dolar amerykański"
	CurrCode string `json:"code"`     // kod waluty np: "USD"
	Rates    []rate `json:"rates"`    // tebalu kursów walut: rates[]
}

type rate struct {
	NbpNum string  `json:"no"`            // "219/A/NBP/2021"
	Date   string  `json:"effectiveDate"` // data kursu np: "2021-11-12"
	Value  float64 `json:"mid"`           // wartość waluty np: 4.0559
}

// NbpApiProvider have to implemet ApiProvider interface:
var _ ApiProvider = NbpApiProvider{}

type NbpApiProvider struct{}

func (nbp NbpApiProvider) GetCurrentRate(ctx context.Context, reqData UserRequest) (ApiResponce, error) {
	type getResult struct {
		resp ApiResponce
		err  error
	}
	resp := make(chan getResult, 1)

	go func() {
		req, err := http.NewRequest("GET", getApiCall(reqData), nil)
		if err != nil {
			resp <- getResult{ApiResponce{}, errors.Wrap(err, "nie mogę utworzyć żądania http")}
			return
		}
		req.Header.Add("User-Agent", "2pln")
		req.Header.Add("Accept", "*/*")

		client := &http.Client{}
		httpResp, err := client.Do(req)
		defer httpResp.Body.Close()

		if err != nil {
			resp <- getResult{ApiResponce{}, errors.Wrap(err, "błąd pobrania danych z NBP (http.Get)")}
			return
		}

		if httpResp.StatusCode != http.StatusOK {
			s := fmt.Sprintf("status odpowiedzi: %s", httpResp.Status)
			resp <- getResult{ApiResponce{}, fmt.Errorf(s)}
			return
		}

		body, err := ioutil.ReadAll(httpResp.Body)
		if err != nil {
			resp <- getResult{ApiResponce{}, errors.Wrap(err, "błąd odczytu pobranych danych (ReadAll)")}
			return
		}

		nbpResp := &nbpResp{}
		if err := json.Unmarshal(body, nbpResp); err != nil {
			resp <- getResult{ApiResponce{}, errors.Wrap(err, "błąd Unmarshal JSON")}
			return
		}

		if len(nbpResp.Rates) == 0 {
			resp <- getResult{ApiResponce{}, errors.New("pusta lista kursów walut")}
			return
		}

		apiResp := ApiResponce{
			Date:     nbpResp.Rates[0].Date,
			Value:    nbpResp.Rates[0].Value,
			CurrCode: nbpResp.CurrCode,
			CurrName: nbpResp.CurrName,
		}
		resp <- getResult{apiResp, nil}
	}()

	result := getResult{}
	select {
	case result = <-resp:
	case <-ctx.Done():
		return ApiResponce{}, ctx.Err()
	}
	return result.resp, result.err
}

// getApiCall return get query for given data
// http://api.nbp.pl/api/exchangerates/rates/a/usd?format=json
func getApiCall(r UserRequest) string {
	return fmt.Sprintf("%s/%s/%s?format=%s",
		adresNBPApi,
		tableCode,
		r.Curr.String(),
		responceFormat,
	)
}
