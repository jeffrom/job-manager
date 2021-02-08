package resource

type Stats struct {
	Queued           int64 `json:"queued"`
	Running          int64 `json:"running"`
	Complete         int64 `json:"complete"`
	Dead             int64 `json:"dead"`
	Cancelled        int64 `json:"cancelled"`
	Invalid          int64 `json:"invalid"`
	Failed           int64 `json:"failed"`
	LongestUnstarted int64 `json:"longest_unstarted_secs"`
}
