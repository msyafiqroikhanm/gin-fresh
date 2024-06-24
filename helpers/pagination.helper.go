package helpers

import (
	"fmt"
	"math"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type Pagination struct {
	Page         int         `json:"page"`
	Limit        int         `json:"limit"`
	TotalPages   int         `json:"total_pages"`
	TotalRows    int         `json:"total_rows"`
	FirstPage    string      `json:"first_page"`
	PreviousPage string      `json:"previous_page"`
	NextPage     string      `json:"next_page"`
	LastPage     string      `json:"last_page"`
	FromRow      int         `json:"from_row"`
	ToRow        int         `json:"to_row"`
	Rows         interface{} `json:"rows"`
}

func Paginate(c *gin.Context) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
		if page <= 0 {
			page = 1
		}

		limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
		switch {
		case limit <= 0:
			limit = 10
		}

		offset := (page - 1) * limit
		return db.Offset(offset).Limit(limit)
	}
}

func GeneratePaginatedQuery(c *gin.Context, totalRows int64, data []interface{}) *Pagination {
	// initilize required variable
	var nextPage, previousPage string
	var fromRow, toRow int
	totalRow := int(totalRows)

	// getting and setting page
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	if page <= 0 {
		page = 1
	}

	// getting and setting page
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	switch {
	case limit <= 0:
		limit = 10
	}

	// Calculate total page using totalRow [len(data)] and limit
	totalPages := int(math.Ceil(float64(totalRow) / float64(limit)))

	// Set url for first and last page
	firstPage := fmt.Sprintf("%s?page=1&limit=%d", c.Request.URL.Path, limit)
	lastPage := fmt.Sprintf("%s?page=%d&limit=%d", c.Request.URL.Path, totalPages, limit)

	// Set url for previous and next page
	if page > 1 {
		previousPage = fmt.Sprintf("%s?page=%d&limit=%d", c.Request.URL.Path, page-1, limit)
	}
	if page < totalRow {
		nextPage = fmt.Sprintf("%s?page=%d&limit=%d", c.Request.URL.Path, page+1, limit)
	}

	// Set from and to row (index)
	if page == 1 {
		fromRow = 1
		if limit > totalRow {
			toRow = totalRow
		} else {
			toRow = limit
		}
	} else {
		if page <= totalPages {
			fromRow = (page * limit) + 1
			toRow = (page + 1) * limit
		}
	}

	return &Pagination{
		Page:         page,
		Limit:        limit,
		TotalPages:   totalPages,
		TotalRows:    totalRow,
		FirstPage:    firstPage,
		PreviousPage: previousPage,
		NextPage:     nextPage,
		LastPage:     lastPage,
		FromRow:      fromRow,
		ToRow:        toRow,
		Rows:         data,
	}
}
