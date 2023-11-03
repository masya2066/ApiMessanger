package models

import "go/types"

type Chat struct {
	Name    string      `json:"name"`
	Owners  types.Array `json:"owner"`
	Members types.Array `json:"members"`
}
