package task_repo

import (
	"log"
	"os"
)

func init() {
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)
	log.SetPrefix("[task_repo/GetTasks] ")
	log.SetOutput(os.Stdout)
}
func (tr *TaskRepo) GetTasks() ([]Task, error) {
	var tasks []Task

	query := "SELECT id, date, title, comment, repeat FROM scheduler ORDER BY date LIMIT 10"
	rows, err := tr.DB.Query(query)
	if err != nil {
		log.Printf("Ошибка при выполнении запроса к базе данных: %v", err)
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var t Task
		err := rows.Scan(&t.ID, &t.Date, &t.Title, &t.Comment, &t.Repeat)
		if err != nil {
			log.Printf("Ошибка при чтении данных из базы данных: %v", err)
			return nil, err
		}
		tasks = append(tasks, t)
	}

	err = rows.Err()
	if err != nil {
		log.Printf("Ошибка при чтении данных из базы данных: %v", err)
		return nil, err
	}

	if len(tasks) == 0 {
		tasks = []Task{}
	}
	return tasks, nil

}
