package parser

import (
	"restmux/model"
)

func ParseString(query string) model.RestmuxQuery {
	tokens, err := Tokenize(query)

	if err != nil {
		panic(err)
	}

	ParseRule(tokens)

	return model.RestmuxQuery{}
}
