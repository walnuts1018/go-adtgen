package model

import "go/token"

type DeclarationKind string

const (
	DeclarationKindProduct DeclarationKind = "product"
	DeclarationKindSum     DeclarationKind = "sum"
)

type Declaration struct {
	Kind             DeclarationKind
	Name             string
	Expression       string
	Options          DeclarationOptions
	TypeParameters   []string
	InterfaceMethods []DeclaredInterfaceMethod
	Position         token.Position
	SourceFilename   string
}

type DeclarationOptions struct {
	NoSetter bool
}

type DeclaredInterfaceMethod struct {
	Name      string
	Signature string
}
