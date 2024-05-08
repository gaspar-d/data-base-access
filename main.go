package main

// https://go.dev/doc/tutorial/data-access
import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/lib/pq"
)

var db *sql.DB

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
		Title:  "Miles Davis - Requiem In D Minor",
		Artist: "Miles Davis",
		Price:  59.99,
	})

	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("ID of added album: %v\n", albumID)

	// if err = removeAlbum(5); err != nil {
	// 	log.Fatal(err)
	// }
}
