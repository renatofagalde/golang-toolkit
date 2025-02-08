package config

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"strconv"
)

func Paginate(ctx *gin.Context) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		page := getPage(ctx)
		pageSize := getPageSize(ctx)
		offset := (page - 1) * pageSize

		return db.Offset(offset).Limit(pageSize)
	}
}

func getPage(ctx *gin.Context) int {
	p, err := strconv.Atoi(ctx.DefaultQuery("page", "1"))
	if err != nil || p <= 0 {
		return 1
	}
	return p
}

func getPageSize(ctx *gin.Context) int {
	ps, err := strconv.Atoi(ctx.DefaultQuery("page_size", "100"))
	if err != nil {
		ps = 100
	}
	switch {
	case ps > 100:
		ps = 100
	case ps <= 0:
		ps = 10
	}
	return ps
}
