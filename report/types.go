package report

import (
	"github.com/fafeitsch/open-callopticum/cdrcsv"
	"regexp"
)

type CdrColumn interface {
	get() column
}
type column string

const (
	AccountCode column = "AccountCode"
	Src         column = "Src"
	Dst         column = "Dst"
	Dcontext    column = "Dcontext"
	CallerId    column = "CallerId"
	Channel     column = "Channel"
	DstChannel  column = "DstChannel"
	LastApp     column = "LastApp"
	LastData    column = "LastData"
	Start       column = "Start"
	Answer      column = "Answer"
	End         column = "End"
	Duration    column = "Duration"
	Billsec     column = "Billsec"
	Disposition column = "Disposition"
	AmaFlag     column = "AmaFlag"
	Userfield   column = "Userfield"
	UniqueId    column = "UniqueId"
)

type Matcher interface {
	MatchRecord(*cdrcsv.Record) bool
}

type RegexMatcher struct {
	regex    regexp.Regexp
	provider func(*cdrcsv.Record) string
}

func (r *RegexMatcher) MatchRecord(record *cdrcsv.Record) bool {
	text := r.provider(record)
	return r.regex.MatchString(text)
}

type AndMatcher struct {
	left  Matcher
	right Matcher
}

func (a *AndMatcher) MatchRecord(record *cdrcsv.Record) bool {
	return a.left.MatchRecord(record) && a.right.MatchRecord(record)
}

type OrMatcher struct {
	left  Matcher
	right Matcher
}

func (a *OrMatcher) MatchRecord(record *cdrcsv.Record) bool {
	return a.left.MatchRecord(record) || a.right.MatchRecord(record)
}

type CountingsDefinition struct {
	Name                string
	DisplayName         string
	Formula             Matcher
	ExcludeOtherMatches bool
}

type ReportDefinition struct {
	Countings []CountingsDefinition
}
