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

func (s *Sqlite) CreateStudent(name string, email string, age int) (int64, error) {
	
	stnt,err := s.Db.Prepare("INSERT INTO students (name, email, age) VALUES(?, ?, ?)")
	if err !=nil{
		return 0, err
	}
	defer stnt.Close()

	result, err := stnt.Exec(name, email, age)
	if err!= nil{
		return 0, nil
	}

	lastId, err := result.LastInsertId()
	if err != nil{
		return 0, err
	}
	
	return lastId, nil
}