package hunt

import (
	"database/sql"
	"errors"
	"log"
	"os"

	_ "modernc.org/sqlite"

	"github.com/cuhsat/fox/v4/internal/pkg/types/event"
)

const schema = `
CREATE TABLE IF NOT EXISTS events (
    id INTEGER PRIMARY KEY,
    time INTEGER NOT NULL,
    host TEXT NULL,
    user TEXT NULL,
	message TEXT NULL,
	severity INTEGER NOT NULL,
	UNIQUE(time) ON CONFLICT IGNORE
);
`

const insert = `
INSERT INTO events (
	time,
	host,
	user,
	message,
	severity
) VALUES (?,?,?,?,?);
`

type Database struct {
	path string
	sql  *sql.DB
}

func NewDB(path string) *Database {
	var err error

	db := &Database{
		path: path,
	}

	db.sql, err = sql.Open("sqlite", "file:"+path)

	if err != nil {
		log.Fatal(err)
	}

	_, err = os.Stat(path)

	if errors.Is(err, os.ErrNotExist) {
		_, err = db.sql.Exec(schema)
	}

	if err != nil {
		log.Fatal(err)
	}

	return db
}

func (db *Database) String() string {
	return db.path
}

func (db *Database) Write(evt *event.Event) {
	_, err := db.sql.Exec(insert,
		evt.Time.UTC(),
		evt.Host,
		evt.User,
		evt.Message,
		evt.Severity,
	)

	if err != nil {
		log.Println(err)
	}
}
