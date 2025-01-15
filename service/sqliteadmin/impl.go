package sqliteadmin

import (
	"time"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type sqliteAdmin struct {
	db *gorm.DB
}

func NewSQLiteAdmin(db *gorm.DB) *sqliteAdmin {
	return &sqliteAdmin{
		db: db,
	}
}

func (s *sqliteAdmin) Register(r fiber.Router) {
	r.Get("/tables", s.handleTables)
	r.Get("/data", s.handleGetData)
	r.Post("/run", s.handleRun)
}

func (s *sqliteAdmin) handleRun(c *fiber.Ctx) error {
	sql := string(c.Body())
	useQuery := c.QueryBool("query") || useQuery(sql) // 是否使用强制查询模式
	data := make([]map[string]any, 0)

	var result *gorm.DB
	start := time.Now()
	if useQuery {
		result = s.db.Raw(sql).Scan(&data)
	} else {
		result = s.db.Exec(sql)
	}
	timeUsed := time.Since(start)

	retMap := fiber.Map{"time": timeUsed.String()}
	if result.Error != nil {
		retMap["error"] = result.Error.Error()
		return c.JSON(retMap)
	}
	if useQuery {
		retMap["data"] = data
	} else {
		retMap["rows_affected"] = result.RowsAffected
	}
	return c.JSON(retMap)
}

func (s *sqliteAdmin) handleTables(c *fiber.Ctx) error {
	start := time.Now()
	var tableNames []string
	err := s.db.Raw("SELECT name FROM sqlite_master WHERE type='table' AND name NOT LIKE 'sqlite_%' ORDER BY name;").Scan(&tableNames).Error
	if err != nil {
		return err
	}
	ret := fiber.Map{"tables": tableNames}
	if c.QueryBool("count") {
		countMap := make(map[string]int64)
		for _, tableName := range tableNames {
			var count int64
			err = s.db.Raw("SELECT COUNT(*) FROM " + tableName).Scan(&count).Error
			if err != nil {
				return err
			}
			countMap[tableName] = count
		}
		ret["tables"] = countMap
	}
	ret["time"] = time.Since(start).String()
	return c.JSON(ret)
}

func (s *sqliteAdmin) handleGetData(c *fiber.Ctx) error {
	page := c.QueryInt("page")
	if page <= 0 {
		page = 1
	}
	pageSize := c.QueryInt("size")
	if pageSize <= 0 {
		pageSize = 20
	}
	tableName := c.Query("table")
	if tableName == "" {
		return fiber.ErrBadRequest
	}
	offset := (page - 1) * pageSize
	orderBy := c.Query("order")

	start := time.Now()
	var count int64
	err := s.db.Table(tableName).Count(&count).Error
	timeUsed := time.Since(start)
	ret := fiber.Map{"time": timeUsed.String()}
	if err != nil {
		ret["error"] = err.Error()
		return c.JSON(ret)
	}
	ret["page"] = fiber.Map{
		"page": page,
		"size": pageSize,
		"rows": count,
	}

	data := make([]map[string]any, 0)
	query := s.db.Table(tableName)
	if orderBy != "" {
		query = query.Order(orderBy)
	}
	err = query.Limit(pageSize).Offset(offset).Find(&data).Error
	timeUsed = time.Since(start)
	ret["time"] = timeUsed.String()
	if err != nil {
		ret["error"] = err.Error()
		return c.JSON(ret)
	}

	ret["data"] = data
	return c.JSON(ret)
}
