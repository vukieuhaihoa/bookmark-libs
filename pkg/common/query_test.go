package common

import (
	"errors"
	"reflect"
	"testing"
)

func Test_ParseSortParams(t *testing.T) {

	allowed := map[string]struct{}{
		"created_at": {},
		"updated_at": {},
	}

	tests := []struct {
		name string

		sort    string
		allowed map[string]struct{}
		want    []SortedField
		wantErr error
	}{
		{
			name:    "empty sort string",
			sort:    "",
			allowed: allowed,
			want: []SortedField{
				{Field: "created_at", Direction: SortDesc},
			},
			wantErr: nil,
		},
		{
			name:    "single ascending field",
			sort:    "updated_at",
			allowed: allowed,
			want: []SortedField{
				{Field: "updated_at", Direction: SortAsc},
			},
			wantErr: nil,
		},
		{
			name:    "single descending field",
			sort:    "-created_at",
			allowed: allowed,
			want: []SortedField{
				{Field: "created_at", Direction: SortDesc},
			},
			wantErr: nil,
		},
		{
			name:    "multiple fields",
			sort:    "-created_at,updated_at",
			allowed: allowed,
			want: []SortedField{
				{Field: "created_at", Direction: SortDesc},
				{Field: "updated_at", Direction: SortAsc},
			},
			wantErr: nil,
		},
		{
			name:    "invalid field",
			sort:    "invalidField",
			allowed: allowed,
			want:    nil,
			wantErr: ErrInvalidSortField,
		},
		{
			name:    "mixed valid and invalid field",
			sort:    "created_at,invalidField",
			allowed: allowed,
			want:    nil,
			wantErr: ErrInvalidSortField,
		},
		{
			name:    "field with only dash prefix",
			sort:    "-",
			allowed: allowed,
			want:    nil,
			wantErr: ErrInvalidSortField,
		},
		{
			name:    "empty allowed map",
			sort:    "created_at",
			allowed: map[string]struct{}{},
			want:    nil,
			wantErr: ErrInvalidSortField,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParseSortParams(tt.sort, tt.allowed)
			if !errors.Is(err, tt.wantErr) {
				t.Errorf("ParseSortParams() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ParseSortParams() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_Paging_Process(t *testing.T) {
	tests := []struct {
		name string

		inputPage  int
		inputLimit int

		wantPage  int
		wantLimit int
	}{
		{
			name: "page below minimum is set to 1",

			inputPage:  0,
			inputLimit: 10,

			wantPage:  1,
			wantLimit: 10,
		},
		{
			name: "negative page is set to 1",

			inputPage:  -5,
			inputLimit: 10,

			wantPage:  1,
			wantLimit: 10,
		},
		{
			name: "limit below minimum is set to 1",

			inputPage:  1,
			inputLimit: 0,

			wantPage:  1,
			wantLimit: 1,
		},
		{
			name: "limit equal to 1 stays at 1",

			inputPage:  1,
			inputLimit: 1,

			wantPage:  1,
			wantLimit: 1,
		},
		{
			name: "limit above maximum is capped to 50",

			inputPage:  1,
			inputLimit: 100,

			wantPage:  1,
			wantLimit: 50,
		},
		{
			name: "limit equal to 50 is capped to 50",

			inputPage:  1,
			inputLimit: 50,

			wantPage:  1,
			wantLimit: 50,
		},
		{
			name: "valid page and limit are unchanged",

			inputPage:  3,
			inputLimit: 20,

			wantPage:  3,
			wantLimit: 20,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &Paging{Page: tt.inputPage, Limit: tt.inputLimit}
			p.Process()
			if p.Page != tt.wantPage {
				t.Errorf("Process() Page = %v, want %v", p.Page, tt.wantPage)
			}
			if p.Limit != tt.wantLimit {
				t.Errorf("Process() Limit = %v, want %v", p.Limit, tt.wantLimit)
			}
		})
	}
}

func Test_NewQueryOptions(t *testing.T) {
	tests := []struct {
		name string

		inputPage    int
		inputLimit   int
		inputSorting []SortedField

		wantPage    int
		wantLimit   int
		wantSorting []SortedField
	}{
		{
			name: "valid page and limit are unchanged",

			inputPage:    2,
			inputLimit:   10,
			inputSorting: []SortedField{{Field: "created_at", Direction: SortDesc}},

			wantPage:    2,
			wantLimit:   10,
			wantSorting: []SortedField{{Field: "created_at", Direction: SortDesc}},
		},
		{
			name: "page below minimum is normalized to 1",

			inputPage:    0,
			inputLimit:   10,
			inputSorting: nil,

			wantPage:    1,
			wantLimit:   10,
			wantSorting: nil,
		},
		{
			name: "limit above maximum is normalized to 50",

			inputPage:    1,
			inputLimit:   200,
			inputSorting: nil,

			wantPage:    1,
			wantLimit:   50,
			wantSorting: nil,
		},
		{
			name: "limit below minimum is normalized to 1",

			inputPage:    1,
			inputLimit:   0,
			inputSorting: nil,

			wantPage:    1,
			wantLimit:   1,
			wantSorting: nil,
		},
		{
			name: "multiple sorting fields are preserved",

			inputPage:  1,
			inputLimit: 5,
			inputSorting: []SortedField{
				{Field: "created_at", Direction: SortDesc},
				{Field: "updated_at", Direction: SortAsc},
			},

			wantPage:  1,
			wantLimit: 5,
			wantSorting: []SortedField{
				{Field: "created_at", Direction: SortDesc},
				{Field: "updated_at", Direction: SortAsc},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := NewQueryOptions(tt.inputPage, tt.inputLimit, tt.inputSorting)
			if got.Page != tt.wantPage {
				t.Errorf("NewQueryOptions() Page = %v, want %v", got.Page, tt.wantPage)
			}
			if got.Limit != tt.wantLimit {
				t.Errorf("NewQueryOptions() Limit = %v, want %v", got.Limit, tt.wantLimit)
			}
			if !reflect.DeepEqual(got.Sorting, tt.wantSorting) {
				t.Errorf("NewQueryOptions() Sorting = %v, want %v", got.Sorting, tt.wantSorting)
			}
		})
	}
}
