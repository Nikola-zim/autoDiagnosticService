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
		INSERT INTO requests (chatID, image_path_name, status_code)
		VALUES ($1, $2, $3);`

	getTasks = `
		SELECT id, chatid, image_path_name FROM requests
		WHERE status_code = $1
		LIMIT $2;`

	changeToRecognized = `
		UPDATE requests
		SET detected_path_name = $1,
		description = $2,
		status_code = $3
		WHERE id = $4;`

	getAnswers = `
		SELECT id, chatid, detected_path_name, description
		FROM requests
		WHERE status_code = 3;`

	changeToDone = `
		UPDATE requests
		SET status_code = 4
		WHERE id in (%s);`
)
