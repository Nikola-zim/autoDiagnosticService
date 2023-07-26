package repo

type Status int8

const (
	Ready      Status = 1
	InWork     Status = 2
	Recognized Status = 3
	Done       Status = 4
)

const (
	newTask = `
		INSERT INTO requests (chat_id, image_path_name, status_code)
		VALUES ($1, $2, $3);`

	getTasks = `
		SELECT id, chat_id, image_path_name FROM requests
		WHERE status_code = $1
		LIMIT $2;`

	changeToRecognized = `
		UPDATE requests
		SET detected_path_name = $1,
		description = $2,
		status_code = $3
		WHERE id = $4;`

	getAnswers = `
		SELECT id, chat_id, detected_path_name, description
		FROM requests
		WHERE status_code = 3;`

	changeToDone = `
		UPDATE requests
		SET status_code = 4
		WHERE id in (%s);`

	addUser = `
		INSERT INTO users (login, password)
		VALUES ($1, $2)
		RETURNING user_id;`
	login = `
		SELECT user_id, count(user_id) FROM users
		WHERE login = $1 AND password = $2
		GROUP BY user_id;`
)
