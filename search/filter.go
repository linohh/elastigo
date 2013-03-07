package search

import (
	"encoding/json"
	"fmt"

	. "github.com/araddon/gou"
)

var (
	_         = DEBUG
	ANDCLAUSE = LogicClause("and")
)

// A logical clause  [and, or, not]
type LogicClause string

// Filter clause is either a LogicClause or FilterOp
type FilterClause interface {
	String() string
}

// A wrapper to allow for custom serialization
//     A filter is a nested set of operators (not,and,or) and the
//     filter clauses that are contained in them
type FilterWrap struct {
	// The complete map of logic operators to their nested filters
	filtermap map[LogicClause][]interface{}
	// The current LogicClause, can be none (which = "and")
	curClause LogicClause
}

func NewFilterWrap() *FilterWrap {
	return &FilterWrap{curClause: "and"}
}

func (f *FilterWrap) String() string {
	return fmt.Sprintf(`fopv: %d:%v`, len(f.filtermap), f.filtermap)
}

// logic associated with adding filters.   
func (f *FilterWrap) addFilters(filterList []interface{}) {
	// Default clause is "and"
	if len(f.curClause) == 0 {
		f.curClause = LogicClause("and")
	}
	if len(f.filtermap) == 0 {
		f.filtermap = make(map[LogicClause][]interface{})
	}
	for _, filterClause := range filterList {
		switch clauseOp := filterClause.(type) {
		case LogicClause, string:
			f.curClause = LogicClause(filterClause.(string))
		case *QueryDsl:
			if len(f.filtermap[f.curClause]) == 0 {
				f.filtermap[f.curClause] = make([]interface{}, 0)
			}
			f.filtermap[f.curClause] = append(f.filtermap[f.curClause], map[string]*QueryDsl{"query": clauseOp})
		case *FilterOp:
			if len(f.filtermap[f.curClause]) == 0 {
				f.filtermap[f.curClause] = make([]interface{}, 0)
			}
			if clauseOp.not {
				m := map[string]*FilterOp{"not": clauseOp}
				f.filtermap[f.curClause] = append(f.filtermap[f.curClause], m)
			} else {
				f.filtermap[f.curClause] = append(f.filtermap[f.curClause], clauseOp)
			}
		default:
			Logf(ERROR, "Unkown Filter Clause? %v", clauseOp)
		}
	}
}

// Custom marshalling to support the query dsl 
func (f *FilterWrap) MarshalJSON() ([]byte, error) {

	var root interface{}
	_, hasAnd := f.filtermap[ANDCLAUSE]

	// If we only have one clause and its an and we don't need to
	// send as that is default for elasticsearch
	if len(f.filtermap) == 1 && hasAnd && len(f.filtermap[ANDCLAUSE]) == 1 {
		root = f.filtermap[ANDCLAUSE]
		return json.Marshal(root)
	}
	return json.Marshal(f.filtermap)
}

/*
	"filter": {
		"range": {
		  "@timestamp": {
		    "from": "2012-12-29T16:52:48+00:00",
		    "to": "2012-12-29T17:52:48+00:00"
		  }
		}
	}
	"filter": {
	    "missing": {
	        "field": "repository.name"
	    }
	}

	"filter" : {
	    "terms" : {
	        "user" : ["kimchy", "elasticsearch"],
	        "execution" : "bool",
	        "_cache": true
	    }
	}

	"filter" : {
	    "term" : { "user" : "kimchy"}
	}

	"filter" : {
	    "and" : [
	        {
	            "range" : { 
	                "postDate" : { 
	                    "from" : "2010-03-01",
	                    "to" : "2010-04-01"
	                }
	            }
	        },
	        {
	            "prefix" : { "name.second" : "ba" }
	        }
	    ]
	}

{
  "filter": 
    { "and" : [
      { "term": { "actor_attributes.location": "portland" } } 
      , { "not": 
             { "filter": { "term": {"repository.has_wiki": true} } }
      }
    ]
  }
}
*/

// Filter Operation
//
//   Filter().Term("user","kimchy")
//    
//   // we use variadics to allow n arguments, first is the "field" rest are values
//   Filter().Terms("user", "kimchy", "elasticsearch")
// 
//   Filter().Exists("repository.name")
//
func Filter() *FilterOp {
	return &FilterOp{}
}

type FilterOp struct {
	curField   string
	not        bool                         // The not operator
	TermsMap   map[string][]interface{}     `json:"terms,omitempty"`
	Range      map[string]map[string]string `json:"range,omitempty"`
	Exist      map[string]string            `json:"exists,omitempty"`
	MissingVal map[string]string            `json:"missing,omitempty"`
}

// A range is a special type of Filter operation
//
//    Range().Exists("repository.name")
func Range() *FilterOp {
	return &FilterOp{Range: make(map[string]map[string]string)}
}

func (f *FilterOp) Field(fld string) *FilterOp {
	f.curField = fld
	if _, ok := f.Range[fld]; !ok {
		m := make(map[string]string)
		f.Range[fld] = m
	}
	return f
}

// Filter Terms
//
//   Filter().Terms("user","kimchy")
//    
//   // we use variadics to allow n arguments, first is the "field" rest are values
//   Filter().Terms("user", "kimchy", "elasticsearch")
//
func (f *FilterOp) Terms(field string, values ...interface{}) *FilterOp {
	if len(f.TermsMap) == 0 {
		f.TermsMap = make(map[string][]interface{})
	}
	for _, val := range values {
		f.TermsMap[field] = append(f.TermsMap[field], val)
	}

	return f
}
func (f *FilterOp) From(from string) *FilterOp {
	f.Range[f.curField]["from"] = from
	return f
}
func (f *FilterOp) To(to string) *FilterOp {
	f.Range[f.curField]["to"] = to
	return f
}
func (f *FilterOp) Exists(name string) *FilterOp {
	f.Exist = map[string]string{"field": name}
	return f
}
func (f *FilterOp) Missing(name string) *FilterOp {
	f.MissingVal = map[string]string{"field": name}
	return f
}
func (f *FilterOp) Not() *FilterOp {
	f.not = true
	return f
}

// Add another Filterop, "combines" two filter ops into one
func (f *FilterOp) Add(fop *FilterOp) *FilterOp {
	// TODO, this is invalid, refactor
	if len(fop.Exist) > 0 {
		f.Exist = fop.Exist
	}
	if len(fop.MissingVal) > 0 {
		f.MissingVal = fop.MissingVal
	}
	if len(fop.Range) > 0 {
		f.Range = fop.Range
	}
	return f
}

/*
// Custom marshalling to support the query dsl 
func (f *FilterOp) MarshalJSON() ([]byte, error) {
	// Ok, hope this doesn't break with go 1.04 as it is probably 
	// capitializing on a bug:   to prevent this from becoming a recursive lock
	// we derefrence, and pointers don't get passed here?
	if f.not {
		m := map[string]FilterOp{"not": *f}
		return json.Marshal(&m)
	}
	return json.Marshal(*f)
}
*/
