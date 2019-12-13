package biofiles

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"strings"
)

// FastaRecord is a representation of a single
// FASTA formatted sequence record
type FastaRecord struct {
	ID          string
	Description string
	Sequence    string
}

// ParseFasta parses a FASTA file into its
// constituent records, returned as a slice
func ParseFasta(r io.Reader) []FastaRecord {

	var records []FastaRecord
	var currentRecord FastaRecord
	var sequence strings.Builder

	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if len(line) == 0 {
			continue
		}
		if strings.HasPrefix(line, ";") {
			continue
		}
		if strings.HasPrefix(line, ">") {
			if len(currentRecord.ID) > 0 {
				currentRecord.Sequence = sequence.String()
				records = append(records, currentRecord)
			}
			currentRecord = FastaRecord{}
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

// WriteFasta writes a slice of FastaRecords to the given writer
func WriteFasta(recs []FastaRecord, w io.Writer) {
	var recbytes bytes.Buffer
	for _, rec := range recs {
		recbytes.Reset()
		fmt.Fprintf(&recbytes, ">%s %s\n", rec.ID, rec.Description)
		recbytes.WriteString(wrapString(rec.Sequence, 80))
		recbytes.WriteByte('\n')
		w.Write(recbytes.Bytes())
	}
}
