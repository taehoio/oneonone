package main

import (
	"database/sql"
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
	"time"

	"github.com/go-sql-driver/mysql"
)

type Question struct {
	ID        uint64 `json:"id"`
	Question  string `json:"question"`
	Category  string `json:"category"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

func main() {
	if err := upload(); err != nil {
		panic(err)
	}

	log.Println("done")
}

func parseJSONFile(filepath string) ([]*Question, error) {
	var questions []*Question
	data, err := ioutil.ReadFile(filepath)
	if err != nil {
		return nil, err
	}
	if err := json.Unmarshal(data, &questions); err != nil {
		return nil, err
	}
	return questions, nil
}

type Category struct {
	ID        int
	Name      string
	Questions []*Question
}

func extractCategoryMap(questions []*Question) map[string]*Category {
	categoryMap := make(map[string]*Category)

	for _, q := range questions {
		if _, ok := categoryMap[q.Category]; !ok {
			categoryMap[q.Category] = &Category{
				ID:        len(categoryMap),
				Name:      q.Category,
				Questions: nil,
			}
		}

		c := categoryMap[q.Category]
		c.Questions = append(c.Questions, q)
	}

	return categoryMap
}

func getEnv(key, defaultValue string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	if defaultValue == "" {
		log.Fatalf("a required environment variable missed: %s", key)
	}
	return defaultValue
}

func newMySQLDB() (*sql.DB, error) {
	mysqlCfg := mysql.Config{
		Net:                  getEnv("MYSQL_NETWORK_TYPE", "tcp"),
		Addr:                 getEnv("MYSQL_ADDRESS", "localhost:3306"),
		User:                 getEnv("MYSQL_USER", "taehoio_sa"),
		Passwd:               getEnv("MYSQL_PASSWORD", ""),
		DBName:               getEnv("MYSQL_DATABASE_NAME", "taehoio"),
		AllowNativePasswords: true,
		ParseTime:            true,
		TLSConfig:            "preferred",
	}

	db, err := sql.Open("mysql", mysqlCfg.FormatDSN())
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}

	return db, nil
}

func upload() error {
	questions, err := parseJSONFile("questions.json")
	if err != nil {
		return err
	}

	categoryMap := extractCategoryMap(questions)

	db, err := newMySQLDB()
	if err != nil {
		return err
	}

	currentTime := time.Now().UTC()
	i := 0
	for _, c := range categoryMap {
		log.Println("category:", c.ID, c.Name)

		_, err := db.Exec(
			"INSERT IGNORE INTO category (id, name, created_at, updated_at) VALUES (?, ?, ?, ?)",
			c.ID,
			c.Name,
			sql.NullTime{Time: currentTime, Valid: true},
			sql.NullTime{Time: currentTime, Valid: true},
		)
		if err != nil {
			return err
		}

		for _, q := range c.Questions {
			log.Println("question:", q.ID, q.Question)

			_, err := db.Exec(
				"INSERT IGNORE INTO question (id, question, created_at, updated_at) VALUES (?, ?, ?, ?)",
				q.ID,
				q.Question,
				sql.NullTime{Time: currentTime, Valid: true},
				sql.NullTime{Time: currentTime, Valid: true},
			)
			if err != nil {
				return err
			}

			log.Println("category_question:", c.ID, q.ID)

			_, err = db.Exec(
				"INSERT IGNORE INTO category_question (id, category_id, question_id, created_at, updated_at) VALUES (?, ?, ?, ?, ?)",
				i,
				c.ID,
				q.ID,
				sql.NullTime{Time: currentTime, Valid: true},
				sql.NullTime{Time: currentTime, Valid: true},
			)
			if err != nil {
				return err
			}

			i++
		}
	}

	return nil
}
