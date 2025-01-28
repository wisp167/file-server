package storage

import (
	"fmt"
	"database/sql"
	"context"
	//_ "github.com/go-sql-driver/mysql"
	_ "github.com/lib/pq"
	"github.com/wisp167/file-server/internal/dbQueries"
	"github.com/google/uuid"
	"strings"
)

func (s MySqlStorage) Select(ctx context.Context, columns []string) ([]dbQueries.File, error){
	query := fmt.Sprintf("SELECT %s FROM files", strings.Join(columns, ", "))
	rows, err := s.rawDB.QueryContext(ctx, query)
	if(err != nil){
		return nil, fmt.Errorf("Failed to execute query: %w\n", err)
	}
	defer rows.Close()
	var files []dbQueries.File
	for rows.Next(){
		var b dbQueries.File
		
		scanArgs := make([]interface{}, len(columns))
		for i, col := range columns{
			switch col{
				case "id":
					scanArgs[i] = &b.ID
				case "file_name":
					scanArgs[i] = &b.FileName
				case "file_data":
					scanArgs[i] = &b.FileData
				case "create_time":
					var createTime sql.NullTime
					scanArgs[i] = &createTime
					if createTime.Valid{
						b.CreateTime = createTime.Time
					}
				default:
					return nil, fmt.Errorf("Uknown column in select query: %w\n", err)
			}
			
			
		}
		err := rows.Scan(scanArgs...)
		if err != nil {
			return nil, fmt.Errorf("Failed to scan row: %w\n", err)
		}
		for i, col := range columns {
			if col == "create_time" {
				var createTime sql.NullTime
				createTime = *(scanArgs[i].(*sql.NullTime))
				if createTime.Valid {
					b.CreateTime = createTime.Time
				}
			}
		}
		files = append(files, b)
	}
	if err = rows.Err(); err != nil{
		return nil, fmt.Errorf("Error during iterations %w\n", err)
	}
	return files, nil
}




func (s MySqlStorage) GetFile(ctx context.Context, id uuid.UUID) (dbQueries.File, error) {
	query := `SELECT id, file_name, file_data, create_time, update_time
	          FROM files
	          WHERE id = $1`
	row := s.rawDB.QueryRowContext(ctx, query, id)
	var file dbQueries.File
	err := row.Scan(&file.ID, &file.FileName, &file.FileData, &file.CreateTime, &file.UpdateTime)
	return file, err
}
