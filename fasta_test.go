package biofiles

import (
	"fmt"
	"strings"
)

func ExampleParseFasta() {
	var fastaExample = `>seq1 description of seq1
	ATGCGAGATAGATCATACTGAGCTCCCTACAGGGAATCA
	>seq2 desccription of seq2
	ATGCCCATGGACGACTATGACCCGAGCTACTA
	`
	recs := ParseFasta(strings.NewReader(fastaExample))
	fmt.Println(len(recs))
	fmt.Println(recs[0].ID, recs[1].ID)
	fmt.Println(recs[0].Description)
	fmt.Println(recs[1].Sequence)
	// Output:
	// 2
	// seq1 seq2
	// description of seq1
	// ATGCCCATGGACGACTATGACCCGAGCTACTA
}
