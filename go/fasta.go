package biofiles

import (
	"bufio"
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

	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if len(line) == 0 {
			continue
		}
		if string(line[0]) == ">" {
			if len(currentRecord.ID) > 0 {
				records = append(records, currentRecord)
			}
			currentRecord = FastaRecord{}
			fields := strings.Fields(line[1:])
			currentRecord.ID = fields[0]
			if len(fields) > 1 {
				currentRecord.Description =
					strings.Join(fields[1:], " ")
			}
		} else {
			currentRecord.Sequence =
				strings.Join([]string{currentRecord.Sequence, line}, "")
		}
	}
	if len(currentRecord.ID) > 0 {
		records = append(records, currentRecord)
	}
	return records
}
