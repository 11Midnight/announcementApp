package dbOp

import (
	"database/sql"
	"github.com/Masterminds/squirrel"
)

//Структура того, что читаем с базы данных.
type taskstr struct {
	id   int
	Date string
	Task string
}

func InsertTask(db *sql.DB, registerDate, registerTask string) error {
	sql, args, err := squirrel.Insert("tasks").Columns("date", "task").Values(registerDate, registerTask).ToSql()
	_, err = db.Exec(sql, args[0], args[1])
	return err

}

func Connect(dataSoureName string) (*sql.DB, error) {
	db, err := sql.Open("mysql", dataSoureName)
	return db, err
}

//Считывает таски в структуру.
func ReadTasks(db *sql.DB) (tasksw []taskstr, err error) {
	sql, _, err := squirrel.Select("*").From("golang.tasks").ToSql()
	if err != nil {
		return nil, err
	}
	//rows, err := db.Query("select * from golang.tasks")
	rows, err := db.Query(sql)
	if err != nil {
		return nil, err
	}

	defer rows.Close()
	tasks := []taskstr{}
	for rows.Next() {
		t := taskstr{}
		err := rows.Scan(&t.id, &t.Date, &t.Task)
		if err != nil {
			return nil, err
		}
		tasks = append(tasks, t)
	}
	return tasks, nil
}
