package search

import (
	//"encoding/json"
	. "github.com/araddon/gou"
	"testing"
)

func TestFilters(t *testing.T) {
	// search for docs that are missing repository.name
	qry := Search("github").Filter(
		Filter().Exists("repository.name"),
	)
	out, err := qry.Result()
	Assert(err == nil, t, "should not have error")
	Assert(out.Hits.Len() == 10, t, "Should have 10 docs %v", out.Hits.Len())
	Assert(out.Hits.Total == 7305, t, "Should have 7305 total= %v", out.Hits.Total)

	qry = Search("github").Filter(
		Filter().Missing("repository.name"),
	)
	out, _ = qry.Result()
	Assert(out.Hits.Len() == 10, t, "Should have 10 docs %v", out.Hits.Len())
	Assert(out.Hits.Total == 306, t, "Should have 306 total= %v", out.Hits.Total)

	//actor_attributes: {type: "User",
	qry = Search("github").Filter(
		Filter().Terms("actor_attributes.location", "portland"),
	)
	out, _ = qry.Result()
	Debug(out)
	Assert(out.Hits.Len() == 10, t, "Should have 10 docs %v", out.Hits.Len())
	Assert(CloseInt(out.Hits.Total, 66), t, "Should have ~66 total= %v", out.Hits.Total)

	/*
		Should this be an AND by default?
	*/
	qry = Search("github").Filter(
		Filter().Terms("actor_attributes.location", "portland"),
		Filter().Terms("repository.has_wiki", true),
	)
	out, err = qry.Result()
	Debug(out)
	Assert(err == nil, t, "should not have error")
	Assert(out.Hits.Len() == 10, t, "Should have 10 docs %v", out.Hits.Len())
	Assert(CloseInt(out.Hits.Total, 42), t, "Should have 42 total= %v", out.Hits.Total)

	// NOW, lets try with two query calls instead of one
	qry = Search("github").Filter(
		Filter().Terms("actor_attributes.location", "portland"),
	)
	qry.Filter(
		Filter().Terms("repository.has_wiki", true),
	)
	out, err = qry.Result()
	Debug(out)
	Assert(err == nil, t, "should not have error")
	Assert(out.Hits.Len() == 10, t, "Should have 10 docs %v", out.Hits.Len())
	Assert(CloseInt(out.Hits.Total, 42), t, "Should have ~42 total= %v", out.Hits.Total)

	qry = Search("github").Filter(
		"or",
		Filter().Terms("actor_attributes.location", "portland"),
		Filter().Terms("repository.has_wiki", true),
	)
	out, err = qry.Result()
	Assert(err == nil, t, "should not have error")
	Assert(out.Hits.Len() == 10, t, "Should have 10 docs %v", out.Hits.Len())
	Assert(CloseInt(out.Hits.Total, 6344), t, "Should have !6344 total= %v", out.Hits.Total)
}

func TestFilterRange(t *testing.T) {

	// now lets filter range for repositories with more than 100 forks
	out, _ := Search("github").Size("25").Filter(
		Range().Field("repository.forks").From("100"),
	).Result()
	if out == nil || &out.Hits == nil {
		t.Fail()
		return
	}

	Assert(out.Hits.Len() == 25, t, "Should have 25 docs %v", out.Hits.Len())
	Assert(out.Hits.Total == 686, t, "Should have total=686 but was %v", out.Hits.Total)
}

func TestFilterLogicClause(t *testing.T) {
	// test the and,or,not operators including nesting
	qry := Search("github").Filter(
		Filter().Terms("actor_attributes.location", "portland"),
		Filter().Not().Terms("repository.has_wiki", true),
	)
	out, err := qry.Result()
	Assert(err == nil, t, "should not have error")
	Assert(out.Hits.Len() == 10, t, "Should have 10 docs %v", out.Hits.Len())
	Assert(CloseInt(out.Hits.Total, 24), t, "Should have ~24 total= %v", out.Hits.Total)
}

func TestFilterQueryClause(t *testing.T) {
	// test passing a Query Clause into Filters
	qry := Search("github").Filter(
		Query().Search("add"),
	)
	out, err := qry.Result()
	Assert(err == nil, t, "should not have error")
	Assert(out.Hits.Len() == 10, t, "Should have 10 docs %v", out.Hits.Len())
	Assert(CloseInt(out.Hits.Total, 483), t, "Should have ~483 total= %v", out.Hits.Total)

	// combine filter, query clauses
	qry = Search("github").Filter(
		Filter().Terms("actor_attributes.location", "portland"),
		Filter().Not().Terms("repository.has_wiki", true),
		Query().Search("add"),
	)
	out, err = qry.Result()
	Assert(err == nil, t, "should not have error")
	Assert(out.Hits.Len() == 3, t, "Should have 3 docs %v", out.Hits.Len())
	Assert(CloseInt(out.Hits.Total, 3), t, "Should have ~3 total= %v", out.Hits.Total)
}
