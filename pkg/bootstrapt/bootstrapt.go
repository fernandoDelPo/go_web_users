package bootstrapt

import (
	_ "github.com/go-sql-driver/mysql"
	//"github.com/fernandoDelPo/go_web_users/internal/domain"
	//"github.com/fernandoDelPo/go_web_users/internal/user"
	"database/sql"
	"log"
	"os"
)

func NewLogger() *log.Logger {
	return log.New(os.Stdout, "", log.LstdFlags|log.Lshortfile)
}

func NewDB() (*sql.DB, error) {

	db, err := sql.Open("mysql", "root:root@tcp(127.0.0.1:3336)/go_course_users")

	if err != nil {
		return nil, err
	}

	return  db, nil	
}
