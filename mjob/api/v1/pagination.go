package v1

import "github.com/jeffrom/job-manager/mjob/resource"

func PaginationToResource(page *Pagination) *resource.Pagination {
	if page == nil {
		return nil
	}

	return &resource.Pagination{
		Limit:  page.Limit,
		LastID: page.LastId,
	}
}

func PaginationToProto(page *resource.Pagination) *Pagination {
	if page == nil {
		return nil
	}

	return &Pagination{
		Limit:  page.Limit,
		LastId: page.LastID,
	}
}
