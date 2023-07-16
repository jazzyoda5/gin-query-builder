package filter

import (
	"gorm.io/gorm"
	"net/url"
	"strings"
)

type ColType byte

const (
	Text = iota
)

type Filter struct {
	Name      string
	FieldName string
	ColType   ColType
}

func BuildQuery(db *gorm.DB, params url.Values, fs []Filter) *gorm.DB {
	var filterMap = map[string]Filter{}
	for _, f := range fs {
		filterMap[f.Name] = f
	}

	processedParams := extractLookups(params)
	for paramName, lookupsAndValues := range processedParams {
		f := filterMap[paramName]
		switch f.ColType {
		case Text:
			db = StringFilter(db, paramName, lookupsAndValues)
		}
	}
	return db
}

type ProcessedParams map[string]map[string][]string

type LookupType byte

const (
	Exact LookupType = iota
	NotExact
	Contains
	IContains
)

func (lookup LookupType) String() string {
	switch lookup {
	case Exact:
		return "ex"
	case NotExact:
		return "nex"
	case Contains:
		return "con"
	case IContains:
		return "icon"
	default:
		return ""
	}
}

func extractLookups(params url.Values) ProcessedParams {
	processedParams := make(ProcessedParams)
	for paramName, values := range params {
		for _, fullValue := range values {
			splitValue := strings.SplitN(fullValue, ":", 2)

			lookup := "ex"
			value := splitValue[0]
			if len(splitValue) > 1 {
				lookup = splitValue[0]
				value = splitValue[1]
			}

			if processedParams[paramName] == nil {
				processedParams[paramName] = make(map[string][]string)
			}
			if processedParams[paramName][lookup] == nil {
				processedParams[paramName][lookup] = []string{}
			}

			processedParams[paramName][lookup] = append(processedParams[paramName][lookup], value)
		}
	}
	return processedParams
}
