package main

import (
	"testing"

	"github.com/mczechyra/2pln/api"
	"github.com/stretchr/testify/assert"
)

func TestDecodeValidUserInput(t *testing.T) {
	var Tests = []struct {
		in  string
		out api.UserRequest
	}{
		{in: "0", out: api.UserRequest{Value: 0, Curr: api.EUR}},
		{in: "1", out: api.UserRequest{Value: 1, Curr: api.EUR}},
		{in: "1 eur", out: api.UserRequest{Value: 1, Curr: api.EUR}},
		{in: "10 eur", out: api.UserRequest{Value: 10, Curr: api.EUR}},
		{in: "101 Eur", out: api.UserRequest{Value: 101, Curr: api.EUR}},
		{in: "101 EUr", out: api.UserRequest{Value: 101, Curr: api.EUR}},
		{in: "101 EUR", out: api.UserRequest{Value: 101, Curr: api.EUR}},
		{in: "101.00 EUR", out: api.UserRequest{Value: 101.00, Curr: api.EUR}},
		{in: "101,00 EUR", out: api.UserRequest{Value: 101.00, Curr: api.EUR}},
		{in: "101,00 eur", out: api.UserRequest{Value: 101.00, Curr: api.EUR}},
		{in: "0101,00 EUR", out: api.UserRequest{Value: 101.00, Curr: api.EUR}},
		{in: "10,00 eur", out: api.UserRequest{Value: 10, Curr: api.EUR}},
		{in: "11,00 uSD", out: api.UserRequest{Value: 11, Curr: api.USD}},
		{in: "10.00 usD", out: api.UserRequest{Value: 10, Curr: api.USD}},
		{in: "10,0 usD", out: api.UserRequest{Value: 10, Curr: api.USD}},
		{in: "1. gbp", out: api.UserRequest{Value: 1, Curr: api.GBP}},
		{in: ".1 usD", out: api.UserRequest{Value: 0.1, Curr: api.USD}},
		{in: "1 chf", out: api.UserRequest{Value: 1, Curr: api.CHF}},
		{in: "99999. eur", out: api.UserRequest{Value: 99999, Curr: api.EUR}},
		{in: "99 999. eur", out: api.UserRequest{Value: 99999, Curr: api.EUR}},
		{in: "99 999.00 eur", out: api.UserRequest{Value: 99999, Curr: api.EUR}},
		{in: "9 9 99 9 ,00 eur", out: api.UserRequest{Value: 99999, Curr: api.EUR}},
		{in: "10.1eur", out: api.UserRequest{Value: 10.1, Curr: api.EUR}},
		{in: " 1 0 . 1 e u r ", out: api.UserRequest{Value: 10.1, Curr: api.EUR}},
	}
	for _, tt := range Tests {
		t.Run(tt.in, func(t *testing.T) {
			r, e := decodeUserInput(tt.in)
			assert.NoError(t, e)
			assert.Equal(t, tt.out.Value, r.Value, "błędna wartość")
			assert.Equal(t, tt.out.Curr, r.Curr, "błędna waluta")
		})
	}
}

func TestDecodeInvalidUserInput(t *testing.T) {
	var Tests = []struct {
		in  string
		out api.UserRequest
	}{
		{in: "", out: api.UserRequest{}},
		{in: "euro", out: api.UserRequest{}},
		{in: "+", out: api.UserRequest{}},
		{in: "-+", out: api.UserRequest{}},
		{in: "+-", out: api.UserRequest{}},
	}
	for _, tt := range Tests {
		t.Run(tt.in, func(t *testing.T) {
			_, e := decodeUserInput(tt.in)
			assert.Error(t, e)
		})
	}
}
