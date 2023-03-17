package queryparameter

import (
	"github.com/evgeniy-dammer/marketplace-api/pkg/pagination"
	"github.com/evgeniy-dammer/marketplace-api/pkg/sort"
	"time"
)

type QueryParameter struct {
	Search     string
	Sorts      sort.Sorts
	Pagination pagination.Pagination
	StartDate  time.Time
	EndDate    time.Time
}
