package parser

import "strings"

type TokenType int

const (
	TOKEN_ERROR TokenType = iota
	TOKEN_EOF

	TOKEN_DIRECTIVE

	TOKEN_KEYWORD
	TOKEN_NUMBER
	TOKEN_STRING
	TOKEN_EQUALS
)

const EOF rune = 0
const KEYWORD_CHARS = "abcdefghijklmnopqrstuvwxyz-_"
const DIGIT_CHARS = "0123456789"
const LETTER_CHARS = "abcdefghijklmnopqrstuvwxyz"

type Token struct {
	Type   TokenType
	Value  string
	Line   int
	Column int
}

func isDigit(ch string) bool {
	return len(ch) == 1 && strings.Contains(DIGIT_CHARS, ch)
}

func isLetter(ch string) bool {
	return len(ch) == 1 && strings.Contains(LETTER_CHARS, strings.ToLower(ch))
}

func isKeywordChar(ch string) bool {
	return len(ch) == 1 && strings.Contains(KEYWORD_CHARS, strings.ToLower(ch))
}

func Tokenize(query string) ([]Token, error) {

	tokens := []Token{}
	index := 0
	done := false

	hasNextToken := false
	var nextToken Token
	buffer := ""

	line := 1
	col := 1

	for !done {
		ch := string(query[index])

		if hasNextToken {
			switch nextToken.Type {
			case TOKEN_DIRECTIVE:
				if ch == "\n" {
					nextToken.Value = buffer
					tokens = append(tokens, nextToken)
					hasNextToken = false
				} else {
					buffer += ch
				}
			case TOKEN_NUMBER:
				if isDigit(ch) {
					buffer += ch
				} else {
					nextToken.Value = buffer
					tokens = append(tokens, nextToken)
					hasNextToken = false
				}
			case TOKEN_STRING:
				if ch != "\"" {
					buffer += ch
				} else {
					nextToken.Value = buffer
					tokens = append(tokens, nextToken)
					hasNextToken = false
				}
			case TOKEN_KEYWORD:
				if isKeywordChar(ch) {
					buffer += ch
				} else {
					nextToken.Value = buffer
					tokens = append(tokens, nextToken)
					hasNextToken = false
				}
			}
		} else {
			if ch == "#" {
				buffer = ch
				nextToken = Token{Type: TOKEN_DIRECTIVE, Line: line, Column: col}
				hasNextToken = true
			} else if ch == "=" {
				tokens = append(tokens, Token{Type: TOKEN_EQUALS, Line: line, Column: col})
			} else if ch == "\"" {
				buffer = ""
				nextToken = Token{Type: TOKEN_STRING, Line: line, Column: col}
				hasNextToken = true
			} else if isDigit(ch) {
				buffer = ch
				nextToken = Token{Type: TOKEN_NUMBER, Line: line, Column: col}
				hasNextToken = true
			} else if isLetter(ch) {
				buffer = ch
				nextToken = Token{Type: TOKEN_KEYWORD, Line: line, Column: col}
				hasNextToken = true
			}
		}

		if query[index] == '\n' {
			col = 1
			line++
		} else {
			col++
		}

		index++
		if index >= len(query) {
			done = true
			if hasNextToken {
				nextToken.Value = buffer
				tokens = append(tokens, nextToken)
			}
			tokens = append(tokens, Token{Type: TOKEN_EOF, Line: line, Column: col})
		}
	}

	return tokens, nil

}
