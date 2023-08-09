package oplog

import (
	"fmt"
	"strings"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

func generateInsertStatement(tableName string, data map[string]interface{}) string {
	var columns []string
	var values []string
	for key, value := range data {
		columns = append(columns, key)
		values = append(values, fmt.Sprintf("%v", value))
	}
	columnsStr := strings.Join(columns, ", ")
	valuesStr := strings.Join(values, ", ")

	insertStatement := fmt.Sprintf("INSERT INTO %s (%s) VALUES (%s);", tableName, columnsStr, valuesStr)
	return insertStatement
}

func ParseDocument(tablename string, obj map[string]interface{}) []string {
	var queries []string
	data := make(map[string]interface{})
	for key, value := range obj {
		switch v := value.(type) {
		case primitive.ObjectID:
			data[key] = v.Hex()
		case map[string]interface{}:
			query := generateInsertStatement(key, v)
			queries = append(queries, query)

		case primitive.A:
			for _, item := range v {
				query := generateInsertStatement(key, item.(map[string]interface{}))
				queries = append(queries, query)
			}

		default:
			data[key] = value
		}
	}
	query := generateInsertStatement(tablename, data)
	queries = append(queries, query)
	return queries
}
