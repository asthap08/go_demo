package storage

import 
	"github.com/asthap08/Golang_notes/internal/types"


type Storage interface {
	CreateStudent(name string, email string, age int) (int64, error)
	GetStudentByID(id int64) (types.Student, error)
	GetStudents()([]types.Student, error)
}

type Student struct {
	ID    int64
	Name  string
	Email string
	Age   int
}