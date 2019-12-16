package fasta

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"strings"
)

// Record is a representation of a single
// FASTA formatted sequence record
type Record struct {
	ID          string
	Description string
	Sequence    string
}

func (r *Record) String() string {
	i := len(r.Sequence)
	fmtstr := ">%s %s\n%s\n"
	if i > 10 {
		i = 10
		fmtstr = ">%s %s\n%s...\n"
	}
	return fmt.Sprintf(fmtstr,
		r.ID, r.Description, r.Sequence[:i])
}

// Parse parses a FASTA file into its
// constituent records, returned as a slice
func Parse(r io.Reader) []*Record {

	var records []*Record
	var sequence strings.Builder
	var currentRecord *Record = new(Record)

	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		// Empty line
		if len(line) == 0 {
			continue
		}
		// Comment line
		if strings.HasPrefix(line, ";") {
			continue
		}
		// Start new record
		if strings.HasPrefix(line, ">") {
			if len(currentRecord.ID) > 0 {
				currentRecord.Sequence = sequence.String()
				records = append(records, currentRecord)
			}
			currentRecord = new(Record)
			sequence.Reset()
			fields := strings.Fields(line[1:])
			currentRecord.ID = fields[0]
			if len(fields) > 1 {
				currentRecord.Description =
					strings.Join(fields[1:], " ")
			}
		} else {
			sequence.Grow(len(line) * 2)
			sequence.WriteString(line)
		}
	}
	// Process last record
	if len(currentRecord.ID) > 0 {
		currentRecord.Sequence = sequence.String()
		records = append(records, currentRecord)
	}
	return records
}

func wrapString(s string, l int) string {
	var ws strings.Builder
	var i int
	for i = 0; len(s[i:]) > l; i += l {
		ws.WriteString(s[i : i+l])
		ws.WriteByte('\n')
	}
	ws.WriteString(s[i:]) // Process last bit of string
	return ws.String()
}

// Write writes string represeentations of single fasta.Record
// to the given io.Writer
func (r *Record) Write(w io.Writer) {
	var b bytes.Buffer
	fmt.Fprintf(&b, ">%s %s\n", r.ID, r.Description)
	b.WriteString(wrapString(r.Sequence, 80))
	b.WriteByte('\n')
	w.Write(b.Bytes())
}

// WriteAll writes string representations of slice of Records
// to the given io.Writer
func WriteAll(recs []*Record, w io.Writer) {
	for _, rec := range recs {
		rec.Write(w)
	}
}

// ToFastaDict converts sequence of fasta.Record to map of
// Record indexed by ID
func ToFastaDict(recs []*Record) map[string]*Record {
	dict := make(map[string]*Record)
	for _, rec := range recs {
		dict[rec.ID] = rec
	}
	return dict
}
