package vcf

// GenoAD represents read depth for each allele
type GenoAD struct {
	value []int64
}

// GenoADF represents read depth for each allele on forward strand
type GenoADF struct {
	value []int64
}

// GenoADR represents read depth for each allele on reverse strand
type GenoADR struct {
	value []int64
}

// GenoDP represents total read depth
type GenoDP struct {
	value int64
}

// GenoEC represents expected alternate allele counts
type GenoEC struct {
	value []int64
}

// GenoFT represents a filter indicating if these genotype was "called"
type GenoFT struct {
	value string
}

// GenoGL represents genotype likelihoods
type GenoGL struct {
	value []float64
}

// GenoGP represents genotype posterior probabilities
type GenoGP struct {
	value []float64
}

// GenoGQ represents conditional genotype quality
type GenoGQ struct {
	value int64
}

// GenoGT represents called genotype
type GenoGT struct {
	value string
}

// GenoHQ represents haplotype quality
type GenoHQ struct {
	value [2]int64
}

// GenoMQ represents RMS mapping quality
type GenoMQ struct {
	value int64
}

// GenoPL represents Phred-scaled genotype likelihoods rounded to
// closest integer
type GenoPL struct {
	value int64
}

// GenoPQ represents phasing quality
type GenoPQ struct {
	value int64
}

// GenoPS represents phase set
type GenoPS struct {
	value int64
}
