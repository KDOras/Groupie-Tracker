package Game

import (
	"database/sql"
)

func AddPoint(db *sql.DB, userId int) {
	data, _ := db.Query(`SELECT score from ROOM_USERS WHERE id_user=?`, userId)
	userScore := 0
	for data.Next() {
		data.Scan(&userScore)
	}
	userScore++

	db.Exec(`UPDATE ROOM_USERS SET score=? WHERE id_user=?`, userScore, userId)
}
