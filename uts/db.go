package uts

import (
	"database/sql"
	"os"
	"os/exec"
	"strconv"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

func BackupMysql() {
	_args := []string{}
	_args = append(_args, "-u"+Conf.String("sql_account"))
	_args = append(_args, "-p"+Conf.String("sql_pwd"))
	_args = append(_args, Conf.String("sql_db_name"))

	if _, err := os.Stat("./mysql/backup"); err != nil && os.IsNotExist(err) {
		err := os.MkdirAll("./mysql/backup", os.ModePerm)
		ChkErr(err)
	}
	if _, err := os.Stat("./mysql/backup"); err != nil && os.IsNotExist(err) {
		Log("./mysql/backup 创建失败")
	} else {
		out, err := os.OpenFile("./mysql/backup/dump"+time.Now().Format("2006-01-02-15-04-05")+".sql", os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0644)
		if ChkErr(err) {
			Log("mysql backup fail!!")
		} else {
			cmd := exec.Command("mysqldump", _args...)
			cmd.Stdout = out
			_err := cmd.Run()
			if ChkErr(_err) {
				Log("mysql backup fail!!")
			} else {
				Log("mysql backup success!!")
			}
		}
		out.Close()
	}
}

func AutoBackupMysql() {
	_args := []string{}
	_args = append(_args, "-u"+Conf.String("sql_account"))
	_args = append(_args, "-p"+Conf.String("sql_pwd"))
	_args = append(_args, Conf.String("sql_db_name"))

	if _, err := os.Stat("./mysql"); err != nil && os.IsNotExist(err) {
		os.Mkdir("./mysql", os.ModePerm)
	} else if _, err := os.Stat("./mysql/dump.sql"); err == nil || os.IsExist(err) {
		os.RemoveAll("./mysql/dump.sql")
	}

	out, err := os.OpenFile("./mysql/dump.sql", os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		LogF("Cache2DB1", err.Error())
	}
	defer out.Close()
	cmd := exec.Command("mysqldump", _args...)
	cmd.Stdout = out
	_err := cmd.Run()
	if _err != nil {
		LogF("Cache2DB1", _err.Error())
	}
	Log("mysql auto backup success!!")
}

func ExecMySQLCommand(sqlStr string) sql.Result {
	mDB := GetDBConnect(Conf.String("sql_account"), Conf.String("sql_pwd"), Conf.String("sql_host"), Conf.String("sql_port"), Conf.String("sql_db_name"))
	if mDB != nil {
		defer mDB.Close()
		result, err := mDB.Exec(sqlStr)
		if !ChkErr(err) {
			Log(result.RowsAffected())
			return result
		}
	}
	return nil
}

func ExecMySQLCommandInsert(sqlStr string, values ...any) bool {
	mDB := GetDBConnect(Conf.String("sql_account"), Conf.String("sql_pwd"), Conf.String("sql_host"), Conf.String("sql_port"), Conf.String("sql_db_name"))
	if mDB != nil {
		defer mDB.Close()
		result, err := mDB.Prepare(sqlStr)
		if !ChkErr(err) {
			_, err := result.Exec(values...)
			if !ChkErr(err) {
				return true
			}
		}
	}
	return false
}

func GetDBConnect(user string, pwd string, host string, port string, dbName string) *sql.DB {
	mDB, err := sql.Open("mysql", user+":"+pwd+"@tcp("+host+":"+port+")/"+dbName+"?multiStatements=true")
	if !ChkErr(err) {
		err = mDB.Ping()
	}
	if ChkErr(err) {
		mDB.Close()
		mDB = nil
	}
	return mDB
}
func FindAllTable() map[string]string {
	db := GetDBConnect(Conf.String("sql_account"), Conf.String("sql_pwd"), Conf.String("sql_host"), Conf.String("sql_port"), Conf.String("sql_db_name"))

	var result = make(map[string]string)
	if db == nil {
		return result
	}
	defer db.Close()

	sqlStr := `SELECT table_name tableName,TABLE_COMMENT tableDesc
				FROM INFORMATION_SCHEMA.TABLES 
				WHERE UPPER(table_type)='BASE TABLE'
				AND LOWER(table_schema) = ? 
				ORDER BY table_name asc`

	rows, _ := db.Query(sqlStr, Conf.String("sql_db_name"))

	for rows.Next() {
		var tableName, tableDesc string
		_ = rows.Scan(&tableName, &tableDesc)

		// if len(tableDesc) == 0 {
		// 	tableDesc = tableName
		// }
		result[tableName] = tableDesc
	}
	return result
}

func FindDBTableField(tableName string) []Field {
	sql_str := `select 	column_name fName,
						column_type fType,
						data_type dType,
						column_comment fDesc,
						column_default fDefault,
						is_nullable isNull
			from information_schema.COLUMNS 
			where table_schema = ? and table_name = ?
			order by ordinal_position;`

	var result []Field
	db := GetDBConnect(Conf.String("sql_account"), Conf.String("sql_pwd"), Conf.String("sql_host"), Conf.String("sql_port"), Conf.String("sql_db_name"))
	if db == nil {
		return result
	}
	defer db.Close()

	rows, _ := db.Query(sql_str, Conf.String("sql_db_name"), tableName)

	for rows.Next() {
		var f Field
		_ = rows.Scan(&f.FieldName, &f.FieldType, &f.DataType, &f.FieldDesc, &f.FieldDefault, &f.IsNull)
		result = append(result, f)
	}
	return result
}

type Field struct {
	FieldName    string
	FieldType    string
	DataType     string
	FieldDesc    string
	FieldDefault string
	IsNull       string
}

func FindDBTableData(tableName string) []map[string]string {
	sql_str := "select * from " + tableName + ";"

	var result []map[string]string
	db := GetDBConnect(Conf.String("sql_account"), Conf.String("sql_pwd"), Conf.String("sql_host"), Conf.String("sql_port"), Conf.String("sql_db_name"))
	if db == nil {
		return result
	}
	defer db.Close()

	rows, err := db.Query(sql_str)
	if !ChkErrNormal(err) {
		result, err = getTableDataMap(rows)
		ChkErrNormal(err)
	}
	return result
}

func getTableDataMap(query *sql.Rows) ([]map[string]string, error) {
	column, err := query.Columns()
	if err != nil {
		return nil, err
	}
	values := make([][]byte, len(column))
	scans := make([]interface{}, len(column))
	for i := range values {
		scans[i] = &values[i]
	}
	results := []map[string]string{}
	i := 0
	for query.Next() {
		if err := query.Scan(scans...); err != nil {
			return nil, err
		}
		row := make(map[string]string)
		for k, v := range values {
			key := column[k]
			row[key] = string(v)
		}
		results = append(results, row)
		i++
	}
	return results, nil
}

func GetFieldsValue(fieldST Field, skipErrPrint ...bool) interface{} {
	if fieldST.FieldType == "int" {
		if fieldST.FieldDefault == "" {
			return 0
		}
		_v, err := strconv.ParseInt(fieldST.FieldDefault, 10, 32)
		if len(skipErrPrint) > 0 {
			if ChkErrNormal(err) {
				return ""
			} else {
				return _v
			}
		} else {
			if ChkErr(err) {
				return ""
			} else {
				return _v
			}
		}
	} else if fieldST.FieldType == "string" {
		return fieldST.FieldDefault
	} else if fieldST.FieldType == "bool" {
		if fieldST.FieldDefault == "" {
			return false
		}
		_v, err := strconv.ParseBool(fieldST.FieldDefault)
		if len(skipErrPrint) > 0 {
			if ChkErrNormal(err) {
				return ""
			} else {
				return _v
			}
		} else {
			if ChkErr(err) {
				return ""
			} else {
				return _v
			}
		}
	}
	return ""
}
