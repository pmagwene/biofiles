package gff

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"net/url"
	"strconv"
	"strings"

	"github.com/pmagwene/biofiles/fasta"
)

/*
Description of GFF3 file format from: http://gmod.org/wiki/GFF3

GFF3 Format

GFF3 format is a flat tab-delimited file. The first line of the file is a comment that identifies the file format and version. This is followed by a series of data lines, each one of which corresponds to an annotation.Here is a miniature GFF3 file:

##gff-version 3
ctg123  .  exon  1300  1500  .  +  .  ID=exon00001
ctg123  .  exon  1050  1500  .  +  .  ID=exon00002
ctg123  .  exon  3000  3902  .  +  .  ID=exon00003
ctg123  .  exon  5000  5500  .  +  .  ID=exon00004
ctg123  .  exon  7000  9000  .  +  .  ID=exon00005

The ##gff-version 3 line is required and must be the first line of the file. It introduces the annotation section of the file.

The 9 columns of the annotation section are as follows:

Column 1: "seqid"

    The ID of the landmark used to establish the coordinate system for the current feature. IDs may contain any characters, but must escape any characters not in the set [a-zA-Z0-9.:^*$@!+_?-|]. In particular, IDs may not contain unescaped whitespace and must not begin with an unescaped ">".

    To escape a character in this, or any of the other GFF3 fields, replace it with the percent sign followed by its hexadecimal representation. For example, ">" becomes "%E3". See URL Encoding (or: 'What are those "%20" codes in URLs?') for details.

Column 2: "source"

    The source is a free text qualifier intended to describe the algorithm or operating procedure that generated this feature. Typically this is the name of a piece of software, such as "Genescan" or a database name, such as "Genbank." In effect, the source is used to extend the feature ontology by adding a qualifier to the type creating a new composite type that is a subclass of the type in the type column. It is not necessary to specify a source. If there is no source, put a "." (a period) in this field.

Column 3: "type"

    The type of the feature (previously called the "method"). This is constrained to be either: (a) a term from the "lite" sequence ontology, SOFA; or (b) a SOFA accession number. The latter alternative is distinguished using the syntax SO:000000. This field is required.

Columns 4 & 5: "start" and "end"

    The start and end of the feature, in 1-based integer coordinates, relative to the landmark given in column 1. Start is always less than or equal to end.

    For zero-length features, such as insertion sites, start equals end and the implied site is to the right of the indicated base in the direction of the landmark. These fields are required.

Column 6: "score"

    The score of the feature, a floating point number. As in earlier versions of the format, the semantics of the score are ill-defined. It is strongly recommended that E-values be used for sequence similarity features, and that P-values be used for ab initio gene prediction features. If there is no score, put a "." (a period) in this field.

Column 7: "strand"

    The strand of the feature. + for positive strand (relative to the landmark), - for minus strand, and . for features that are not stranded. In addition, ? can be used for features whose strandedness is relevant, but unknown.

Column 8: "phase"

    For features of type "CDS", the phase indicates where the feature begins with reference to the reading frame. The phase is one of the integers 0, 1, or 2, indicating the number of bases that should be removed from the beginning of this feature to reach the first base of the next codon. In other words, a phase of "0" indicates that the next codon begins at the first base of the region described by the current line, a phase of "1" indicates that the next codon begins at the second base of this region, and a phase of "2" indicates that the codon begins at the third base of this region. This is NOT to be confused with the frame, which is simply start modulo 3. If there is no phase, put a "." (a period) in this field.

    For forward strand features, phase is counted from the start field. For reverse strand features, phase is counted from the end field.

    The phase is required for all CDS features.

Column 9: "attributes"

    A list of feature attributes in the format tag=value. Multiple tag=value pairs are separated by semicolons. URL escaping rules are used for tags or values containing the following characters: ",=;". Spaces are allowed in this field, but tabs must be replaced with the %09 URL escape. This field is not required.

Column 9 Tags

Column 9 tags have predefined meanings:

ID
    Indicates the unique identifier of the feature. IDs must be unique within the scope of the GFF file.

Name
    Display name for the feature. This is the name to be displayed to the user. Unlike IDs, there is no requirement that the Name be unique within the file.

Alias
    A secondary name for the feature. It is suggested that this tag be used whenever a secondary identifier for the feature is needed, such as locus names and accession numbers. Unlike ID, there is no requirement that Alias be unique within the file.

Parent
    Indicates the parent of the feature. A parent ID can be used to group exons into transcripts, transcripts into genes, and so forth. A feature may have multiple parents. Parent can *only* be used to indicate a partof relationship.

Target
    Indicates the target of a nucleotide-to-nucleotide or protein-to-nucleotide alignment. The format of the value is "target_id start end [strand]", where strand is optional and may be "+" or "-". If the target_id contains spaces, they must be escaped as hex escape %20.

Gap
    The alignment of the feature to the target if the two are not collinear (e.g. contain gaps). The alignment format is taken from the CIGAR format described in the Exonerate documentation. http://cvsweb.sanger.ac.uk/cgi-bin/cvsweb.cgi/exonerate?cvsroot=Ensembl). See the GFF3 specification for more information.

Derives_from
    Used to disambiguate the relationship between one feature and another when the relationship is a temporal one rather than a purely structural "part of" one. This is needed for polycistronic genes. See the GFF3 specification for more information.

Note
    A free text note.

Dbxref
    A database cross reference. See the GFF3 specification for more information.

Ontology_term
    A cross reference to an ontology term. See the GFF3 specification for more information.

Multiple attributes of the same type are indicated by separating the values with the comma "," character, as in:

Parent=AF2312,AB2812,abc-3

Note that attribute names are case sensitive. "Parent" is not the same as "parent".

All attributes that begin with an uppercase letter are reserved for later use. Attributes that begin with a lowercase letter can be used freely by applications. You can stash any semi-structured data into the database by using one or more unreserved (lowercase) tags.

*/

