package album

import (
	"time"

	"github.com/tupyy/gophoto/internal/domain/filters/album"
)

type Predicate func() album.Filter

func BeforeDate(date time.Time) Predicate {
	return func() album.Filter {
		f, _ := album.GenerateFilterFuncs(album.FilterBeforeDate, date)
		return f
	}
}

func AfterDate(date time.Time) Predicate {
	return func() album.Filter {
		f, _ := album.GenerateFilterFuncs(album.FilterAfterDate, date)
		return f
	}
}

func Owner(id []string) Predicate {
	return func() album.Filter {
		f, _ := album.GenerateFilterFuncs(album.FilterByOwnerID, id)
		return f
	}
}
