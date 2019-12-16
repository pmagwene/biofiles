package gff

import (
	"fmt"
	"strings"
)

func ExampleParseGFFRecord() {
	var oneGFF = "chrI	SGD	telomere	1	801	.	-	.	ID=TEL01L;Name=TEL01L"
	rec, _ := ParseGFFRecord(oneGFF)
	fmt.Println(rec.Start)
	fmt.Println(rec.End)
	fmt.Println(rec.ID)
	// Output:
	// 1
	// 801
	// TEL01L
}

func ExampleParseGFF() {
	var manyGFF = `
chrI	SGD	chromosome	1	230218	.	.	.	ID=chrI;dbxref=NCBI:NC_001133;Name=chrI
chrI	SGD	telomere	1	801	.	-	.	ID=TEL01L;Name=TEL01L
chrI	SGD	gene	335	649	.	+	.	ID=YAL069W;Name=YAL069W;
chrI	SGD	CDS	335	649	.	+	0	Parent=YAL069W_mRNA;Name=YAL069W_CDS;
chrI	SGD	mRNA	335	649	.	+	.	ID=YAL069W_mRNA;Name=YAL069W_mRNA;Parent=YAL069W
`
	recs, _ := ParseGFF(strings.NewReader(manyGFF))
	fmt.Println(len(recs))
	fmt.Println(recs[0].Type)
	fmt.Println(recs[len(recs)-1].Type)
	fmt.Println(recs[len(recs)-1].Parent)
	// Output:
	// 5
	// chromosome
	// mRNA
	// YAL069W
}
