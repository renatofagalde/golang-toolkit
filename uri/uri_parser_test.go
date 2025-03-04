package uri

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

type ParsedURI struct {
	Method string `json:"method"`
	URI    string `json:"uri"`
}

func parseURI(method, uri string) (ParsedURI, error) {

	return ParsedURI{Method: method, URI: uri}, nil
}

func TestParseURI(t *testing.T) {
	testCases := []struct {
		method   string
		uri      string
		expected ParsedURI
	}{
		{
			"GET", "/pessoas/4152bd3d-9121-4889-b683-6ee72b8e3719",
			ParsedURI{"GET", "/pessoas/4152bd3d-9121-4889-b683-6ee72b8e3719"},
		},
		{
			"POST", "/pessoas/4152bd3d-9121-4889-b683-6ee72b8e3719/telefones",
			ParsedURI{"POST", "/pessoas/4152bd3d-9121-4889-b683-6ee72b8e3719/telefones"},
		},
		{
			"GET", "/pessoas/4152bd3d-9121-4889-b683-6ee72b8e3719/enderecos/65020de3-8a6c-416f-97d2-9755b0111bdc",
			ParsedURI{
				"GET", "/pessoas/4152bd3d-9121-4889-b683-6ee72b8e3719/enderecos/65020de3-8a6c-416f-97d2-9755b0111bdc",
			},
		},
		{
			"GET", "/pessoas/4152bd3d-9121-4889-b683-6ee72b8e3719/enderecos/65020de3-8a6c-416f-97d2-9755b0111bdc/tipos",
			ParsedURI{
				"GET",
				"/pessoas/4152bd3d-9121-4889-b683-6ee72b8e3719/enderecos/65020de3-8a6c-416f-97d2-9755b0111bdc/tipos",
			},
		},
		{
			"GET",
			"/pessoas/4152bd3d-9121-4889-b683-6ee72b8e3719/enderecos/65020de3-8a6c-416f-97d2-9755b0111bdc/tipos/326ab051-589e-404d-b674-808b3268e93d",
			ParsedURI{
				"GET",
				"/pessoas/4152bd3d-9121-4889-b683-6ee72b8e3719/enderecos/65020de3-8a6c-416f-97d2-9755b0111bdc/tipos/326ab051-589e-404d-b674-808b3268e93d",
			},
		},
	}

	for _, testCase := range testCases {
		t.Run(fmt.Sprintf("Parsing %s %s", testCase.method, testCase.uri), func(t *testing.T) {
			parsed, err := parseURI(testCase.method, testCase.uri)
			assert.Nil(t, err)

			parsedJSON, _ := json.Marshal(parsed)
			expectedJSON, _ := json.Marshal(testCase.expected)

			assert.JSONEq(t, string(expectedJSON), string(parsedJSON))
		})
	}
}
