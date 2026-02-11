package sqlite

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"maps"
	"os"

	_ "modernc.org/sqlite"

	"github.com/cuhsat/fox/v4/internal/pkg/data/storage"
	"github.com/cuhsat/fox/v4/internal/pkg/types/event"
)

const schema = `
CREATE TABLE IF NOT EXISTS Events (
    ID INTEGER PRIMARY KEY,
    Time INTEGER NOT NULL,
    Host TEXT NULL,
    User TEXT NULL,
	Message TEXT NULL,
	Severity INTEGER NOT NULL,
	Sequence TEXT NULL,
	Source TEXT NOT NULL,
	Category TEXT NULL,
	Service TEXT NULL
);

CREATE TABLE IF NOT EXISTS Fields (
    ID INTEGER PRIMARY KEY,
	EventID INTEGER REFERENCES Events,
    Key TEXT NOT NULL,
    Value TEXT NULL
);

CREATE INDEX EventID ON Fields(EventID);
`

const events = `
INSERT OR IGNORE INTO Events (
	Time,
	Host,
	User,
	Message,
	Severity,
	Sequence,
	Source,
	Category,
	Service
) VALUES (?,?,?,?,?,?,?,?,?);
`

const fields = `
INSERT OR IGNORE INTO Fields (
	EventID,	
	Key,
	Value
) VALUES (?,?,?);
`

type Database struct {
	path string
	sql  *sql.DB
}

func New(name string) storage.Storage {
	var err error

	name = fmt.Sprintf("%s.sqlite", name)

	db := &Database{path: name}
	db.sql, err = sql.Open("sqlite", fmt.Sprintf("file:%s", name))

	if err != nil {
		log.Fatalln(err)
	}

	_, err = os.Stat(name)

	if errors.Is(err, os.ErrNotExist) {
		_, err = db.sql.Exec(schema)
	}

	if err != nil {
		log.Fatalln(err)
	}

	return db
}

func (db *Database) String() string {
	return db.path
}

func (db *Database) Write(evt *event.Event) error {
	res, err := db.sql.Exec(events,
		evt.Time.UTC(),
		evt.Host,
		evt.User,
		evt.Message,
		evt.Severity,
		evt.Sequence,
		evt.Source,
		evt.Category,
		evt.Service,
	)

	if err != nil {
		return err
	}

	if len(evt.Fields) > 0 {
		id, err := res.LastInsertId()

		if err != nil {
			return err
		}

		for k, v := range maps.All(evt.Fields) {
			_, err := db.sql.Exec(fields, id, k, v)

			if err != nil {
				return err
			}
		}
	}

	return nil
}

func (db *Database) Close() error {
	return nil
}
