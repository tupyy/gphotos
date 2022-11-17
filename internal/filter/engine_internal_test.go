package filter

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/tupyy/gophoto/internal/entity"
)

func TestFilterEngine(t *testing.T) {
	data := []struct {
		expr     string
		album    entity.Album
		expected bool
		hasError bool
	}{
		// {
		// 	expr:     "name = 'test'",
		// 	album:    entity.Album{Name: "test"},
		// 	expected: true,
		// },
		// {
		// 	expr:     "name != 'test'",
		// 	album:    entity.Album{Name: "test"},
		// 	expected: false,
		// },
		// {
		// 	expr:     "name = 'test' and description != 'toto' and location = 'loc'",
		// 	album:    entity.Album{Name: "test", Description: "titi", Location: "loc"},
		// 	expected: true,
		// },
		// {
		// 	expr:     "name != 'test' and description != 'toto'",
		// 	album:    entity.Album{Name: "t", Description: "titi", Location: "loc"},
		// 	expected: true,
		// },
		// {
		// 	expr:     "name = 'test' and description != 'toto' or location = 'loc'",
		// 	album:    entity.Album{Name: "test", Description: "titi", Location: "loc"},
		// 	expected: true,
		// },
		// {
		// 	expr:     "name = 'test' and tag = 'tag'",
		// 	album:    entity.Album{Name: "test", Tags: []entity.Tag{{Name: "tag"}}},
		// 	expected: true,
		// },
		// {
		// 	expr:     "name = 'test' and tag = 'tag2'",
		// 	album:    entity.Album{Name: "test", Tags: []entity.Tag{{Name: "tag"}}},
		// 	expected: false,
		// },
		// {
		// 	expr:     "name = 'test' and tag > 'tag2'",
		// 	album:    entity.Album{Name: "test", Tags: []entity.Tag{{Name: "tag"}}},
		// 	expected: false,
		// 	hasError: true,
		// },
		{
			expr:     "date > '01/01/2022' and date < '01/03/2022'",
			album:    entity.Album{CreatedAt: createDate(2022, 04, 01)},
			expected: false,
		},
		{
			expr:     "date > '11/01/2021' and date < '16/11/2021'",
			album:    entity.Album{CreatedAt: createDate(2021, 11, 15)},
			expected: true,
		},
		{
			expr:     "date > '01/01/2022' and date < '01/03/2022'",
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
			expr:     "name = 'titi' and date != '01/02/2022'",
			album:    entity.Album{Name: "titi", CreatedAt: createDate(2022, 02, 02)},
			expected: true,
		},
		{
			expr:     "blabla = 'titi' and date != '01/02/2022'",
			album:    entity.Album{Name: "titi", CreatedAt: createDate(2022, 02, 02)},
			expected: false,
			hasError: true,
		},
		{
			expr:     "location like 't.t.'",
			album:    entity.Album{Name: "toto", Location: "tata", CreatedAt: createDate(2022, 02, 02)},
			expected: true,
		},
		{
			expr:     "location like 't[a-z]{1}t.'",
			album:    entity.Album{Name: "toto", Location: "tata", CreatedAt: createDate(2022, 02, 02)},
			expected: true,
		},
		{
			expr:     "location in ['loc1','loc2']",
			album:    entity.Album{Name: "toto", Location: "loc1", CreatedAt: createDate(2022, 02, 02)},
			expected: true,
		},
		{
			expr:     "location in ['loc1','loc2']",
			album:    entity.Album{Name: "toto", Location: "loc3", CreatedAt: createDate(2022, 02, 02)},
			expected: false,
		},
		{
			expr: "permissions.user = 'toto'",
			album: entity.Album{Name: "toto", Location: "loc3", CreatedAt: createDate(2022, 02, 02), UserPermissions: []entity.AlbumPermission{
				{
					OwnerID: "toto",
				},
			}},
			expected: true,
		},
		{
			expr: "permissions.user = 'toto2'",
			album: entity.Album{Name: "toto", Location: "loc3", CreatedAt: createDate(2022, 02, 02), UserPermissions: []entity.AlbumPermission{
				{
					OwnerID: "toto",
				},
			}},
			expected: false,
		},
		{
			expr: "permissions.group = 'toto'",
			album: entity.Album{Name: "toto", Location: "loc3", CreatedAt: createDate(2022, 02, 02), GroupPermissions: []entity.AlbumPermission{
				{
					OwnerID: "toto",
				},
			}},
			expected: true,
		},
		{
			expr: "permissions.group = 'toto2'",
			album: entity.Album{Name: "toto", Location: "loc3", CreatedAt: createDate(2022, 02, 02), GroupPermissions: []entity.AlbumPermission{
				{
					OwnerID: "toto",
				},
			}},
			expected: false,
		},
	}

	for _, d := range data {
		t.Run(d.expr, func(t *testing.T) {
			filterEngine, err := New(d.expr)
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
