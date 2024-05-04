package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/lib/pq"
)

var db *sql.DB

type Album struct {
	ID     int64
	Title  string
	Artist string
	Price  float32
}

func main() {
	cfg := struct {
		User     string
		Password string
		Host     string
		Port     int
		Database string
	}{
		User:     os.Getenv("DBUSER"),
		Password: os.Getenv("DBPASS"),
		Host:     "localhost",
		Port:     5432,
		Database: "recordings",
	}

	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		cfg.Host, cfg.Port, cfg.User, cfg.Password, cfg.Database)

	var err error
	db, err = sql.Open("postgres", psqlInfo)
	if err != nil {
		log.Fatal(err)
	}

	if err = db.Ping(); err != nil {
		log.Fatal(err)
	}
	fmt.Println("Connected to Postgres!")

	albums, err := albumsByArtist("John Coltrane")
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Albums found:", albums)

	album, err := albumsByID(3)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Album found:", album)

	albumID, err := addAlbum(Album{
		Title:  "The Modern Sound of Betty Carter",
		Artist: "Betty Carter",
		Price:  49.99,
	})

	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("ID of added album: %v\n", albumID)

	// if err = removeAlbum(9); err != nil {
	// 	log.Fatal(err)
	// }
}

// get all albums for an artist
func albumsByArtist(name string) ([]Album, error) {
	var albums []Album

	rows, err := db.Query("SELECT * FROM album WHERE artist = $1", name)
	if err != nil {
		return nil, fmt.Errorf("albumsByArtist %q: %v", name, err)
	}

	// try to move this defer to the bottom of the function, before the return
	defer rows.Close()
	for rows.Next() {
		var album Album
		if err := rows.Scan(&album.ID, &album.Title, &album.Artist, &album.Price); err != nil {
			return nil, fmt.Errorf("albumsByArtist %q: %v", name, err)
		}
		albums = append(albums, album)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("albumsByArtist %q: %v", name, err)
	}

	return albums, nil
}

// get one album by id
func albumsByID(id int64) (Album, error) {
	var album Album

	row := db.QueryRow("SELECT * FROM album WHERE id = $1", id)

	if err := row.Scan(&album.ID, &album.Title, &album.Artist, &album.Price); err != nil {
		if err == sql.ErrNoRows {
			return album, fmt.Errorf("albumsById %d: no such album", id)
		}
		return album, fmt.Errorf("albumsById %d: %v", id, err)
	}
	return album, nil
}

func addAlbum(album Album) (int64, error) {
	var id int64
	err := db.QueryRow(
		"INSERT INTO album (title, artist, price) VALUES ($1, $2, $3) RETURNING id", album.Title, album.Artist, album.Price).Scan(&id)

	if err != nil {
		return 0, fmt.Errorf("addAlbum: %v", err)
	}

	return id, nil
}

func removeAlbum(id int64) error {
	_, err := db.Exec("DELETE FROM album WHERE id = $1", id)
	if err != nil {
		return fmt.Errorf("removeAlbum: %v", err)
	}

	fmt.Printf("Album %d deleted\n", id)
	return nil
}
