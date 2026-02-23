package database

func DeleteSession(sessionId string) error {
	// return database.Database.QueryRow(query).Err()
	// db.exec vs db.queryrow in golang sqlite
	// queryrow not working with delete statement
	_, err := Database.Exec(
		"DELETE FROM sessions WHERE id = ?",
		sessionId) // returns result
	return err
}
