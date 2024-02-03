package main

import (
	"database/sql"
	"fmt"
	"reflect"
	"strings"
)

type ORM struct {
	db *sql.DB
}

func NewORM(db *sql.DB) *ORM {
	return &ORM{db}
}

func (o *ORM) Insert(table string, data interface{}) error {
	// 获取数据的类型和值
	valueType := reflect.TypeOf(data)
	value := reflect.ValueOf(data)

	// 构建INSERT语句
	columns := make([]string, 0)
	placeholders := make([]string, 0)
	values := make([]interface{}, 0)

	// 原理就是遍历data的所有属性和字段名 拼接成sql
	for i := 0; i < valueType.NumField(); i++ {
		field := valueType.Field(i)
		column := field.Tag.Get("db")
		if column == "" {
			column = strings.ToLower(field.Name)
		}
		if column == "id" && field.Tag.Get("auto") == "true" {
			continue
		}

		columns = append(columns, column)
		// 使用占位符 而不是拼接防止sql注入
		placeholders = append(placeholders, "?")
		values = append(values, value.Field(i).Interface())
	}
	query := fmt.Sprintf("INSERT INTO %s (%s) VALUES (%s)", table, strings.Join(columns, ", "), strings.Join(placeholders, ", "))
	fmt.Println(query, values)
	_, err := o.db.Exec(query, values...)
	if err != nil {
		return err
	}

	return nil
}

func (o *ORM) Update(table string, data interface{}, where string, args ...interface{}) error {
	valueType := reflect.TypeOf(data)
	value := reflect.ValueOf(data)

	updates := make([]string, 0)
	values := make([]interface{}, 0)

	for i := 0; i < valueType.NumField(); i++ {
		field := valueType.Field(i)
		column := field.Tag.Get("db")
		if column == "" {
			column = strings.ToLower(field.Name)
		}
		if column == "id" && field.Tag.Get("auto") == "true" {
			continue
		}

		updates = append(updates, fmt.Sprintf("%s = ?", column))
		values = append(values, value.Field(i).Interface())
	}

	query := fmt.Sprintf("UPDATE %s SET %s WHERE %s", table, strings.Join(updates, ", "), where)
	_, err := o.db.Exec(query, append(values, args...)...)
	if err != nil {
		return err
	}

	return nil
}

func (o *ORM) FindAll(table string, where string, args ...interface{}) ([]map[string]interface{}, error) {
	query := fmt.Sprintf("SELECT * FROM %s WHERE %s", table, where)
	rows, err := o.db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	columns, err := rows.Columns()
	if err != nil {
		return nil, err
	}

	var result []map[string]interface{}

	for rows.Next() {
		values := make([]interface{}, len(columns))
		valuePointers := make([]interface{}, len(columns))
		for i := range columns {
			valuePointers[i] = &values[i]
		}

		if err := rows.Scan(valuePointers...); err != nil {
			return nil, err
		}

		rowData := make(map[string]interface{})
		for i, colName := range columns {
			val := values[i]
			if val == nil {
				rowData[colName] = nil
			} else {
				rowData[colName] = val
			}
		}

		result = append(result, rowData)
	}

	return result, nil
}

func (o *ORM) Delete(table string, where string, args ...interface{}) error {
	query := fmt.Sprintf("DELETE FROM %s WHERE %s", table, where)
	_, err := o.db.Exec(query, args...)
	if err != nil {
		return err
	}

	return nil
}
