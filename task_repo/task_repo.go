package task_repo

import "database/sql"

type TaskRepo struct {
	DB *sql.DB
}
