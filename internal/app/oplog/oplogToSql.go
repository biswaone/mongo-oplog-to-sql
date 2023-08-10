package oplog

import (
	"fmt"
	"strings"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

func generateSqlEmbeddedDocument(tableName string, data map[string]interface{}, sqlOp string, objectId string) string {
	var columns []string
	var values []string
	columns = append(columns, "_id")
	values = append(values, objectId)
	for key, value := range data {
		columns = append(columns, key)
		values = append(values, fmt.Sprintf("%v", value))
	}
	columnsStr := strings.Join(columns, ", ")
	valuesStr := strings.Join(values, ", ")

	insertStatement := fmt.Sprintf("%s INTO %s (%s) VALUES (%s);", sqlOp, tableName, columnsStr, valuesStr)
	return insertStatement
}

func GenerateSqlDocument(tablename string, obj map[string]interface{}, operation string) []string {
	var queries []string
	data := make(map[string]interface{})
	objectId := obj["_id"].(primitive.ObjectID).Hex()
	for key, value := range obj {
		switch v := value.(type) {
		case primitive.ObjectID:
			data[key] = v.Hex()
		case map[string]interface{}:
			query := generateSqlEmbeddedDocument(key, v, operation, objectId)
			queries = append(queries, query)
		case primitive.A:
			for _, item := range v {
				query := generateSqlEmbeddedDocument(key, item.(map[string]interface{}), operation, objectId)
				queries = append(queries, query)
			}

		default:
			data[key] = value
		}
	}
	query := generateSqlEmbeddedDocument(tablename, data, operation, objectId)
	queries = append(queries, query)
	return queries
}
