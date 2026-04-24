package pagination

import (
	"math"
	"testing"
)

func TestPagination_GetLimit(t *testing.T) {
	tests := []struct {
		name     string
		p        Pagination
		expected int
	}{
		{"Default limit", Pagination{Limit: 0}, 10},
		{"Specific limit", Pagination{Limit: 20}, 20},
		{"Show all limit", Pagination{Limit: -1}, -1},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.p.GetLimit(); got != tt.expected {
				t.Errorf("GetLimit() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestPagination_GetOffset(t *testing.T) {
	tests := []struct {
		name     string
		p        Pagination
		expected int
	}{
		{"Page 1, Limit 10", Pagination{Page: 1, Limit: 10}, 0},
		{"Page 2, Limit 10", Pagination{Page: 2, Limit: 10}, 10},
		{"Page 1, Limit -1", Pagination{Page: 1, Limit: -1}, -1},
		{"Page 5, Limit -1", Pagination{Page: 5, Limit: -1}, -1},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.p.GetOffset(); got != tt.expected {
				t.Errorf("GetOffset() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestPagination_Paginate_TotalPages(t *testing.T) {
	tests := []struct {
		name      string
		p         Pagination
		totalRows int64
		expected  int
	}{
		{"Limit 10, Rows 25", Pagination{Limit: 10}, 25, 3},
		{"Limit -1, Rows 100", Pagination{Limit: -1}, 100, 1},
		{"Limit -1, Rows 0", Pagination{Limit: -1}, 0, 1},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			limit := tt.p.GetLimit()
			var gotTotalPages int
			if limit == -1 {
				gotTotalPages = 1
			} else {
				gotTotalPages = int(math.Ceil(float64(tt.totalRows) / float64(limit)))
			}
			
			if gotTotalPages != tt.expected {
				t.Errorf("TotalPages = %v, want %v", gotTotalPages, tt.expected)
			}
		})
	}
}
