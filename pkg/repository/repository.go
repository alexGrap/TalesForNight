package repository

import (
	"fmt"
	"fsm/config"
	"fsm/internal/models"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"log"
)

type ConnectBd struct {
	Database *sqlx.DB
}

var Connection ConnectBd

func InitPsqlDB(c *config.Config) (*sqlx.DB, error) {
	connectionUrl := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		*c.Postgres.Host,
		*c.Postgres.Port,
		*c.Postgres.User,
		*c.Postgres.Password,
		*c.Postgres.DbName)
	fmt.Println(connectionUrl)
	database, err := sqlx.Connect("postgres", connectionUrl)
	if err != nil {
		return nil, err
	}
	return database, nil
}

func InitTables() error {
	_, err := Connection.Database.Exec(`CREATE TABLE users (
    	userid BIGSERIAL PRIMARY KEY,
    	genre TEXT NOT NULL,
    	sounder TEXT NOT NULL ,
    	book TEXT NOT NULL,
    	counter INT DEFAULT 0
		);`)
	if err != nil {
		log.Println(err)
	}
	return err
}

func DropTable() error {
	_, err := Connection.Database.Exec(`DROP TABLE users`)
	if err != nil {
		log.Fatal(err)
	}
	return err
}

func CreateUser(body models.User) {
	Connection.Database.QueryRowx(`INSERT INTO users
	( userid, genre, sounder, book, counter)
	VALUES ($1, $2, $3, $4, $5)`, body.UserId, "Сказка", "Текст", body.Book, body.Counter)
}

func GetUser(userId int64) models.User {
	var body models.User
	err := Connection.Database.Get(&body, `SELECT * FROM users WHERE userid = $1`, userId)
	if err != nil {
		log.Println(err)
	}
	return body
}

func UpdateSounder(id int64, sounder string) error {
	_, err := Connection.Database.Queryx("UPDATE users SET sounder = $1 WHERE userid = $2", sounder, id)
	if err != nil {
		log.Println(err)
	}
	return err
}

func UpdateGenre(id int64, genre string) error {
	_, err := Connection.Database.Queryx("UPDATE users SET genre = $1 WHERE userid = $2", genre, id)
	if err != nil {
		log.Println(err)
	}
	return err
}

func UpdateBook(id int64, title string) error {
	_, err := Connection.Database.Queryx("UPDATE users SET book = $1 WHERE userid = $2", title, id)
	if err != nil {
		log.Println(err)
	}
	return err
}

func UpdateCounter(id int64, count int) error {
	var body models.User
	err := Connection.Database.Get(&body, "SELECT * FROM users WHERE userid = $1;", id)
	if err != nil {
		log.Println(err)
	}
	if count == 1 {
		body.Counter++
	} else {
		body.Counter = 0
	}

	_, err = Connection.Database.Queryx("UPDATE users SET counter = $1 WHERE userid = $2", body.Counter, id)
	if err != nil {
		log.Println(err)
	}
	return err
}

func GetAdminInfo() (int, int) {
	var count int
	Connection.Database.Get(&count, "SELECT count(*) FROM users")

	var rows *sqlx.Rows
	var err error
	var tmp models.User
	yandexCount := 0
	rows, err = Connection.Database.Queryx(`SELECT * FROM users`)
	if err != nil {
		log.Println(err)
	}
	for rows.Next() {
		if err := rows.Scan(&tmp.UserId, &tmp.Genre, &tmp.Sounder, &tmp.Book, &tmp.Format, &tmp.Counter); err != nil {
			log.Println(err)
		}
		yandexCount += tmp.Counter
	}

	return count, yandexCount
}

func GetAllId(flag bool) []int64 {
	var rows *sqlx.Rows
	var err error
	var result []int64
	var tmp models.User
	rows, err = Connection.Database.Queryx(`SELECT * FROM users`)
	if err != nil {
		log.Println(err)
	}
	for rows.Next() {
		if err := rows.Scan(&tmp.UserId, &tmp.Genre, &tmp.Sounder, &tmp.Book, &tmp.Format, &tmp.Counter); err != nil {
			log.Println(err)
		}
		result = append(result, tmp.UserId)
		if flag {
			UpdateCounter(tmp.UserId, 0)
		}
	}
	return result
}
