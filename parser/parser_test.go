package parser

import "testing"

func TestParseOneClause(t *testing.T) {

	query, err := ParseString(`
from user
with id = 1
`)

	if err != nil {
		t.Fatal(err)
	}

	if len(query.Blocks) != 1 {
		t.Fatalf("Expected query to have %d blocks, but it had %d", 1, len(query.Blocks))
	}

	if query.Blocks[0].From != "user" {
		t.Fatalf("Expected from to be %q but was %q", "user", query.Blocks[0].From)
	}

	if query.Blocks[0].With[0].Field != "id" {
		t.Fatalf("Expected with field to be %q but was %q", "id", query.Blocks[0].With[0].Field)
	}

	if query.Blocks[0].With[0].Value != "1" {
		t.Fatalf("Expected with field to be %q but was %q", "1", query.Blocks[0].With[0].Value)
	}
}
