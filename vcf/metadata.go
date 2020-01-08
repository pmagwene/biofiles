package vcf

import (
	"fmt"
	"strings"
)

// Metadata represents data for a VCF metafield
type Metadata struct {
	Class       string
	Value       string
	ID          string
	Description string
	Number      string
	Type        DatatypeType
	Source      string
	Version     string
	OtherFields map[string]string
}

// NewMetadata is a constructor for the metafield struct
func NewMetadata() *Metadata {
	var result Metadata
	result.OtherFields = make(map[string]string)
	return &result
}

func (m *Metadata) String() string {
	s := ""
	if m.ID == "" {
		s = fmt.Sprintf("(%s, %s)", m.Class, m.Value)
	} else {
		s = fmt.Sprintf("(%s, %s, %s, %s, %s)",
			m.Class, m.ID, m.Description, m.Number, m.Type)
	}
	return s
}

// ParseMetadata parses a string to a metadata struct
func ParseMetadata(s string) (*Metadata, error) {
	result := NewMetadata()
	if !strings.HasPrefix(s, "##") {
		return result, fmt.Errorf("invalid metadata line; must start with ##")
	}
	classval := strings.SplitN(s, "=", 2)
	if len(classval) != 2 {
		return result, fmt.Errorf("invalid metadata line; no ##field=value")
	}

	result.Class = classval[0][2:]
	result.Value = classval[1]
	val := classval[1]

	if !strings.HasPrefix(val, "<") || !strings.HasSuffix(val, ">") {
		return result, nil
	}

	fieldstr := strings.Trim(classval[1], "<>")
	matches := keyvalRe.FindAllStringSubmatch(fieldstr, -1)

	for _, match := range matches {
		// note that match[0] is the entire match
		// while match[1] and match[2] are the desired subgroups
		key := match[1]
		value := match[2]
		switch key {
		case "ID":
			result.ID = value
		case "Description":
			result.Description = value
		case "Number":
			result.Number = value
		case "Source":
			result.Source = value
		case "Version":
			result.Version = value
		case "Type":
			switch value {
			case "Integer":
				if result.Number == "1" {
					result.Type = IntegerType
				} else {
					result.Type = IntegerVectorType
				}
			case "Float":
				if result.Number == "1" {
					result.Type = FloatType
				} else {
					result.Type = FloatVectorType
				}
			case "String":
				if result.Number == "1" {
					result.Type = StringType
				} else {
					result.Type = StringVectorType
				}
			case "Character":
				if result.Number == "1" {
					result.Type = CharacterType
				} else {
					result.Type = CharacterVectorType
				}
			case "Flag":
				result.Type = FlagType
			}
		default:
			result.OtherFields[key] = value
		}
	}
	return result, nil
}
