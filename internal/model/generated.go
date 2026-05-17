package model

import "go/types"

type GeneratedType struct {
	Name           string
	TypeParameters []string
	Fields         []GeneratedField
}

type GeneratedField struct {
	Name      string
	Type      types.Type
	Tag       string
	Anonymous bool
}
