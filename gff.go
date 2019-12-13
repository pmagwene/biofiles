package biofiles

import (
	"bufio"
	"fmt"
	"io"
	"net/url"
	"regexp"
	"strconv"
	"strings"
)

// GFFRecord represents the information in a GFF3 record
type GFFRecord struct {

	// The nine standard GFF fields
	SeqID      string
	Source     string
	Type       string
	Start      uint
	End        uint
	Score      float64
	Strand     string
	Phase      string
	Attributes string

	// the ten standard reserved attributes
	ID           string
	Name         string
	Alias        string
	Parent       string
	Target       string
	Gap          string
	DerivesFrom  string
	Note         string
	Dbxref       string
	OntologyTerm string

	// some attributes I've added for convenience
	AttributeDict map[string]string
	Sequence      string
	IsGene        bool
	Children      []*GFFRecord
	Cds           []*GFFRecord
	Parts         []*GFFRecord
	Exons         []*GFFRecord
	Introns       []*GFFRecord
}

func (f GFFRecord) String() string {
	return fmt.Sprintf("(%s, %s, %s, %d, %d, %s )",
		f.ID, f.Type, f.SeqID, f.Start, f.End, f.Strand)
}

// ParseGFFRecord turns a string into a single GFF record
func ParseGFFRecord(s string) (GFFRecord, error) {

	var f GFFRecord

	parts := strings.Split(strings.TrimSpace(s), "\t")
	if len(parts) != 9 {
		return f, fmt.Errorf("Invalid GFF record string")
	}
	f.SeqID = parts[0]
	f.Source = parts[1]
	f.Type = parts[2]
	start, err := strconv.ParseUint(parts[3], 10, 0)
	if err == nil {
		f.Start = uint(start)
	}
	end, err := strconv.ParseUint(parts[4], 10, 0)
	if err == nil {
		f.End = uint(end)
	}
	score, err := strconv.ParseFloat(parts[5], 64)
	if err == nil {
		f.Score = score
	}
	f.Strand = parts[6]
	f.Phase = parts[7]
	f.Attributes = parts[8]

	if f.Type == "gene" {
		f.IsGene = true
	}

	// parse attributes field, setting standard attribute
	f.AttributeDict = parseAttributes(parts[8])
	f.ID = f.AttributeDict["ID"]
	f.Name = f.AttributeDict["Name"]
	f.Parent = f.AttributeDict["Parent"]
	f.Target = f.AttributeDict["Target"]
	f.Gap = f.AttributeDict["Gap"]
	f.DerivesFrom = f.AttributeDict["Derives_from"]
	f.Note = f.AttributeDict["Note"]
	f.Dbxref = f.AttributeDict["Dbxref"]
	f.OntologyTerm = f.AttributeDict["Ontology_term"]

	return f, nil

}

var attributeRe = regexp.MustCompile(`(\S+)=(.+)`)

func parseAttributes(s string) map[string]string {
	// Parse the attributes field (column 9) of a GFF3 record.
	attribDict := make(map[string]string)
	if len(s) == 0 {
		return attribDict
	}
	if s == "." {
		return attribDict
	}
	fields := strings.Split(s, ";")
	for _, val := range fields {
		attribval := strings.SplitN(val, "=", 2)
		if len(attribval) < 2 {
			continue
		}
		unescaped, _ := url.QueryUnescape(attribval[1])
		attribDict[attribval[0]] = unescaped
		// matches := attributeRe.FindStringSubmatch(val)
		// if matches != nil {
		// 	// matches[0] is the entire match
		// 	// matches[1]... are the submatches
		// 	unescaped, _ := url.QueryUnescape(matches[2])
		// 	attribDict[matches[1]] = unescaped
		// }
	}
	return attribDict
}

// Populate the Children field of a slice of GFFRecord structs
func populateChildren(recs []GFFRecord) {
	var tbl = make(map[string]*GFFRecord)
	for i, rec := range recs {
		tbl[rec.ID] = &recs[i]
	}
	for i, rec := range recs {
		if rec.Parent != "" {
			therec := &recs[i]
			parent := tbl[therec.Parent]
			parent.Children = append(parent.Children, therec)
		}
	}
}

// ParseGFF parses GFF records, returning a slice of
// GFFRecord and a string (possibly empty) with associated Fasta file
func ParseGFF(r io.Reader) ([]GFFRecord, []FastaRecord) {
	var records []GFFRecord
	var isFasta bool
	var fastaStr strings.Builder

	input := bufio.NewScanner(r)
	for input.Scan() {
		line := strings.TrimSpace(input.Text())
		if len(line) == 0 {
			continue
		}
		// Comment lines
		if strings.HasPrefix(line, "#") {
			if strings.HasPrefix(line, "##FASTA") { // Marks beginning of FASTA section
				isFasta = true
			}
			continue
		}
		// Process FASTA
		if isFasta {
			fastaStr.WriteString(line)
			fastaStr.WriteByte('\n')
			continue
		}
		// If not FASTA, process GFFRecord
		rec, err := ParseGFFRecord(line)
		if err == nil {
			records = append(records, rec)
		}
	}
	populateChildren(records)

	var fastaRecs []FastaRecord
	if fastaStr.Len() > 0 {
		fastaRecs = ParseFasta(strings.NewReader(fastaStr.String()))
	}
	return records, fastaRecs
}
