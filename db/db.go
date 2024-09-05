package db

import (
	"context"
	"errors"
	"fmt"
	"log"
	"strconv"
	"strings"
)

func create(db Queryable, querySQL string, args ...any) (int64, error) {
	ctx, cancel := context.WithTimeout(context.Background(), defaultTimeout)
	defer cancel()

	stmt, err := db.PrepareContext(ctx, querySQL)
	if err != nil {
		return 0, err
	}
	defer stmt.Close()

	res, err := stmt.ExecContext(ctx, args...)
	if err != nil {
		return 0, err
	}

	_, err = res.RowsAffected()
	if err != nil {
		return 0, err
	}

	newId, err := res.LastInsertId()
	if err != nil {
		return 0, err
	}

	// 关系表没有id
	if newId <= 0 {
		return 0, errors.New("newId error " + strconv.Itoa(int(newId)))
	}

	return newId, nil
}

func update(db Queryable, querySQL string, args ...any) error {
	ctx, cancel := context.WithTimeout(context.Background(), defaultTimeout)
	defer cancel()
	stmt, err := db.PrepareContext(ctx, querySQL)
	if err != nil {
		return err
	}
	defer stmt.Close()

	res, err := stmt.ExecContext(ctx, args...)
	if err != nil {
		return err
	}

	rows, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if rows <= 0 {
		return errors.New("not found")
	}
	return nil
}

// unqQuerySQL 唯一查询组合SQL
func unqQuerySQL(unqFields []string, values map[string]interface{}) (sqlStr string, params []interface{}) {
	for i, x := range unqFields {
		v, ok := values[x]
		if ok && v != nil {
			if sqlStr == "" {
				sqlStr += fmt.Sprintf("%s = $%d", x, i+1)
			} else {
				sqlStr += fmt.Sprintf(" and %s = $%d", x, i+1)
			}
			params = append(params, v)
		}
	}
	return
}

// multipleInsert 批量插入
// rows: 插入几条数据
func multipleInsert(db Queryable, tableName string, rows int, fileds []string, values []interface{}) error {
	// INSERT INTO table(field1,field2,field3...) VALUES ($1, $2, $3....), ($1, $2, $3....);
	insertSQL := "INSERT INTO " + tableName + "(" + strings.Join(fileds, ",") + ") VALUES "

	for i := 0; i < rows; i++ {
		queryVals := []string{}
		for range fileds {
			// NOTE: 如果采用 $1 语法，这里是不正确的，获取不到正确的下标，这种情况可以直接传一个模板格式
			queryVals = append(queryVals, "?")
		}

		if i < rows-1 {
			insertSQL += "(" + strings.Join(queryVals, ",") + "),"
		} else {
			insertSQL += "(" + strings.Join(queryVals, ",") + ");"
		}
	}

	log.Println(insertSQL)

	ctx, cancel := context.WithTimeout(context.Background(), defaultTimeout)
	defer cancel()
	stmt, err := db.PrepareContext(ctx, insertSQL)
	if err != nil {
		return err
	}
	defer stmt.Close()
	res, err := stmt.ExecContext(ctx, values...)
	if err != nil {
		return err
	}
	affRows, err := res.RowsAffected()
	if err != nil {
		return err
	}
	log.Println("RowsAffected: ", affRows)
	return nil
}
