
package storage

import (
	"fmt"
	"database/sql"
	//_ "github.com/go-sql-driver/mysql"
	_ "github.com/lib/pq"
	"github.com/wisp167/file-server/internal/dbQueries"
)

type MySqlStorage struct{
	DB *dbQueries.Queries
	rawDB *sql.DB
}

type MySqlConfig struct{
	Username string
	Password string
	DbName string
	Port uint
	Host string
}

func NewMySqlStorage(conf MySqlConfig) (MySqlStorage, error){
	dsn := fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=disable", conf.Username, conf.Password, conf.Host, conf.Port, conf.DbName)
	db, err := sql.Open("postgres", dsn)
	if err != nil{
		return MySqlStorage{}, fmt.Errorf("impossible to open %w\n", err)
	}
	
	err = db.Ping()
	if err != nil{
		return MySqlStorage{}, fmt.Errorf("impossible to ping %w\n", err)	
	}
	return MySqlStorage{
		DB: dbQueries.New(db),
		rawDB: db,
	}, nil
}

/*
func (s MySqlStorage) Create(ctx context.Context, b book.Book) (book.Book, error){
	query:= "INSERT INTO book (create_time, name, author_name) VALUES (?, ?, ?)"
	InsertResult, err := s.db.ExecContext(ctx, query, b.CreateTime, b.Name, b.AuthorName)
	if err != nil{
		return b, fmt.Errorf("error while insert: %w\n", err)
	}
	id, err := InsertResult.LastInsertId()
	if err != nil{
		return b, fmt.Errorf("error while getint last insert id: %w\n", err)
	}
	b.ID = id
	return b, err
}

*/

