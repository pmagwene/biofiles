package vcf

import (
	"bufio"
	"fmt"
	"io"
	"strconv"
	"strings"
)

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
	Info      map[string]string
	Genotypes *GenotypeData
	HasQual   bool
}

// GenotypeData represents the genotype information in the genotype
// field of a VCF record
type GenotypeData struct {
	Format  string
	Samples []string
}

// MetaInformation represents metainformation associated with VCF file
type MetaInformation struct {
	Key   string
	Value string
}

// File represents a collection of VCF records and associated
// metainformation
type File struct {
	Fileformat      string
	MetaInformation []MetaInformation
	Records         []*Record
}

// ParseFile parses a VCF file, returning a vcf.File struct
func ParseFile(r io.Reader) (File, error) {
	var records []*Record
	var metainfo []MetaInformation

	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		// Empty line
		if len(line) == 0 {
			continue
		}
		// comment, header, or metainformation
		if strings.HasPrefix(line, "#") {
			// header line
			if strings.HasPrefix(line, "#CHROM") {
				continue
			}
			// metainformation
			if strings.HasPrefix(line, "##") {
				matches := strings.SplitN(line, "=", 2)
				if matches != nil {
					metainfo = append(metainfo,
						MetaInformation{matches[0][2:], matches[1]})
				}
				continue
			}
		}
		// line with fileds
		parts := strings.Split(line, "\t")
		if len(parts) < 8 {
			file := File{MetaInformation: metainfo, Records: records}
			return file, fmt.Errorf("invalid VCF line")
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
		r.Info = parseInfoFields(parts[7])
		if len(parts) >= 9 {
			r.Genotypes = parseGenotypeFields(parts[8:])
		}
		records = append(records, r)
	}
	file := File{MetaInformation: metainfo, Records: records}
	return file, nil
}

func parseMetaInfo(m MetaInformation) map[string]map[string]string {
	typemap := make(map[string]map[string]string)

	return typemap
}

func parseInfoFields(s string) map[string]string {
	infomap := make(map[string]string)

	fields := strings.Split(strings.TrimSpace(s), ";")
	for _, field := range fields {
		keyvalue := strings.Split(field, "=")
		if len(keyvalue) < 2 {
			continue
		}
		key := keyvalue[0]
		value := keyvalue[1]
		infomap[key] = value
	}

	return infomap
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
