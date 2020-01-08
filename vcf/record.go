package vcf

import (
	"fmt"
	"strconv"
	"strings"
)

// Record is a representation of the fields of a single VCF record
type Record struct {
	// The 8 required fields
	Chrom     string
	Pos       int
	ID        string
	Ref       string
	Alt       string
	Qual      float64
	Filter    string
	Info      map[string]string
	Format    []string
	Genotypes [][]string
	HasQual   bool
}

// NewRecord constructs a vcf.Record
func NewRecord() *Record {
	var result Record
	result.Info = make(map[string]string)
	return &result
}

// ParseRecord parses a string to a vcf.Record struct
func ParseRecord(s string) (*Record, error) {
	r := NewRecord()
	parts := strings.Split(s, "\t")
	if len(parts) < 8 {
		return r, fmt.Errorf("invalid VCF line")
	}
	r.Chrom = parts[0]
	pos, err := strconv.ParseInt(parts[1], 10, 64)
	if err != nil {
		return r, fmt.Errorf("invalid parsing of Pos")
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
	r.parseInfo(parts[7])

	// FORMAT field
	if len(parts) > 8 {
		r.Format = strings.Split(parts[8], ":")
	}

	// Sample genotypes
	if len(parts) > 9 {
		r.parseGenotypes(parts[9:])
	}
	return r, err
}

func (r *Record) parseGenotypes(fields []string) {
	for _, field := range fields {
		r.Genotypes = append(r.Genotypes, strings.Split(field, ":"))
	}
}

func (r *Record) parseInfo(s string) {

	for _, field := range strings.Split(strings.TrimSpace(s), ";") {
		if !strings.Contains(field, "=") { // Flag field
			r.Info[field] = field
			continue
		}
		idval := strings.SplitN(field, "=", 2)
		id := idval[0]
		val := idval[1]
		r.Info[id] = val
	}
}
