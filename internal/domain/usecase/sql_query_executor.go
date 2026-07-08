package usecase

import "database/sql"

// SqlQueryExecutor executes raw SQL queries
type SqlQueryExecutor struct {
	db *sql.DB
}

// NewSqlQueryExecutor creates a new SqlQueryExecutor
func NewSqlQueryExecutor(db *sql.DB) *SqlQueryExecutor {
	return &SqlQueryExecutor{db: db}
}

// Execute executes a raw SQL query and returns the result rows as maps
// VULNERABLE: Executes raw SQL without parameterization
func (e *SqlQueryExecutor) Execute(sqlQuery string) ([]map[string]interface{}, error) {
	rows, err := e.db.Query(sqlQuery)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	columns, err := rows.Columns()
	if err != nil {
		return nil, err
	}

	var results []map[string]interface{}
	for rows.Next() {
		values := make([]interface{}, len(columns))
		valuePtrs := make([]interface{}, len(columns))
		for i := range columns {
			valuePtrs[i] = &values[i]
		}

		if err := rows.Scan(valuePtrs...); err != nil {
			return nil, err
		}

		row := make(map[string]interface{})
		for i, col := range columns {
			row[col] = values[i]
		}
		results = append(results, row)
	}

	return results, nil
}
