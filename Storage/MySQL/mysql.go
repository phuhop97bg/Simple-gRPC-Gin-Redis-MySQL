package Storage

import (
	"errors"
	"os"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"golang.org/x/net/context"

	"log"
)

func NewMySQLConnect(ctx context.Context) *sqlx.DB {
	//dsn := "dbUser:dbPassword@(dbURL:PORT)/dbName"
	dsn := "root:1234@(localhost:3306)/TestDB"
	db, err := sqlx.ConnectContext(ctx, "mysql", dsn)
	if err != nil {
		log.Fatalln(err)
	}

	if errors.Is(err, context.DeadlineExceeded) {
		log.Println("ContextDeadlineExceeded: true")
	}
	if os.IsTimeout(err) {
		log.Println("IsTimeoutError: true")

	}

	return db
}
