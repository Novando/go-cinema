package helper

import (
	"github.com/novando/go-cinema/pkg/logger"
	"reflect"
)

type PaginationData struct {
	Total  uint64   `json:"total"`
	Keys   []string `json:"keys"`
	Values any      `json:"values"`
}

func CreateListResponse(tot uint64, obj interface{}, l ...*logger.Logger) PaginationData {
	emptyList := PaginationData{
		Total:  0,
		Keys:   []string{},
		Values: [][]string{},
	}
	log := logger.Call()
	if len(l) > 0 {
		log = l[0]
	}

	// Get the type of the struct
	ot := reflect.TypeOf(obj)

	// Ensure it's a struct
	if ot.Kind() != reflect.Slice {
		if log != nil {
			log.Warnf("Expect slice, got %v", ot.Kind())
		}
		return emptyList
	}
	et := ot.Elem()
	if et.Kind() != reflect.Struct {
		if log != nil {
			log.Warnf("Expect struct, got %v", et.Kind())
		}
		return emptyList
	}
	// Extract field names as keys
	var keys []string
	for i := 0; i < et.NumField(); i++ {
		field := et.Field(i)

		// Check if a JSON tag exists, otherwise use the field name
		jsonTag := field.Tag.Get("json")
		if jsonTag != "" {
			keys = append(keys, jsonTag)
		} else {
			keys = append(keys, field.Name)
		}
	}

	// Extract values from each struct in the slice
	v := reflect.ValueOf(obj)
	var rv any
	var values [][]interface{}
	for i := 0; i < v.Len(); i++ {
		var row []interface{}
		elem := v.Index(i)
		for j := 0; j < elem.NumField(); j++ {
			row = append(row, elem.Field(j).Interface())
		}
		values = append(values, row)
	}
	if values == nil {
		rv = [][]string{}
	} else {
		rv = values
	}
	return PaginationData{Total: tot, Keys: keys, Values: rv}
}
