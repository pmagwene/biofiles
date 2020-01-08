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
func NewTable() *Table {
	var tbl Table
	tbl.Info = make(map[string]*Metadata)
	tbl.Format = make(map[string]*Metadata)
	return &tbl
}

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

// NewMetadata is a constructor for the metafield struct
func NewMetadata() *Metadata {
	var result Metadata
	result.OtherFields = make(map[string]string)
	return &result
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
	// if len(matches) < 3 {
	// 	return result, fmt.Errorf("No key=value pairs found in: %s, %v", fieldstr, matches)
	// }

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

var keyvalRe = regexp.MustCompile(`(?P<key>\w+)=(?P<value>["].+?["]|[^,]+?)(?:$|,)`)

func splitKeyEqualValueString(s string) (string, string) {
	matches := keyvalRe.FindStringSubmatch(s)
	if len(matches) != 3 {
		return "", ""
	}
	return matches[1], matches[2]
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
	Info      map[string]Datatype
	Format    []string
	Genotypes [][]Datatype
	HasQual   bool
	Parent    *Table
}

// NewRecord constructs a vcf.Record
func NewRecord(parent *Table) *Record {
	var result Record
	result.Parent = parent
	result.Info = make(map[string]Datatype)
	return &result
}

// ParseRecord parses a string to a vcf.Record struct
func (r *Record) ParseRecord(s string) error {
	parts := strings.Split(s, "\t")
	if len(parts) < 8 {
		return fmt.Errorf("invalid VCF line")
	}
	r.Chrom = parts[0]
	pos, err := strconv.ParseInt(parts[1], 10, 64)
	if err != nil {
		return fmt.Errorf("invalid parsing of Pos")
	}
	r.Pos = int(pos)
	r.ID = parts[2]
	r.Ref = parts[3]
	r.Alt = parts[4]
	qual, err := strconv.ParseFloat(parts[5], 64)
	if err == nil {
		r.Qual = qual
		r.HasQual = true
	}
	r.Filter = parts[6]
	// INFO field
	err = r.parseInfo(parts[7])
	if err != nil {
		return err
	}
	// FORMAT field
	if len(parts) > 8 {
		r.Format = strings.Split(parts[8], ":")
	}
	// Sample genotypes
	if len(parts) > 9 {
		err = r.parseGenotypes(parts[9:])
		if err != nil {
			return err
		}
	}
	return err
}

func (r *Record) parseInfo(s string) error {

	var datatype Datatype
	var err error

	if len(s) < 1 {
		return fmt.Errorf("empty info string")
	}

	for _, field := range strings.Split(strings.TrimSpace(s), ";") {
		if !strings.Contains(field, "=") { // Flag field
			f := Flag(true)
			r.Info[field] = f
			continue
		}
		idval := strings.SplitN(field, "=", 2)
		id := idval[0]
		val := idval[1]

		if metadata, ok := r.Parent.Info[id]; ok {
			datatype, err = ParseDatatype(metadata.Type, val)
		} else {
			datatype, err = ParseString(val)
		}
		if err != nil {
			return fmt.Errorf("error parsing Info field: %s", id)
		}
		r.Info[id] = datatype
	}
	return nil
}

func (r *Record) parseGenotypes(fields []string) error {
	var key string
	var typetype DatatypeType
	var datatype Datatype
	var err error

	for _, field := range fields {
		var sample []Datatype
		parts := strings.Split(field, ":")
		for i, part := range parts {
			if part == "." {
				datatype, err = ParseString(part)
			} else {
				key = r.Format[i]
				typetype = r.Parent.Format[key].Type
				datatype, err = ParseDatatype(typetype, part)
			}
			sample = append(sample, datatype)
		}
		r.Genotypes = append(r.Genotypes, sample)
	}
	return err
}

// ParseFile parses a VCF file, returning a vcf.Table struct
func ParseFile(r io.Reader) (*Table, error) {
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
				meta, err := ParseMetadata(line)
				if err != nil {
					return table, err
				}
				table.Metadata = append(table.Metadata, meta)
				switch meta.Class {
				case "fileformat":
					table.Fileformat = meta.Value
				case "INFO":
					if meta.ID != "" {
						table.Info[meta.ID] = meta
					}
				case "FORMAT":
					if meta.ID != "" {
						table.Format[meta.ID] = meta
					}
				}
				continue
			}
		}
		r := NewRecord(table)
		r.ParseRecord(line)
		records = append(records, r)
	}
	table.Records = records
	return table, nil
}
