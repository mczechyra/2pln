package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

// GET https://api.exchangeratesapi.io/latest?base=PLN
// RESP {"rates":{"EUR":0.2273088901,"USD":0.2683608756,"GBP":0.203452822,...},"base":"PLN","date":"2020-08-27"}

const (
	// file name of the cache in temp folder
	name = "2pln_cache.txt"
	// url to api.exchangeratesapi.io
	apiURL = "https://api.exchangeratesapi.io/latest?base=PLN"
)

// file path to cache file inside temp folder
var cacheFileName = filepath.Join(os.TempDir(), name)

type respCurr struct {
	Rates map[string]float64 `json:"rates"`
	Base  string             `json:"base"`
	Date  string             `json:"date"`
}

func main() {
	if len(os.Args) < 3 {
		log.Println("Not enough prameters")
		printUsage()
		return
	}

	// store data with currency in JSON format
	var rawJSON []byte

	// check for temp data insice system temp dir:
	rawJSON = checkForCache()

	// current currency data
	var rc respCurr

	currVal, err := strconv.ParseFloat(os.Args[1], 64)
	if err != nil {
		log.Println(err)
		return
	}

	currType := strings.ToUpper(os.Args[2])

	if len(rawJSON) > 0 {
		if err := json.Unmarshal(rawJSON, &rc); err != nil {
			log.Fatalf("Can't unmarchal json data: %v", err)
			return
		}
		printResult(rc, currVal, currType)
		return
	}

	// read current data from exchangeratesapi:
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// temporary type to pass data back from goroutine:
	type respResult struct {
		resp *http.Response
		err  error
	}

	chanResp := make(chan respResult)
	go func() {
		req, err := http.NewRequest("GET", apiURL, nil)
		if err != nil {
			chanResp <- respResult{
				resp: nil,
				err:  fmt.Errorf("can not create request: %v", err),
			}
		}
		req = req.WithContext(ctx)
		httpClient := &http.Client{}
		resp, err := httpClient.Do(req)
		if err != nil {
			chanResp <- respResult{
				resp: nil,
				err:  fmt.Errorf("could not dowload currency data: %v", err),
			}
		}
		chanResp <- respResult{resp: resp, err: nil}
	}()

	select {
	case rr := <-chanResp:
		defer rr.resp.Body.Close()
		if rr.err != nil {
			log.Fatal(err)
			return
		}
		rawJSON, err = ioutil.ReadAll(rr.resp.Body)
		if err != nil {
			log.Fatalf("can not read body: %v", err)
			return
		}

		if err := json.Unmarshal(rawJSON, &rc); err != nil {
			log.Fatalf("Can't unmarchal json data: %v", err)
			return
		}
		printResult(rc, currVal, currType)

		if err := saveCache(rawJSON); err != nil {
			log.Printf("can not update cache: %v", err)
		}
	case <-ctx.Done():
		log.Fatal("context deadline exceeded")
	}
}

func (r respCurr) Pln(pln float64, curr string) float64 {
	val, ok := r.Rates[curr]
	if !ok {
		return -1.0
	}
	return pln / val
}

func printUsage() {
	fmt.Println("USAGE: 2pln Value CurrCode")
	fmt.Println("Example: 2pln 15.00 EUR")
}

// Function opens temp file and checks save date.
// First line inside cache file contains save date.
// If save date is equal today then remove this date and returns rest content of the file
// If file is older, then returns empty string.
func checkForCache() []byte {
	b, err := ioutil.ReadFile(cacheFileName)
	if err != nil {
		log.Printf("can not read cache file %s: %v", cacheFileName, err)
		return []byte{}
	}

	// read cache data:
	buf := bytes.NewBuffer(b)
	cacheTimeStr, err := buf.ReadString('\n')
	if err != nil {
		log.Printf("can not read cache data: %v", err)
		return []byte{}
	}
	cacheTimeStr = strings.Trim(cacheTimeStr, "\r\n ")
	cacheDate, err := time.Parse("2006-01-02", cacheTimeStr)
	if err != nil {
		log.Printf("can not convert cache date: %v", err)
		return []byte{}
	}
	localTimeZone := time.Now().Location()
	cacheTime := time.Date(cacheDate.Year(), cacheDate.Month(), cacheDate.Day(), 1, 2, 3, 0, localTimeZone)
	localTime := time.Date(time.Now().Year(), time.Now().Month(), time.Now().Day(), 1, 2, 3, 0, localTimeZone)
	if !cacheTime.Equal(localTime) {
		return []byte{}
	}
	return buf.Next(buf.Len())
}

func saveCache(b []byte) error {
	buf := new(bytes.Buffer)
	// before write add current date before json content
	buf.WriteString(fmt.Sprintf("%s\n", time.Now().Format("2006-01-02")))
	_, err := buf.Write(b)
	if err != nil {
		return err
	}
	err = ioutil.WriteFile(cacheFileName, buf.Bytes(), 0644)
	if err != nil {
		return err
	}
	return nil
}

func printResult(r respCurr, currVal float64, currType string) {
	fmt.Printf("[%s] %.4f %s = %.4f PLN\n", r.Date, currVal, currType, r.Pln(currVal, currType))
}
