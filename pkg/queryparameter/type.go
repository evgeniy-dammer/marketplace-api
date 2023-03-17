package queryparameter

import (
	"time"

	"github.com/evgeniy-dammer/marketplace-api/pkg/pagination"
	"github.com/evgeniy-dammer/marketplace-api/pkg/sort"
)

type QueryParameter struct {
	Search     string
	Sorts      sort.Sorts
	Pagination pagination.Pagination
	StartDate  time.Time
	EndDate    time.Time
}
