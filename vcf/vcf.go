package vcf

import (
	"bufio"
	"fmt"
	"io"
	"regexp"
	"strconv"
	"strings"
)

// Table represents a collection of VCF records and associated
// metadata
type Table struct {
	Fileformat string
	Metadata   []*Metadata
	Info       map[string]*Metadata
	Format     map[string]*Metadata
	Records    []*Record
}

// NewTable initializes a vcf.Table struct
func NewTable() Table {
	var tbl Table
	tbl.Info = make(map[string]*Metadata)
	tbl.Format = make(map[string]*Metadata)
	return tbl
}

// Metadata represents data for a VCF metafield
type Metadata struct {
	class       string
	ID          string
	Description string
	Number      string
	Type        MetadataType
	Source      string
	Version     string
	Other       map[string]string
	Value       string
	isKeyValue  bool
}

// NewMetadata is a constructor for the metafield struct
func NewMetadata() Metadata {
	var result Metadata
	result.Other = make(map[string]string)
	return result
}

// MetadataType represents the data types associated with different
// VCF metainformation
type MetadataType int

// enum for MetaDataType
const (
	None MetadataType = iota
	Integer
	Float
	String
	Character
	Flag
	NumMetadataTypes
)

func (i MetadataType) String() string {
	return [...]string{"None", "Integer", "Float", "String", "Character", "Flag"}[i]
}

func parseMetadataType(s string) (MetadataType, error) {
	switch s {
	case "Integer":
		return Integer, nil
	case "Float":
		return Float, nil
	case "String":
		return String, nil
	case "Character":
		return Character, nil
	case "Flag":
		return Flag, nil
	default:
		return None, fmt.Errorf("invalid INFO type")
	}
}

// Record is a representation of a VCF record
type Record struct {
	// The 8 required fields
	Chrom     string
	Pos       int
	ID        string
	Ref       string
	Alt       string
	Qual      float64
	Filter    string
	Info      *InfoData
	Genotypes *GenotypeData
	HasQual   bool
}

// InfoData is a struct that represents the INFO field in a VCF record
type InfoData struct {
	Integers   map[string][]int64
	Floats     map[string][]float64
	Strings    map[string][]string
	Characters map[string][]byte
	Flags      map[string]bool
}

// NewInfoData is a constructor for an InfoData struct
func NewInfoData() InfoData {
	var result InfoData
	result.Integers = make(map[string][]int64)
	result.Floats = make(map[string][]float64)
	result.Strings = make(map[string][]string)
	result.Characters = make(map[string][]byte)
	result.Flags = make(map[string]bool)
	return result
}

var convArray = [...]{

}

func parseInfoData(
	infomap map[string]*Metadata, s string) (InfoData, error) {

	result := NewInfoData()
	var key, value string

	for _, field := range strings.Split(strings.TrimSpace(s), ";") {
		if strings.Contains(field, "=") {
			key, value = splitKeyEqualValueString(field)
		} else {
			key = field
		}
		if _, ok := infomap[key]; ok {
			datatype, err := parseMetadataType(key)
			if err != nil {
				continue
			}
			switch datatype {
			case None:
				panic("Invalid metadata type")
			case Integer:
				value, err := strconv.ParseInt(value, 10, 64)
				if err != nil {
					continue
				}
				result.Integers[key] = append(result.Integers[key], value)
			case Float:
				value, err := strconv.ParseFloat(value, 64)
				if err != nil {
					continue
				}
				result.Floats[key] = append(result.Floats[key], value)
			case String:
				result.Strings[key] = append(result.Strings[key], value)
			case Character:
				result.Characters[key] = append(result.Characters[key], byte(value[0]))
			case Flag:
				result.Flags[key] = true
			}
		} else {
			result.Strings[key] = append(result.Strings[key], value)
		}
	}
	return result, nil
}

// GenotypeData represents the genotype information in the genotype
// field of a VCF record
type GenotypeData struct {
	Format  string
	Samples []string
}

func parseGenotypeFields(fields []string) *GenotypeData {
	if fields == nil {
		return nil
	}
	gd := new(GenotypeData)
	gd.Format = fields[0]
	gd.Samples = fields[1:]
	return gd
}

var keyvalRe = regexp.MustCompile(`(?P<key>\w+)=(?P<value>["].+?["]|[^,]+?)(?:$|,)`)

func splitKeyEqualValueString(s string) (string, string) {
	matches := keyvalRe.FindStringSubmatch(s)
	if len(matches) != 3 {
		return "", ""
	}
	return matches[1], matches[2]
}

func parseMetadata(s string) (Metadata, error) {
	result := NewMetadata()
	result.Value = s
	if !strings.HasPrefix(s, "<") || !strings.HasSuffix(s, ">") {
		result.isKeyValue = false
		return result, nil
	}
	s = strings.Trim(s, "<>")

	matches := keyvalRe.FindAllStringSubmatch(s, -1)
	if len(matches) < 1 {
		return result, fmt.Errorf("No key=value pairs found in: <%s>", s)
	}

	for _, match := range matches {
		// note that match[0] is the entire match
		// while match[1] and match[2] are the desired groups
		key := match[1]
		value := match[2]
		switch key {
		case "ID":
			result.ID = value
		case "Description":
			result.ID = value
		case "Number":
			result.Number = value
		case "Type":
			dtype, _ := parseMetadataType(value)
			result.Type = dtype
		case "Source":
			result.Source = value
		case "Version":
			result.Version = value
		default:
			result.Other[key] = value
		}
	}
	return result, nil
}

// ParseFile parses a VCF file, returning a vcf.Table struct
func ParseFile(r io.Reader) (Table, error) {
	var records []*Record
	table := NewTable()

	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		// Empty line
		if len(line) == 0 {
			continue
		}
		// comment, header, or metainformation
		if strings.HasPrefix(line, "#") {
			if strings.HasPrefix(line, "#CHROM") { // header line
				continue
			}
			// metainformation
			if strings.HasPrefix(line, "##") {
				key, value := splitKeyEqualValueString(line)
				key = key[2:]
				meta, err := parseMetadata(value)
				if err != nil {
					return table, err
				}
				switch key {
				case "fileformat":
					table.Fileformat = value
				case "INFO":
					if meta.ID != "" {
						table.Info[meta.ID] = &meta
					}
				case "FORMAT":
					if meta.ID != "" {
						table.Format[meta.ID] = &meta
					}
				default:
					table.Metadata =
						append(table.Metadata, &meta)
				}
				continue
			}
		}
		// line with fields
		parts := strings.Split(line, "\t")
		if len(parts) < 8 {
			return table, fmt.Errorf("invalid VCF line")
		}
		r := new(Record)
		r.Chrom = parts[0]
		pos, err := strconv.ParseInt(parts[1], 0, 64)
		if err == nil {
			r.Pos = int(pos)
		}
		r.ID = parts[2]
		r.Ref = parts[3]
		r.Alt = parts[4]
		qual, err := strconv.ParseFloat(parts[5], 64)
		if err == nil {
			r.Qual = qual
			r.HasQual = true
		}
		r.Filter = parts[6]
		info, err := parseInfoData(table.Info, parts[7])
		if err == nil {
			r.Info = &info
		}
		if len(parts) >= 9 {
			r.Genotypes = parseGenotypeFields(parts[8:])
		}
		records = append(records, r)
	}
	table.Records = records
	return table, nil
}