// Record represents the information in a GFF3 record
type Record struct {

	// The nine standard GFF fields
	SeqID      string
	Source     string
	Type       string
	Start      int
	End        int
	Score      float64
	Strand     string
	Phase      string
	Attributes map[string]string

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
	Sequence    string
	IsGene      bool
	ScoreExists bool
	Children    []*Record
	Cds         []*Record
	Parts       []*Record
	Exons       []*Record
	Introns     []*Record
}

func (r *Record) String() string {
	return fmt.Sprintf("(%s, %s, %s, %d, %d, %s)",
		r.ID, r.Type, r.SeqID, r.Start, r.End, r.Strand)
}

// ParseRecord turns a string into a single GFF record
func ParseRecord(s string) (*Record, error) {

	var r Record

	parts := strings.Split(strings.TrimSpace(s), "\t")
	if len(parts) != 9 {
		return &r, fmt.Errorf("Invalid GFF record string")
	}
	r.SeqID = parts[0]
	r.Source = parts[1]
	r.Type = parts[2]
	start, err := strconv.ParseInt(parts[3], 10, 0)
	if err == nil {
		r.Start = int(start)
	}
	end, err := strconv.ParseInt(parts[4], 10, 0)
	if err == nil {
		r.End = int(end)
	}
	score, err := strconv.ParseFloat(parts[5], 64)
	if err == nil {
		r.Score = score
		r.ScoreExists = true
	}
	r.Strand = parts[6]
	r.Phase = parts[7]

	if r.Type == "gene" {
		r.IsGene = true
	}

	// parse attributes field, setting standard attribute
	r.Attributes = parseAttributes(parts[8])
	r.ID = r.Attributes["ID"]
	r.Name = r.Attributes["Name"]
	r.Parent = r.Attributes["Parent"]
	r.Target = r.Attributes["Target"]
	r.Gap = r.Attributes["Gap"]
	r.DerivesFrom = r.Attributes["Derives_from"]
	r.Note = r.Attributes["Note"]
	r.Dbxref = r.Attributes["Dbxref"]
	r.OntologyTerm = r.Attributes["Ontology_term"]

	return &r, nil

}

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
		unescaped, _ := url.PathUnescape(attribval[1])
		attribDict[attribval[0]] = unescaped
	}
	return attribDict
}

func buildAttributeStr(rec *Record) string {
	var fmtstr string
	var attribstr strings.Builder
	stdattribs := [...]string{"ID", "Name", "Parent", "Target",
		"Gap", "Derives_from", "Note", "Dbxref", "Ontology_term"}

	// Deal with standard attributes in given order
	for _, attrib := range stdattribs {
		val := rec.Attributes[attrib]
		if len(val) < 1 {
			continue
		}
		if attribstr.Len() > 0 {
			fmtstr = ";%s=%s"
		} else {
			fmtstr = "%s=%s"
		}
		attribstr.WriteString(
			fmt.Sprintf(fmtstr, attrib, url.PathEscape(val)))

	}
	// Deal with other attributes
	for attrib, val := range rec.Attributes {
		switch attrib {
		case "ID", "Name", "Parent", "Target", "Gap",
			"Derives_from", "Note", "Dbxref", "Ontology_term":
			continue
		default:
			if attribstr.Len() > 0 {
				fmtstr = ";%s=%s"
			} else {
				fmtstr = "%s=%s"
			}
			attribstr.WriteString(
				fmt.Sprintf(fmtstr, attrib, url.PathEscape(val)))
		}
	}
	return attribstr.String()
}

