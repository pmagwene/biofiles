package vcf

// type InfoField struct {
// 	ID          string
// 	Description string
// 	Number      string
// 	Type        string
// 	Other       map[string]string
// }

// func CreateConversionFunc(i InfoField) func(s []string) (interface{}, error) {
// 	var f func(s []string) (interface{}, error)
// 	switch i.Type {
// 	case "Integer":
// 		f = func(s []string) (interface{}, error) {
// 			var result []int64
// 			for _, each := range s {
// 				val, err := strconv.ParseInt(each, 10, 64)
// 				if err != nil {
// 					return result, err
// 				}
// 				result = append(result, val)
// 			}
// 			return result, nil
// 		}
// 	case "Float":
// 		f = func(s []string) (interface{}, error) {
// 			var result []float64
// 			for _, each := range s {
// 				val, err := strconv.ParseFloat(each, 64)
// 				if err != nil {
// 					return result, err
// 				}
// 				result = append(result, val)
// 			}
// 			return result, nil
// 		}
// 	case "Char":
// 		f = func(s []string) (interface{}, error) {
// 			return s, nil
// 		}
// 	case "String":
// 		f = func(s []string) (interface{}, error) {
// 			return s, nil
// 		}
// 	case "Flag":
// 		f = func(s []string) (interface{}, error) {
// 			return true, nil
// 		}
// 	default:
// 		panic("Invalid type!")
// 	}
// 	return f
// }

// // InfoAA represents the Ancestral allele
// type InfoAA struct {
// 	value string
// }

// // InfoAC represents allele counts for each alternate allele (not including ref)
// type InfoAC struct {
// 	value []int64
// }

// // InfoAD represents total read depth for each allele
// type InfoAD struct {
// 	value []int64
// }

// // InfoADF represents the read depth for each alternate allele on the
// // forward strand
// type InfoADF struct {
// 	value []int64
// }

// // InfoADR represents the read depth for each alternate allele on the
// // reverse strand
// type InfoADR struct {
// 	value []int64
// }

// //InfoAF represents the allele frequency for each allele
// type InfoAF struct {
// 	value []float64
// }

// //InfoAN represents the total number of alleles in called genotypes
// type InfoAN struct {
// 	value int64
// }

// // InfoBQ represents RMS base quality
// type InfoBQ struct {
// 	value float64
// }

// // InfoCIGAR represents a CIGAR  string describing how to align
// // an alternate allele to the reference allele
// type InfoCIGAR struct {
// 	value []string
// }

// // InfoDB represents whether variant is found in dbSNP
// type InfoDB struct {
// 	value bool
// }

// // InfoDP represents the combined depth across all the samples
// type InfoDP struct {
// 	value int64
// }

// // InfoEND represents the end position on CHROM
// type InfoEND struct {
// 	value int64
// }

// // InfoH2 represents whether variant is in HapMap2
// type InfoH2 struct {
// 	value bool
// }

// // InfoH3 represents whether variant is in HapMap3
// type InfoH3 struct {
// 	value bool
// }

// // InfoMQ represents RMS mapping quality
// type InfoMQ struct {
// 	value float64
// }

// // InfoMQ0 represents number of MQ == 0 reads
// type InfoMQ0 struct {
// 	value int64
// }

// // InfoNS represents the number of samples with data
// type InfoNS struct {
// 	value int64
// }

// // InfoSB represents the strand bias
// type InfoSB struct {
// 	value [4]int64
// }

// // InfoSOMATIC represents whether variant is a somatic mutation
// type InfoSOMATIC struct {
// 	value bool
// }

// // InfoVALIDATED represents whether there has been follow-up experiments
// type InfoVALIDATED struct {
// 	value bool
// }

// // InfoType defines interface for Info fields
// type InfoType interface {
// 	isInfoType()
// }

// func (InfoAA) isInfoType()        {}
// func (InfoAC) isInfoType()        {}
// func (InfoAD) isInfoType()        {}
// func (InfoADF) isInfoType()       {}
// func (InfoADR) isInfoType()       {}
// func (InfoAF) isInfoType()        {}
// func (InfoAN) isInfoType()        {}
// func (InfoBQ) isInfoType()        {}
// func (InfoCIGAR) isInfoType()     {}
// func (InfoDB) isInfoType()        {}
// func (InfoDP) isInfoType()        {}
// func (InfoEND) isInfoType()       {}
// func (InfoH2) isInfoType()        {}
// func (InfoH3) isInfoType()        {}
// func (InfoMQ) isInfoType()        {}
// func (InfoMQ0) isInfoType()       {}
// func (InfoNS) isInfoType()        {}
// func (InfoSB) isInfoType()        {}
// func (InfoSOMATIC) isInfoType()   {}
// func (InfoVALIDATED) isInfoType() {}
