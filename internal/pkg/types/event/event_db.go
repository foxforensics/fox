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
	Source INTEGER NOT NULL
);

CREATE TABLE IF NOT EXISTS Extensions (
    ID INTEGER PRIMARY KEY,
	EventID INTEGER REFERENCES Events,
    Key TEXT NOT NULL,
    Value TEXT NULL
);

CREATE INDEX EventID ON Extensions(EventID);
`

const events = `
INSERT OR IGNORE INTO Events (
	Time,
	Host,
	User,
	Message,
	Severity,
    Source
) VALUES (?,?,?,?,?,?);
`

const extensions = `
INSERT OR IGNORE INTO Extensions (
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

func (db *Database) Upsert(evt *Event) {
	res, err := db.sql.Exec(events,
		evt.Time.UTC(),
		evt.Host,
		evt.User,
		evt.Message,
		evt.Severity,
		evt.Source,
	)

	if err != nil {
		log.Println(err)
		return
	}

	if len(evt.Extension) > 0 {
		id, err := res.LastInsertId()

		if err != nil {
			log.Println(err)
			return
		}

		for k, v := range maps.All(evt.Extension) {
			_, err := db.sql.Exec(extensions, id, k, v)

			if err != nil {
				log.Println(err)
			}
		}
	}
}
