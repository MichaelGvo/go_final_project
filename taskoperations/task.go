package taskoperations

type Task struct {
	ID      string `json:"id"`      // id задачи
	Date    string `json:"date"`    // дата дедлайна
	Title   string `json:"title"`   // заголовок
	Comment string `json:"comment"` // комментарий
	Repeat  string `json:"repeat"`  // правило, по которому задачи будут повторяться
}
