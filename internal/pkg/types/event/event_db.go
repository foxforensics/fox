package event

import (
	"database/sql"
	"errors"
	"log"
	"maps"
	"os"

	_ "modernc.org/sqlite"
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

func NewDB(path string) *Database {
	var err error

	db := &Database{path: path}
	db.sql, err = sql.Open("sqlite", "file:"+path)

	if err != nil {
		log.Fatalln(err)
	}

	_, err = os.Stat(path)

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

func (db *Database) Upsert(evt *Event) error {
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
