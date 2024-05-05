package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"
)

// Config contains all the necessary data to establish a connection to the database.
type Config struct {
	User     string
	Password string
	Host     string
	Port     int
	Database string
}

// Database is a struct that holds the established connection to the database.
type Database struct {
	conn *sql.DB
}

// AlbumsRepository is an interface that defines the methods needed to interact with the database.
type AlbumsRepository interface {
	GetByArtist(name string) ([]Album, error)
	GetById(id int64) (Album, error)
	Add(album Album) (int64, error)
	Delete(id int64) error
}

// PsqlAlbumsRepository is a struct that implements the AlbumsRepository interface and contains a reference to the database.
type PsqlAlbumsRepository struct {
	db *Database
}

// Album is a struct that represents an album in the database.
type Album struct {
	ID     int64
	Title  string
	Artist string
	Price  float32
}

// GetByArtist retrieves albums from the database based on the artist name provided.
func (r *PsqlAlbumsRepository) GetByArtist(name string) ([]Album, error) {
	// Query to retrieve albums based on the artist name.
	query := "SELECT * FROM album WHERE artist = $1"
	return r.fetchAlbums(query, name)
}

// GetById retrieves an album from the database based on its ID.
func (r *PsqlAlbumsRepository) GetById(id int64) (Album, error) {
	// Query to retrieve an album based on its ID.
	query := "SELECT * FROM album WHERE id = $1"
	return r.fetchAlbum(query, id)
}

// Add inserts a new album into the database and returns its ID.
func (r *PsqlAlbumsRepository) Add(album Album) (int64, error) {
	// Query to insert a new album into the database.
	query := "INSERT INTO album (title, artist, price) VALUES ($1, $2, $3) RETURNING id"
	return r.insertAlbum(query, album)
}

// Delete deletes an album from the database based on its ID.
func (r *PsqlAlbumsRepository) Delete(id int64) error {
	// Query to delete an album from the database.
	query := "DELETE FROM album WHERE id = $1"
	_, err := r.db.conn.Exec(query, id)
	if err != nil {
		return fmt.Errorf("removeAlbum: %v", err)
	}

	return nil
}

// fetchAlbums retrieves multiple albums from the database based on a query and name.
func (r *PsqlAlbumsRepository) fetchAlbums(query, name string) ([]Album, error) {
	rows, err := r.db.conn.Query(query, name)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var albums []Album
	for rows.Next() {
		var album Album
		err := rows.Scan(&album.ID, &album.Title, &album.Artist, &album.Price)
		if err != nil {
			return nil, err
		}
		albums = append(albums, album)
	}

	return albums, rows.Err()
}

// fetchAlbum retrieves a single album from the database based on a query and ID.
func (r *PsqlAlbumsRepository) fetchAlbum(query string, id int64) (Album, error) {
	row := r.db.conn.QueryRow(query, id)
	var album Album
	err := row.Scan(&album.ID, &album.Title, &album.Artist, &album.Price)
	if err == sql.ErrNoRows {
		return album, fmt.Errorf("albumsById %d: no such album", id)
	}

	if err != nil {
		return album, err
	}

	return album, nil
}

// insertAlbum inserts a new album into the database and returns its ID.
func (r *PsqlAlbumsRepository) insertAlbum(query string, album Album) (int64, error) {
	var id int64
	err := r.db.conn.QueryRow(query, album.Title, album.Artist, album.Price).Scan(&id)
	if err != nil {
		return 0, fmt.Errorf("addAlbum: %v", err)
	}

	return id, nil
}

// main is the entry point of the program. It establishes a connection to the database and performs some operations.
func main() {
	// Configuration for the database.
	config := Config{
		User:     os.Getenv("DBUSER"),
		Password: os.Getenv("DBPASS"),
		Host:     "localhost",
		Port:     5432,
		Database: "recordings",
	}
	// Establish a connection to the database.
	db, err := ConnectToDB(config)
	if err != nil {
		log.Fatal(err)
	}

	// Create a new PsqlAlbumsRepository instance with the established connection.
	repo := PsqlAlbumsRepository{db: db}

	// Retrieve albums by artist.
	albums, err := repo.GetByArtist("John Coltrane")
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Albums found:", albums)

	// Retrieve an album by its ID.
	album, err := repo.GetById(3)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Album found:", album)

	// Add a new album to the database and retrieve its ID.
	albumID, err := repo.Add(Album{
		Title:  "The Modern Sound of Betty Carter",
		Artist: "Betty Carter",
		Price:  49.99,
	})
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("ID of added album: %v\n", albumID)
}

// ConnectToDB establishes a connection to the database using the provided configuration.
func ConnectToDB(config Config) (*Database, error) {
	connStr := psqlInfo(config)
	conn, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, err
	}

	err = conn.Ping()
	if err != nil {
		return nil, err
	}

	return &Database{conn: conn}, nil
}

// psqlInfo generates a connection string for the database based on the provided configuration.
func psqlInfo(config Config) string {
	return fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		config.Host, config.Port, config.User, config.Password, config.Database)
}