// Populate the Children field of a slice of Record structs
func populateChildren(recs []*Record) {
	var tbl = make(map[string]*Record)
	for i, rec := range recs {
		tbl[rec.ID] = recs[i]
	}
	for i, rec := range recs {
		if rec.Parent != "" {
			therec := recs[i]
			parent := tbl[therec.Parent]
			parent.Children = append(parent.Children, therec)
		}
	}
}

// ParseAll parses GFF records, returning a slice of
// *Record and a slice(possibly empty) with associated fasta.Records
func ParseAll(r io.Reader) ([]*Record, []*fasta.Record) {
	var records []*Record
	var isFasta bool
	var fastastr strings.Builder

	input := bufio.NewScanner(r)
	for input.Scan() {
		line := strings.TrimSpace(input.Text())
		if len(line) == 0 {
			continue
		}
		// Comment lines
		if strings.HasPrefix(line, "#") {
			// Marks beginning of FASTA section
			if strings.HasPrefix(line, "##FASTA") {
				isFasta = true
			}
			continue
		}
		// Process FASTA
		if isFasta {
			fastastr.WriteString(line)
			fastastr.WriteByte('\n')
			continue
		}
		// If not FASTA, process Record
		rec, err := ParseRecord(line)
		if err == nil {
			records = append(records, rec)
		}
	}
	populateChildren(records)

	var fastarecs []*fasta.Record
	if fastastr.Len() > 0 {
		fastarecs = fasta.ParseAll(strings.NewReader(fastastr.String()))
	}
	return records, fastarecs
}

// WriteRecord writes a single GFF record to the given Writer
func WriteRecord(r *Record, w io.Writer) {
	var b bytes.Buffer
	scorestr := "."
	if r.ScoreExists {
		scorestr = fmt.Sprintf("%f", r.Score)
	}
	attribstr := buildAttributeStr(r)
	fmtstr := "%s\t%s\t%s\t%d\t%d\t%s\t%s\t%s\t%s\n"
	recstr := fmt.Sprintf(fmtstr,
		r.SeqID, r.Source, r.Type, r.Start, r.End,
		scorestr, r.Strand, r.Phase, attribstr)
	b.WriteString(recstr)
}

// WriteAll writes string representations of GFF records, and
// optional associated Fasta records, to given Writer interface
func WriteAll(recs []*Record, w io.Writer) {
	for _, rec := range recs {
		WriteRecord(rec, w)
	}
}

// WriteFastaSection appends an optional Fasta section to a GFF file.
// WriteFastaSection should be called after gff.Write or gff.WriteAll
func WriteFastaSection(fastarecs []*fasta.Record, w io.Writer) {
	w.Write([]byte("##FASTA\n"))
	fasta.WriteAll(fastarecs, w)
}

// ToFastaRecord generates a fasta.Record that corresponds to the given
// GFF Record. lwindow and rwindow parameters facilitate specification of a
// sequence window around the feature.
func (r Record) ToFastaRecord(fastadict map[string]*fasta.Record,
	lwindow int, rwindow int) *fasta.Record {
	target := fastadict[r.SeqID]
	wstart := r.Start - 1 - lwindow
	wend := r.End - 1 + rwindow + 1 // additional +1 to include rightmost coord
	if wstart < 0 {
		wstart = 0
	}
	if wend > len(target.Sequence)-1 {
		wend = len(target.Sequence) - 1
	}
	idstr := r.ID
	if len(idstr) < 1 {
		idstr = fmt.Sprintf("%s_%s_%s_%d_%d",
			r.SeqID, r.Source, r.Type, r.Start, r.End)
	}
	return &fasta.Record{
		ID: idstr,
		Description: fmt.Sprintf("%s:%d..%d",
			r.SeqID, wstart+1, wend),
		Sequence: target.Sequence[wstart:wend]}
}
