package cluster

type NodeStatsReponse struct {
	ClusterName string `json:"cluster_name"`
	Nodes       map[string]NodeStatsNodeResponse
}
type NodeStatsNodeResponse struct {
	Name             string                   `json:"name"`
	Timestamp        int64                    `json:"timestamp"`
	TransportAddress string                   `json:"transport_address"`
	Hostname         string                   `json:"hostname"`
	Indices          NodeStatsIndicesResponse `json:"indices"`
}

type NodeStatsIndicesResponse struct {
	Docs     NodeStatsIndicesDocsResponse
	Store    NodeStatsIndicesStoreResponse
	Indexing NodeStatsIndicesIndexingResponse
	Get      NodeStatsIndicesGetResponse
	Search   NodeStatsIndicesStoreResponse
}

type NodeStatsOSResponse struct {
	Timestamp int64                   `json:"timestamp"`
	Uptime    int64                   `json:"uptime_in_millis"`
	LoadAvg   []float64               `json:"load_average"`
	CPU       NodeStatsOSCPUResponse  `json:"cpu"`
	Mem       NodeStatsOSMemResponse  `json:"mem"`
	Swap      NodeStatsOSSwapResponse `json:"swap"`
}

type NodeStatsOSMemResponse struct {
	Free       int64 `json:"free_in_bytes"`
	Used       int64 `json:"used_in_bytes"`
	ActualFree int64 `json:"actual_free_in_bytes"`
	ActualUsed int64 `json:"actual_used_in_bytes"`
}

type NodeStatsOSSwapResponse struct {
	Used int64 `json:"used_in_bytes"`
	Free int64 `json:"free_in_bytes"`
}

type NodeStatsOSCPUResponse struct {
	Sys   int64 `json:"sys"`
	User  int64 `json:"user"`
	Idle  int64 `json:"idle"`
	Steal int64 `json:"stolen"`
}

type NodeStatsIndicesDocsResponse struct {
	Count   int64 `json:"count"`
	Deleted int64 `json:"deleted"`
}

type NodeStatsIndicesStoreResponse struct {
	Size         int64 `json:"size_in_bytes"`
	ThrottleTime int64 `json:"throttle_time_in_millis"`
}

type NodeStatsIndicesIndexingResponse struct {
	IndexTotal    int64 `json:"index_total"`
	IndexTime     int64 `json:"index_time_in_millis"`
	IndexCurrent  int64 `json:"index_current"`
	DeleteTotal   int64 `json:"delete_total"`
	DeleteTime    int64 `json:"delete_time_in_millis"`
	DeleteCurrent int64 `json:"delete_current"`
}

type NodeStatsIndicesGetResponse struct {
	Total        int64 `json:"total"`
	Time         int64 `json:"time_in_millis"`
	ExistsTotal  int64 `json:"exists_total"`
	ExistsTime   int64 `json:"exists_time_in_millis"`
	MissingTotal int64 `json:"missing_total"`
	MissingTime  int64 `json:"missing_time_in_millis"`
	Current      int64 `json:"current"`
}

type NodeStatsIndicesSearchResponse struct {
	OpenContext  int64 `json:"open_contexts"`
	QueryTotal   int64 `json:"query_total"`
	QueryTime    int64 `json:"query_time_in_millis"`
	QueryCurrent int64 `json:"query_current"`
	FetchTotal   int64 `json:"fetch_total"`
	FetchTime    int64 `json:"fetch_time_in_millis"`
	FetchCurrent int64 `json:"fetch_current"`
}
