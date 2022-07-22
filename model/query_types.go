package model

type RestmuxQuery struct {
	Name      string
	Endpoints []RestmuxEndpoint
	Blocks    []RestmuxQueryBlock
}

type RestmuxEndpoint struct {
	Alias string
	Url   string
}

type RestmuxQueryBlock struct {
	From string
	With []RestmuxQueryWithBlock
}

type RestmuxQueryWithBlock struct {
	Field string
	Value string
}
