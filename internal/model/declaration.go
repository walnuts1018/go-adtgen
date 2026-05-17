package model

import "go/token"

type Declaration struct {
	Name           string
	Expression     string
	TypeParameters []string
	Position       token.Position
}
