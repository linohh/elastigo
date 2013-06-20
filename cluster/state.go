package cluster

import (
	"encoding/json"
	"fmt"
	"github.com/mattbaird/elastigo/api"
	"strings"
)

type ClusterStateFilter struct {
	FilterNodes        bool
	FilterRoutingTable bool
	FilterMetadata     bool
	FilterBlocks       bool
	FilterIndices      []string
}

func (f ClusterStateFilter) Parameterize() []string {
	var parts []string

	if f.FilterNodes {
		parts = append(parts, "filter_nodes=true")
	}

	if f.FilterRoutingTable {
		parts = append(parts, "filter_routing_table=true")
	}

	if f.FilterMetadata {
		parts = append(parts, "filter_metadata=true")
	}

	if f.FilterBlocks {
		parts = append(parts, "filter_blocks=true")
	}

	if f.FilterIndices != nil && len(f.FilterIndices) > 0 {
		parts = append(parts, strings.Join([]string{"filter_indices=", strings.Join(f.FilterIndices, ",")}, ""))
	}

	return parts
}

func ClusterState(pretty bool, filter ClusterStateFilter) (api.ClusterStateResponse, error) {
	var parameters []string
	var url string
	var retval api.ClusterStateResponse

	parameters = filter.Parameterize()

	// prettyfication should be a single parameter somewhere, this is cluttering the method signatures
	if pretty {
		parameters = append(parameters, api.Pretty(pretty))
	}

	url = fmt.Sprintf("/_cluster/state?%s", strings.Join(parameters, "&"))

	body, err := api.DoCommand("GET", url, nil)
	if err != nil {
		return retval, err
	}
	if err == nil {
		// marshall into json
		jsonErr := json.Unmarshal(body, &retval)
		if jsonErr != nil {
			return retval, jsonErr
		}
	}
	return retval, err

}
