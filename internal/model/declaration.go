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
	Options        DeclarationOptions
	TypeParameters []string
	Position       token.Position
	SourceFilename string
}

type DeclarationOptions struct {
	NoSetter bool
}
