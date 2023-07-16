package filter

import (
	"fmt"
	"gorm.io/gorm"
)

func StringFilter(db *gorm.DB, paramName string, lookupsAndValues map[string][]string) *gorm.DB {
	for lookup, values := range lookupsAndValues {
		switch lookup {
		case Exact.String():
			db = filterExact(db, paramName, values)
		case NotExact.String():
			db = filterNotExact(db, paramName, values)
		case Contains.String():
			db = filterContains(db, paramName, values)
		}
	}
	return db
}

func filter(db *gorm.DB, paramName string, op string, values []string) *gorm.DB {
	for _, value := range values {
		whereStatement := fmt.Sprintf("%s %s ?", paramName, op)
		db = db.Where(whereStatement, value)
	}
	return db
}

func filterExact(db *gorm.DB, paramName string, values []string) *gorm.DB {
	return filter(db, paramName, "=", values)
}

func filterNotExact(db *gorm.DB, paramName string, values []string) *gorm.DB {
	return filter(db, paramName, "!=", values)
}

func filterContains(db *gorm.DB, paramName string, values []string) *gorm.DB {
	for _, value := range values {
		whereStatement := fmt.Sprintf("%s LIKE ?", paramName)
		db = db.Where(whereStatement, "%"+value+"%")
	}
	return db
}
