package pagination

import (
	"strings"

	"github.com/pkg/errors"
	"github.com/tguankheng016/go-ecommerce-microservice/internal/pkg/core/helpers"
)

type PageRequest struct {
	MaxResultCount int    `query:"maxResultCount" json:"maxResultCount,omitempty"`
	SkipCount      int    `query:"skipCount" json:"skipCount,omitempty"`
	Sorting        string `query:"sorting" json:"sorting,omitempty"`
	Filters        string `query:"filters" json:"filters,omitempty"`
}

func (p PageRequest) SanitizeSorting(validSortFields ...string) error {
	if p.Sorting == "" {
		return nil
	}

	sortFieldsMap := make(map[string]struct{}) // Initialize the map

	// Populate the map with the strings
	for _, str := range validSortFields {
		sortFieldsMap[strings.ToLower(str)] = struct{}{} // Use struct{} to save memory
	}

	// Remove empty space, asc, desc keywords
	sorting := strings.ToLower(p.Sorting)
	sorting = helpers.ReplaceLast(sorting, " asc", "")
	sorting = helpers.ReplaceLast(sorting, " desc", "")
	sorting = strings.Replace(sorting, " ", "", -1)

	// Check if the requested field is valid
	if _, ok := sortFieldsMap[sorting]; !ok {
		return errors.New("invalid sorting field")
	}

	return nil
}
