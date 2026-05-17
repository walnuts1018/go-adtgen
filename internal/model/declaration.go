package model

import "go/token"

type DeclarationKind string

const (
	DeclarationKindProduct DeclarationKind = "product"
	DeclarationKindSum     DeclarationKind = "sum"
)

type Declaration struct {
	Kind           DeclarationKind
	Name           string
	Expression     string
	TypeParameters []string
	Position       token.Position
}
