package core

import (
<<<<<<< HEAD
	"github.com/araddon/gou"
=======
	u "github.com/araddon/gou"
>>>>>>> upstream/master
	"testing"
)

func TestSearchRequest(t *testing.T) {
	qry := map[string]interface{}{
		"query": map[string]interface{}{
			"wildcard": map[string]string{"actor": "a*"},
		},
	}
	out, err := SearchRequest(true, "github", "", qry, "")
	//log.Println(out)
<<<<<<< HEAD
	gou.Assert(&out != nil && err == nil, t, "Should get docs")
	gou.Assert(out.Hits.Len() == 10, t, "Should have 10 docs but was %v", out.Hits.Len())
	gou.Assert(out.Hits.Total == 589, t, "Should have 589 hits but was %v", out.Hits.Total)
=======
	Assert(&out != nil && err == nil, t, "Should get docs")
	Assert(out.Hits.Len() == 10, t, "Should have 10 docs but was %v", out.Hits.Len())
	Assert(u.CloseInt(out.Hits.Total, 588), t, "Should have 588 hits but was %v", out.Hits.Total)
>>>>>>> upstream/master
}
