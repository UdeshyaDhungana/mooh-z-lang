package object

import (
	"bytes"
	"fmt"
	"strings"

	"github.com/udeshyadhungana/interprerer/app/ast"
)

type ObjectType string

const (
	INTEGER_OBJ       ObjectType = "INTEGER"
	BOOLEAN_OBJ       ObjectType = "BOOLEAN"
	NULL_OBJ          ObjectType = "NULL"
	PATHA_MUJI_OBJ    ObjectType = "RETURN"
	GALAT_MUJI_OBJ    ObjectType = "ERROR"
	KAAM_GAR_MUJI_OBJ ObjectType = "KAAM_GAR"
	STRING            ObjectType = "STRING"
	BUILTIN_OBJECT    ObjectType = "BUILTIN"
	ARRAY_OBJECT      ObjectType = "ARRAY"
	HASHMAP_OBJECT    ObjectType = "HASHMAP"
)

var (
	TRUE  = &Boolean{Value: true}
	FALSE = &Boolean{Value: false}
	NULL  = &Null{}
)

// common for all the data types
type Object interface {
	Type() ObjectType
	Inspect() string
}

// Integer
type Integer struct {
	Value int64
}

func (i *Integer) Inspect() string {
	return fmt.Sprintf("%d", i.Value)
}

func (i *Integer) Type() ObjectType {
	return INTEGER_OBJ
}

// Boolean
type Boolean struct {
	Value bool
}

func (b *Boolean) Type() ObjectType { return BOOLEAN_OBJ }
func (b *Boolean) Inspect() string {
	if b.Value {
		return "sacho_muji"
	} else {
		return "jhut_muji"
	}
}

// Null
type Null struct{}

func (n *Null) Type() ObjectType { return NULL_OBJ }
func (n *Null) Inspect() string  { return "khali_muji" }

// Return Object to signify the end of statements execution
type Return struct {
	Value Object
}

func (r *Return) Type() ObjectType { return PATHA_MUJI_OBJ }
func (r *Return) Inspect() string  { return r.Value.Inspect() }

// Error handling
type Error struct {
	Message string
}

func (e *Error) Type() ObjectType { return GALAT_MUJI_OBJ }
func (e *Error) Inspect() string  { return "ERROR: " + e.Message }

// Kaam gar
type KaamGar struct {
	Parameters []*ast.Identifier
	Body       *ast.BlockStatement
	Env        *Environment
}

func (f *KaamGar) Type() ObjectType { return KAAM_GAR_MUJI_OBJ }
func (f *KaamGar) Inspect() string {
	var out bytes.Buffer
	params := []string{}
	for _, p := range f.Parameters {
		params = append(params, p.String())
	}
	out.WriteString("fn")
	out.WriteString("(")
	out.WriteString(strings.Join(params, ", "))
	out.WriteString(") {\n")
	out.WriteString("...")
	out.WriteString("\n}")
	return out.String()
}

// string
type String struct {
	Value string
}

func (s *String) Inspect() string {
	return s.Value
}

func (s *String) Type() ObjectType {
	return STRING
}

// builtin
type BuiltinFunction func(args ...Object) Object

type Builtin struct {
	Fn BuiltinFunction
}

func (b *Builtin) Type() ObjectType { return BUILTIN_OBJECT }
func (b *Builtin) Inspect() string  { return "builitn function" }

// array
type Array struct {
	Arr []Object
}

func (a *Array) Inspect() string {
	var result bytes.Buffer
	result.WriteString("[")
	elems := []string{}
	for _, e := range a.Arr {
		elems = append(elems, e.Inspect())
	}
	result.WriteString(strings.Join(elems, ", "))
	result.WriteString("]")
	return result.String()
}

func (a *Array) Type() ObjectType {
	return ARRAY_OBJECT
}

// hashmap
type HashMap struct {
	Pairs map[string]Object
}

func (h *HashMap) Inspect() string {
	var result bytes.Buffer
	result.WriteString("{")
	var elems []string
	var current string
	for k, v := range h.Pairs {
		current = fmt.Sprintf("%s : %s", k, v.Inspect())
		elems = append(elems, current)
	}
	result.WriteString(strings.Join(elems, ", "))
	result.WriteString("}")
	return result.String()
}

func (h *HashMap) Type() ObjectType {
	return HASHMAP_OBJECT
}
