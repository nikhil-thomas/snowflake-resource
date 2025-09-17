package models

import (
	"database/sql"
	"reflect"
)

// mapRowsToStruct maps SQL rows to a struct using reflection and struct tags
func mapRowsToStruct(rows *sql.Rows, dest any) error {
	// Get the type and value of the destination struct
	destType := reflect.TypeOf(dest)
	destValue := reflect.ValueOf(dest)

	// Get column names from the SQL rows
	columns, err := rows.Columns()
	if err != nil {
		return err
	}

	// Create a map to store the column values
	columnValues := make([]interface{}, len(columns))
	for i := range columnValues {
		// Create a new sql.NullString for each column
		columnValues[i] = new(sql.NullString)
	}

	// Create a map to store the field values
	fieldValues := make([]reflect.Value, len(columns))
	for i := range fieldValues {
		// Get the field name from the struct tag
		fieldName := destType.Field(i).Tag.Get("db")
		// Get the field value from the struct
		fieldValue := destValue.FieldByName(fieldName)
		// Store the field value in the map
		fieldValues[i] = fieldValue
	}

}
