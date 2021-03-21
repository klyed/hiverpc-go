package types

type HscOperation struct {
	Sender  string   `json:"from"`
	Target string   `json:"target"`
	What      []string `json:"what"`
}
