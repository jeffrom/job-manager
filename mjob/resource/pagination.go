package resource

type Pagination struct {
	Limit  int64  `json:"limit,omitempty"`
	LastID string `json:"last_id,omitempty"`
}
