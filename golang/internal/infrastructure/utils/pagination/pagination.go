package pagination

import (
	"math"
	"net/http"
	"strconv"
	"strings"

	"gorm.io/gorm"
)

type Pagination struct {
	Limit      int         `json:"limit,omitempty" query:"limit"`
	Page       int         `json:"page,omitempty" query:"page"`
	Sort       string      `json:"sort,omitempty" query:"sort"`
	TotalRows  int64       `json:"total_rows"`
	TotalPages int         `json:"total_pages"`
	Rows       interface{} `json:"rows"`
	Search     string      `json:"search,omitempty" query:"search"`
	ServiceId  *uint64     `json:"service_id,omitempty" query:"service_id"`
	RoleId     *uint64     `json:"role_id,omitempty" query:"role_id"`
	UserId 	   *uint64	   `json:"user_id,omitempty" query:"user_id"`
}

func (p *Pagination) GetOffset() int {
	limit := p.GetLimit()
	if limit == -1 {
		return -1
	}
	return (p.GetPage() - 1) * limit
}

func (p *Pagination) GetLimit() int {
	if p.Limit == 0 {
		p.Limit = 10
	}
	return p.Limit
}

func (p *Pagination) GetPage() int {
	if p.Page == 0 {
		p.Page = 1
	}
	return p.Page
}

func (p *Pagination) GetSort() string {
	if p.Sort == "" {
		p.Sort = "created_at desc"
	}
	return p.Sort
}

func GeneratePaginationFromRequest(r *http.Request) *Pagination {
	limit := 6 // Default limit as requested
	page := 1
	sort := "created_at desc"
	search := ""
	var serviceId *uint64
	var roleId *uint64
	var userId *uint64

	query := r.URL.Query()

	for key, value := range query {
		queryValue := value[len(value)-1]

		switch key {
		case "show": // User asked for "show" as limit
			limit, _ = strconv.Atoi(queryValue)
		case "page":
			page, _ = strconv.Atoi(queryValue)
		case "sort":
			// User asked for "desc" or "asc" based on created_at
			if strings.ToLower(queryValue) == "asc" {
				sort = "created_at asc"
			} else {
				sort = "created_at desc"
			}
		case "search":
			search = queryValue
		case "service_id":
			serviceId = safeParseUint(queryValue)
		case "role_id":
			roleId = safeParseUint(queryValue)
		case "user_id":
			userId = safeParseUint(queryValue)
		}
	}

	return &Pagination{
		Limit:     limit,
		Page:      page,
		Sort:      sort,
		Search:    search,
		ServiceId: serviceId,
		RoleId:    roleId,
		UserId: userId,
	}

}

func (p *Pagination) Paginate(value interface{}, db *gorm.DB) func(db *gorm.DB) *gorm.DB {
	var totalRows int64
	db.Model(value).Count(&totalRows)

	p.TotalRows = totalRows
	limit := p.GetLimit()
	if limit == -1 {
		p.TotalPages = 1
	} else {
		p.TotalPages = int(math.Ceil(float64(totalRows) / float64(limit)))
	}

	return func(db *gorm.DB) *gorm.DB {
		return db.Offset(p.GetOffset()).Limit(limit).Order(p.GetSort())
	}
}

func safeParseUint(s string) *uint64 {
	val, err := strconv.ParseUint(s, 10, 64)
	if err != nil {
		return nil
	}
	return &val
}
