package sqlite

import(
	"database/sql"

	"github.com/asthap08/Golang_notes/internal/config"


	//"_" is used bcoz this driver is used behind the scenes and imported only for initialization
	_ "github.com/mattn/go-sqlite3"


)

type Sqlite struct {

	Db *sql.DB
}

func New(cfg *config.Config) (*Sqlite, error) {
	db, err := sql.Open("sqlite3",cfg.StoragePath)
	if err != nil {
		return nil, err
	}

	_,err = db.Exec(`CREATE TABLE IF NOT EXISTS students(
	id INTEGER PRIMARY KEY AUTOINCREMENT,
	name TEXT,
	email TEXT,
	age INTEGER
	)`)

	if err != nil {
		return nil, err
	}

	return &Sqlite{
		Db:db,
	
	}, nil

}