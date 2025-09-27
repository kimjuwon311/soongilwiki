package tool

import (
	"database/sql"
	"log"
	"strings"
	"time"

	_ "github.com/go-sql-driver/mysql"
	jsoniter "github.com/json-iterator/go"
	_ "modernc.org/sqlite"
)

var db_set = make(map[string]string) // 전역 설정 저장

// DB 초기화
func DB_init(get_db_set string) {
	var json = jsoniter.ConfigCompatibleWithStandardLibrary

	other_set := map[string]string{}
	err := json.Unmarshal([]byte(get_db_set), &other_set)
	if err != nil {
		log.Printf("DB_init: JSON unmarshal error: %v", err)
		return
	}

	// 필수 필드 검사
	required := []string{"db_type", "db_name"}
	for _, key := range required {
		if val, ok := other_set[key]; !ok || val == "" {
			log.Printf("DB_init: missing required field: %s", key)
			return
		}
	}

	for k, v := range other_set {
		db_set[k] = v
	}
}

// DB 연결
func DB_connect() *sql.DB {
	dbType := db_set["db_type"]

	if dbType == "" {
		log.Panic("DB_connect: db_type is not set in db_set")
	}

	switch dbType {
	case "sqlite":
		dbName := db_set["db_name"]
		if dbName == "" {
			log.Panic("DB_connect: db_name is not set for sqlite")
		}
		db, err := sql.Open("sqlite", dbName+".db?_journal_mode=WAL&_busy_timeout=5000")
		if err != nil {
			log.Panicf("DB_connect: sqlite open error: %v", err)
		}
		return db

	case "mysql":
		required := []string{"db_mysql_user", "db_mysql_pw", "db_mysql_host", "db_mysql_port", "db_name"}
		for _, key := range required {
			if db_set[key] == "" {
				log.Panicf("DB_connect: missing MySQL config field: %s", key)
			}
		}
		dsn := db_set["db_mysql_user"] + ":" + db_set["db_mysql_pw"] + "@tcp(" + db_set["db_mysql_host"] + ":" + db_set["db_mysql_port"] + ")/" + db_set["db_name"]
		db, err := sql.Open("mysql", dsn)
		if err != nil {
			log.Panicf("DB_connect: mysql open error: %v", err)
		}
		return db

	default:
		log.Panicf("DB_connect: unknown db_type: %s", dbType)
		return nil
	}
}

// DB 닫기
func DB_close(db *sql.DB) {
	db.Close()
}

// DB 타입 가져오기
func Get_DB_type() string {
	return db_set["db_type"]
}

// SQL 문법 변경 (MySQL 호환용)
func DB_change(data string) string {
	if Get_DB_type() == "mysql" {
		data = strings.ReplaceAll(data, "random()", "rand()")
		data = strings.ReplaceAll(data, "collate nocase", "collate utf8mb4_general_ci")
	}
	return data
}

// 쿼리 실행
func Exec_DB(db *sql.DB, query string, values ...any) {
	const retryDelay = 10 * time.Millisecond

	stmt, err := db.Prepare(DB_change(query))
	if err != nil {
		panic(err)
	}
	defer stmt.Close()

	for {
		_, err = stmt.Exec(values...)
		if err == nil {
			return
		}
		if strings.Contains(err.Error(), "database is locked") {
			time.Sleep(retryDelay)
			continue
		}
		panic(err)
	}
}

// 다중 행 쿼리
func QueryRow_DB(db *sql.DB, query string, var_list []any, values ...any) bool {
	const retryDelay = 10 * time.Millisecond

	stmt, err := db.Prepare(DB_change(query))
	if err != nil {
		panic(err)
	}
	defer stmt.Close()

	for {
		row := stmt.QueryRow(values...)
		err := row.Scan(var_list...)
		switch err {
		case nil:
			return true
		case sql.ErrNoRows:
			return false
		}
		if strings.Contains(err.Error(), "database is locked") {
			time.Sleep(retryDelay)
			continue
		}
		panic(err)
	}
}

// 다중 결과 쿼리 함수
func Query_DB(db *sql.DB, query string, values ...any) (*sql.Rows, error) {
	const retryDelay = 10 * time.Millisecond

	stmt, err := db.Prepare(DB_change(query))
	if err != nil {
		return nil, err
	}
	// stmt는 직접 사용자가 닫도록 리턴 후 defer를 생략
	for {
		rows, err := stmt.Query(values...)
		if err == nil {
			return rows, nil
		}
		if strings.Contains(err.Error(), "database is locked") {
			time.Sleep(retryDelay)
			continue
		}
		return nil, err
	}
}






