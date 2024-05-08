package main

import (
	"database/sql"
	"fmt"
)

type Album struct {
	ID     int64
	Title  string
	Artist string
	Price  float32
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
