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
		INSERT INTO requests (chat_id, image_path_name, status_code, login)
		VALUES ($1, $2, $3, $4);`

	checkBalance = `
		SELECT balance FROM users
		WHERE login = $1;`

	addBalance = `
		UPDATE users
		SET balance = balance + $1
		WHERE login = $2;`

	debiting = `
		UPDATE users
		SET balance = balance - 1
		WHERE login = $1;`

	newTaskWEB = `
		INSERT INTO requests (chat_id, image_path_name, status_code, login)
		VALUES ($1, $2, $3, $4);`

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
		WHERE status_code = 3 AND chat_id != 0;`

	getUserAnswers = `
		SELECT id, detected_path_name, description
		FROM requests
		WHERE chat_id = 0 AND login = $1;`

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
