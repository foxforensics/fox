package sqlite

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"maps"
	"os"

	_ "modernc.org/sqlite"

	"github.com/cuhsat/fox/v4/internal/pkg/data/store"
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

CREATE UNIQUE INDEX idx_Events ON Events(Time, Host, Sequence);

CREATE TABLE IF NOT EXISTS Fields (
    ID INTEGER PRIMARY KEY,
	EventID INTEGER REFERENCES Events,
    Key TEXT NOT NULL,
    Value TEXT NULL
);

CREATE UNIQUE INDEX ix_Fields ON Fields(EventID, Key);
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

func New(name string) store.Store {
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

func (db *Database) Store(evt *event.Event) error {
	tx, err := db.sql.Begin()

	if err != nil {
		return err
	}

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
		_ = tx.Rollback()
		return err
	}

	if len(evt.Fields) > 0 {
		id, err := res.LastInsertId()

		if err != nil {
			_ = tx.Rollback()
			return err
		}

		for k, v := range maps.All(evt.Fields) {
			_, err := db.sql.Exec(fields, id, k, v)

			if err != nil {
				_ = tx.Rollback()
				return err
			}
		}
	}

	return tx.Commit()
}

func (db *Database) Close() error {
	return nil
}
