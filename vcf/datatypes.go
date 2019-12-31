package vcf

import (
	"fmt"
	"strconv"
	"strings"
)

// DatatypeType represents the classes of Datatypes
type DatatypeType int8

// enum for Datatypes
const (
	StringType DatatypeType = iota
	StringVectorType
	IntegerType
	IntegerVectorType
	FloatType
	FloatVectorType
	CharacterType
	CharacterVectorType
	FlagType
	numDatatypes
)

func (d DatatypeType) String() string {
	var s string
	switch d {
	case StringType, StringVectorType:
		s = "String"
	case IntegerType, IntegerVectorType:
		s = "Integer"
	case FloatType, FloatVectorType:
		s = "Float"
	case CharacterType, CharacterVectorType:
		s = "Character"
	case FlagType:
		s = "Flag"
	}
	return s
}

// String represents a single string value
type String string

// StringVector represent multiple string values
type StringVector []String

// Integer represents a single integer value
type Integer int64

// IntegerVector represent multiple integer values
type IntegerVector []Integer

// Float represents a single float value
type Float float64

// FloatVector represent multiple float values
type FloatVector []Float

// Character represents a single byte value
type Character byte

// CharacterVector represent multiple byte values
type CharacterVector []Character

// Flag represent a single flag
type Flag bool

// Datatype is the interface implemented by VCF data types
type Datatype interface {
	typeOf() DatatypeType
}

// typeOf interface
func (String) typeOf() DatatypeType          { return StringType }
func (StringVector) typeOf() DatatypeType    { return StringVectorType }
func (Integer) typeOf() DatatypeType         { return IntegerType }
func (IntegerVector) typeOf() DatatypeType   { return IntegerVectorType }
func (Float) typeOf() DatatypeType           { return FloatType }
func (FloatVector) typeOf() DatatypeType     { return FloatVectorType }
func (Character) typeOf() DatatypeType       { return CharacterType }
func (CharacterVector) typeOf() DatatypeType { return CharacterVectorType }
func (Flag) typeOf() DatatypeType            { return FlagType }

// parse interface
func ParseString(str string) (String, error) {
	return String(str), nil
}

func ParseStringVector(str string) (StringVector, error) {
	parts := strings.Split(str, ",")
	var s StringVector
	for _, part := range parts {
		s = append(s, String(part))
	}
	return s, nil
}

func ParseInteger(s string) (Integer, error) {
	val, err := strconv.ParseInt(s, 10, 64)
	return Integer(val), err
}

func ParseIntegerVector(s string) (IntegerVector, error) {
	parts := strings.Split(s, ",")
	var vs IntegerVector
	for _, part := range parts {
		v, err := ParseInteger(part)
		if err != nil {
			return vs, err
		}
		vs = append(vs, v)
	}
	return vs, nil
}

func ParseFloat(s string) (Float, error) {
	val, err := strconv.ParseFloat(s, 64)
	return Float(val), err
}

func ParseFloatVector(s string) (FloatVector, error) {
	parts := strings.Split(s, ",")
	var vs FloatVector
	for _, part := range parts {
		v, err := ParseFloat(part)
		if err != nil {
			return vs, err
		}
		vs = append(vs, v)
	}
	return vs, nil
}

func ParseCharacter(s string) (Character, error) {
	var c Character
	if len(s) < 1 {
		return c, fmt.Errorf("string must have length > 0")
	}
	return Character(s[0]), nil
}

func ParseCharacterVector(s string) (CharacterVector, error) {
	parts := strings.Split(s, ",")
	var vs CharacterVector
	for _, part := range parts {
		v, err := ParseCharacter(part)
		if err != nil {
			return vs, err
		}
		vs = append(vs, v)
	}
	return vs, nil
}

func ParseFlag(s string) (Flag, error) {
	return Flag(true), nil
}

func ParseDatatype(t DatatypeType, s string) (Datatype, error) {
	var datatype Datatype
	var err error
	switch t {
	case IntegerType:
		datatype, err = ParseInteger(s)
	case IntegerVectorType:
		datatype, err = ParseIntegerVector(s)
	case FloatType:
		datatype, err = ParseFloat(s)
	case FloatVectorType:
		datatype, err = ParseFloatVector(s)
	case StringType:
		datatype, err = ParseString(s)
	case StringVectorType:
		datatype, err = ParseStringVector(s)
	case CharacterType:
		datatype, err = ParseCharacter(s)
	case CharacterVectorType:
		datatype, err = ParseCharacterVector(s)
	case FlagType:
		datatype, err = ParseFlag(s)
	}
	return datatype, err
}
