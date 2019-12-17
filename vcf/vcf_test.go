package vcf

import (
	"fmt"
	"strings"
)

var vcfstring string = `
##fileformat=VCFv4.2
##FORMAT=<ID=GT,Number=1,Type=Integer,Description="Genotype">
##FORMAT=<ID=GP,Number=G,Type=Float,Description="Genotype Probabilities">
##FORMAT=<ID=PL,Number=G,Type=Float,Description="Phred-scaled Genotype Likelihoods">
#CHROM	POS	ID	REF	ALT	QUAL	FILTER	INFO	FORMAT	SAMP001	SAMP002
20	1291018	rs11449	G	A	.	PASS	.	GT	0/0	0/1
20	2300608 rs84825 C	T	.	PASS	.	GT:GP	0/1:.	0/1:0.03,0.97,0
20	2301308 rs84823 T	G	.	PASS	.	GT:PL	./.:.	1/1:10,5,0
`

func ExampleParseFile() {
	f, _ := ParseFile(strings.NewReader(vcfstring))
	fmt.Println(f)
}
