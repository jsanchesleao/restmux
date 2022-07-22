package parser

import (
	"io/ioutil"
	"testing"
)

func TestParse(t *testing.T) {
	queryString, err := ioutil.ReadFile("../examples/simple.rmx")
	if err != nil {
		panic(err)
	}
	query := ParseString(string(queryString))

	if query.Name != "Simple query" {
		t.Log("Query name is wrong. Should be Simple query but was", query.Name)
		t.Fail()
	}
}
