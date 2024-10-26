package collect

import (
	"database/sql"
	"fmt"
	"os"
	"strconv"

	_ "github.com/lib/pq"
)

func insertAlbum(name string, artist string, image string) *int64 {
	insertStmt := `INSERT INTO public.album(name, artist, image) VALUES ($1, $2, $3) RETURNING id`
	return insertSelect(insertStmt, name, artist, image)
}

func insertMusic(name string, albumId int64, lyrics string, melon *string, genie *string, flo *string, bugs *string, vibe *string) *int64 {
	insertStmt := `INSERT INTO public.music(name, album_id, lyrics, melon, genie, flo, bugs, vibe) VALUES ($1, $2, $3, $4, $5, $6, $7, $8) RETURNING id`
	return insertSelect(insertStmt, name, albumId, lyrics, melon, genie, flo, bugs, vibe)
}

func insertRank(rank int, types string, musicId int64) *int64 {
	insertStmt := `INSERT INTO public.ranks(ranks, crdate, types, music_id) VALUES ($1, '2024-06-13', $2, $3) RETURNING id`
	return insertSelect(insertStmt, rank, types, musicId)
}

func updateMusic(types string, songNo string, musicId int64) {
	db := connection()
	defer db.Close()

	updateStmt := fmt.Sprintf(`UPDATE music SET %s = $1 WHERE id = $2`, types)
	_, err := db.Exec(updateStmt, songNo, musicId)
	stop(err)
}

func getAlbum(name string, artist string) *int64 {
	selectStmt := `SELECT id FROM album WHERE name % $1 AND artist % $2`
	return insertSelect(selectStmt, name, artist)
}

func getMusic(name string, artist string) *int64 {
	selectStmt := `SELECT m.id FROM music m INNER JOIN album a ON m.album_id = a.id WHERE m.name % $1 AND a.artist % $2`
	return insertSelect(selectStmt, name, artist)
}

func existOnChart(types string, songNo string) bool {
	selectStmt := fmt.Sprintf(`SELECT COUNT(*) FROM music WHERE %s = $1`, types)
	return *insertSelect(selectStmt, songNo) >= 1
}

func insertSelect(query string, args ...interface{}) *int64 {
	db := connection()
	defer db.Close()

	var id sql.NullInt64
	err := db.QueryRow(query, args...).Scan(&id)
	stop(err)

	if id.Valid {
		return &id.Int64
	}
	return nil
}

func connection() *sql.DB {
	port, err := strconv.Atoi(os.Getenv("PG_PORT"))
	stop(err)

	conn := fmt.Sprintf("host=%s port=%d user=%s dbname=%s password=%s",
		os.Getenv("PG_HOST"), port, os.Getenv("PG_USER"), os.Getenv("PG_DATABASE"), os.Getenv("PG_PASSWORD"),
	)

	db, err := sql.Open("postgres", conn)
	stop(err)

	return db
}

func stop(err error) {
	if err != nil && err != sql.ErrNoRows {
		panic(err)
	}
}
