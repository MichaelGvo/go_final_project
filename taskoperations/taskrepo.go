package taskoperations

import "database/sql"

type TaskRepo struct {
	DB *sql.DB
}
