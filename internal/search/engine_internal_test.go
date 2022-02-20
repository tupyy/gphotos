package search

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/tupyy/gophoto/internal/entity"
)

func FilterEngine(t *testing.T) {
	data := []struct {
		expr     string
		album    entity.Album
		expected bool
		hasError bool
	}{
		{
			expr:     "name = 'test'",
			album:    entity.Album{Name: "test"},
			expected: true,
		},
		{
			expr:     "name != 'test'",
			album:    entity.Album{Name: "test"},
			expected: false,
		},
		{
			expr:     "name = 'test' & description != 'toto' & location = 'loc'",
			album:    entity.Album{Name: "test", Description: "titi", Location: "loc"},
			expected: true,
		},
		{
			expr:     "name != 'test' & description != 'toto'",
			album:    entity.Album{Name: "t", Description: "titi", Location: "loc"},
			expected: true,
		},
		{
			expr:     "name = 'test' & description != 'toto' | location = 'loc'",
			album:    entity.Album{Name: "test", Description: "titi", Location: "loc"},
			expected: true,
		},
		{
			expr:     "name = 'test' & tag = 'tag'",
			album:    entity.Album{Name: "test", Tags: []entity.Tag{{Name: "tag"}}},
			expected: true,
		},
		{
			expr:     "name = 'test' & tag = 'tag2'",
			album:    entity.Album{Name: "test", Tags: []entity.Tag{{Name: "tag"}}},
			expected: false,
		},
		{
			expr:     "name = 'test' & tag > 'tag2'",
			album:    entity.Album{Name: "test", Tags: []entity.Tag{{Name: "tag"}}},
			expected: false,
			hasError: true,
		},
		{
			expr:     "date > '01/01/2022' & date < '01/03/2022'",
			album:    entity.Album{CreatedAt: createDate(2022, 04, 01)},
			expected: false,
		},
		{
			expr:     "date > '11/01/2021' & date < '16/11/2021'",
			album:    entity.Album{CreatedAt: createDate(2021, 11, 15)},
			expected: true,
		},
		{
			expr:     "date > '01/01/2022' & date < '01/03/2022'",
			album:    entity.Album{CreatedAt: createDate(2022, 02, 01)},
			expected: true,
		},
		{
			expr:     "date < '01/03/2022'",
			album:    entity.Album{CreatedAt: createDate(2022, 02, 01)},
			expected: true,
		},
		{
			expr:     "date = '01/02/2022'",
			album:    entity.Album{CreatedAt: createDate(2022, 02, 01)},
			expected: true,
		},
		{
			expr:     "date != '01/02/2022'",
			album:    entity.Album{CreatedAt: createDate(2022, 02, 02)},
			expected: true,
		},
		{
			expr:     "name = 'titi' & date != '01/02/2022'",
			album:    entity.Album{Name: "titi", CreatedAt: createDate(2022, 02, 02)},
			expected: true,
		},
		{
			expr:     "blabla = 'titi' & date != '01/02/2022'",
			album:    entity.Album{Name: "titi", CreatedAt: createDate(2022, 02, 02)},
			expected: false,
			hasError: true,
		},
	}

	for _, d := range data {
		t.Run(d.expr, func(t *testing.T) {
			filterEngine, err := NewSearchEngine(d.expr)
			assert.Nil(t, err)

			result, err := filterEngine.Resolve(d.album)
			assert.Equal(t, d.expected, result)
			if d.hasError {
				assert.NotNil(t, err)
			} else {
				assert.Nil(t, err)
			}

		})
	}
}

func createDate(year, month, day int) time.Time {
	return time.Date(year, time.Month(month), day, 0, 0, 0, 0, time.UTC)
}
