package parser

import "fmt"

type tokenCheckFn = func(Token) bool

type ruleMatch struct {
	name   string
	tokens []Token
	rules  []ruleMatch
}

type ruledef struct {
	name  string
	check func([]Token) (*ruleMatch, []Token) // returns (matched tokens, remaining tokens)
}

func tokenType(name string, t TokenType) ruledef {
	return ruledef{
		name: name,
		check: func(tokens []Token) (*ruleMatch, []Token) {
			if len(tokens) < 1 {
				return nil, tokens
			} else if tokens[0].Type == t {
				return &ruleMatch{name: name, tokens: tokens[0:1], rules: nil}, tokens[1:]
			} else {
				return nil, tokens
			}
		},
	}
}

func tokenTypeAndValue(name, value string, t TokenType) ruledef {
	return ruledef{
		name: name,
		check: func(tokens []Token) (*ruleMatch, []Token) {
			if len(tokens) < 1 {
				return nil, tokens
			} else if tokens[0].Type == t && tokens[0].Value == value {
				return &ruleMatch{name: name, tokens: tokens[0:1], rules: nil}, tokens[1:]
			} else {
				return nil, tokens
			}
		},
	}
}

func or(name string, rules ...ruledef) ruledef {
	return ruledef{
		name: name,
		check: func(tokens []Token) (*ruleMatch, []Token) {
			for _, rule := range rules {
				match, remaining := rule.check(tokens)
				if match != nil {
					return &ruleMatch{name: name, tokens: nil, rules: []ruleMatch{*match}}, remaining
				}
			}
			return nil, tokens
		},
	}
}

func seq(name string, rules ...ruledef) ruledef {
	return ruledef{
		name: name,
		check: func(tokens []Token) (*ruleMatch, []Token) {
			remainingTokens := tokens
			matches := []ruleMatch{}
			for _, rule := range rules {
				match, rest := rule.check(remainingTokens)
				if match != nil {
					remainingTokens = rest
					matches = append(matches, *match)
				} else {
					return nil, tokens
				}
			}
			return &ruleMatch{name: name, tokens: nil, rules: matches}, remainingTokens
		},
	}
}

func mult(name string, rule ruledef) ruledef {
	return ruledef{
		name: name,
		check: func(tokens []Token) (*ruleMatch, []Token) {
			remainingTokens := tokens
			matches := []ruleMatch{}
			done := false
			for !done {
				match, rest := rule.check(remainingTokens)
				if match != nil {
					remainingTokens = rest
					matches = append(matches, *match)
				} else if len(matches) > 0 {
					done = true
				} else {
					return nil, tokens
				}
			}
			return &ruleMatch{name: name, tokens: nil, rules: matches}, remainingTokens
		},
	}
}

var QueryRule = seq("Query",
	mult("Directives",
		tokenType("Directive", TOKEN_DIRECTIVE)),
	mult("Blocks",
		seq("Block",
			seq("FromClause",
				tokenTypeAndValue("FromClauseKeyword", "from", TOKEN_KEYWORD),
				tokenType("FromClauseValue", TOKEN_KEYWORD)),
			seq("WithClause",
				tokenTypeAndValue("WithClauseKeyword", "with", TOKEN_KEYWORD),
				mult("WithClauseItems",
					seq("WithClauseItem",
						tokenType("WithClauseItemKey", TOKEN_KEYWORD),
						tokenType("WithClauseItemEqualSign", TOKEN_EQUALS),
						or("WithClauseItemValue",
							tokenType("WithClauseItemValueString", TOKEN_STRING),
							tokenType("WithClauseItemValueNumber", TOKEN_NUMBER))))))))

func ParseRule(tokens []Token) {
	match, remaining := QueryRule.check(tokens)

	fmt.Printf("%+v\n\n", match)
	fmt.Printf("%+v", remaining)

}
