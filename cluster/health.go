package cluster

import (
	"encoding/json"
	"fmt"
	"github.com/mattbaird/elastigo/api"
	"strings"
)

// The cluster health API allows to get a very simple status on the health of the cluster.
// see http://www.elasticsearch.org/guide/reference/api/admin-cluster-health.html
// TODO: implement wait_for_status, timeout, wait_for_relocating_shards, wait_for_nodes
// TODO: implement level (Can be one of cluster, indices or shards. Controls the details level of the health
// information returned. Defaults to cluster.)
func Health(indices ...string) (api.ClusterHealthResponse, error) {
	var url string
	var retval api.ClusterHealthResponse
	if len(indices) > 0 {
		url = fmt.Sprintf("/_cluster/health/%s", strings.Join(indices, ","))
	} else {
		url = fmt.Sprintf("/_cluster/health")
	}
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
	//fmt.Println(body)
	return retval, err
}
