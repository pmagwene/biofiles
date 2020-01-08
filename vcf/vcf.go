package vcf

import (
	"bufio"
	"io"
	"regexp"
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

var keyvalRe = regexp.MustCompile(`(?P<key>\w+)=(?P<value>["].+?["]|[^,]+?)(?:$|,)`)

func splitKeyEqualValueString(s string) (string, string) {
	matches := keyvalRe.FindStringSubmatch(s)
	if len(matches) != 3 {
		return "", ""
	}
	return matches[1], matches[2]
}

// ParseFile parses a VCF file, returning a vcf.Table struct
func ParseFile(r io.Reader) (*Table, error) {
	table := NewTable()

	//adjust the capacity to your need (max characters in line)
	const maxCapacity = 512 * 1024
	buf := make([]byte, maxCapacity)

	scanner := bufio.NewScanner(r)
	scanner.Buffer(buf, maxCapacity)

	var ct int = 0
	for scanner.Scan() {
		ct++
		line := strings.TrimSpace(scanner.Text())

		// Empty line or header line
		if len(line) == 0 || strings.HasPrefix(line, "#CHROM") {
			continue
		}
		// comment or metainformation
		if strings.HasPrefix(line, "#") {
			if strings.HasPrefix(line, "##") {
				meta, err := ParseMetadata(line)
				if err != nil {
					return table, err
				}
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
				default:
					table.Metadata = append(table.Metadata, meta)
				}
			}
			continue
		}
		r, err := ParseRecord(line)
		if err != nil {
			return table, err
		}
		table.Records = append(table.Records, r)
	}
	err := scanner.Err()
	return table, err
}
