package db

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "modernc.org/sqlite"
)

var DB *sql.DB

const createTableSQL = `
CREATE TABLE IF NOT EXISTS scheduler (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    date CHAR(8) NOT NULL DEFAULT "",
    title VARCHAR(255) NOT NULL DEFAULT "",
    comment TEXT,
    repeat VARCHAR(128) NOT NULL DEFAULT ""
);`

const createIndexSQL = `CREATE INDEX IF NOT EXISTS idx_scheduler_date ON scheduler(date);`

func Init(dbFile string) error {
	var install bool
	if _, err := os.Stat(dbFile); err != nil {
		if os.IsNotExist(err) {
			install = true
		} else {
			return fmt.Errorf("ошибка проверки файла БД: %w", err)
		}
	}

	conn, err := sql.Open("sqlite", dbFile)
	if err != nil {
		return fmt.Errorf("ошибка открытия БД: %w", err)
	}
	if err := conn.Ping(); err != nil {
		return fmt.Errorf("ошибка подключения к БД: %w", err)
	}
	DB = conn

	if _, err := DB.Exec(createTableSQL); err != nil {
		return fmt.Errorf("ошибка создания таблицы: %w", err)
	}
	if _, err := DB.Exec(createIndexSQL); err != nil {
		return fmt.Errorf("ошибка создания индекса: %w", err)
	}

	if install {
		log.Println("База данных создана и инициализирована.")
	}
	return nil
}
