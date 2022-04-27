package gomus

import (
	"database/sql"
	"io/fs"
	"log"
	"path/filepath"
	"strings"

	_ "github.com/mattn/go-sqlite3"
)

type trackIndex struct {
	tracks []track
	db     *sql.DB
}

func NewDirTrackIndex(cfg ModelConfig) trackIndex {
	dbPath := filepath.Join(cfg.GomusPath, "gomus.db")
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		log.Fatalf("Failed to open gomus database: %v", err)
	}

	tblStmt := `
	CREATE TABLE IF NOT EXISTS artist (
		artist_id INTEGER PRIMARY KEY AUTOINCREMENT,
		name TEXT UNIQUE
	);

	CREATE TABLE IF NOT EXISTS album (
		album_id INTEGER PRIMARY KEY AUTOINCREMENT,
		artist_id INTEGER NOT NULL REFERENCES artist(artist_id) DEFERRABLE INITIALLY DEFERRED,
		name TEXT
	);

	CREATE UNIQUE INDEX IF NOT EXISTS artist_album ON album(artist_id, name);
	
	CREATE TABLE IF NOT EXISTS track (
		track_id INTEGER PRIMARY KEY AUTOINCREMENT,
		title TEXT NOT NULL,
		album_id INTEGER REFERENCES album(album_id) DEFERRABLE INITIALLY DEFERRED,
		artist_id INTEGER NOT NULL REFERENCES artist(artist_id) DEFERRABLE INITIALLY DEFERRED,
		uri TEXT NOT NULL UNIQUE
	);
	`

	_, err = db.Exec(tblStmt)
	if err != nil {
		log.Fatalf("%q: %s\n", err, tblStmt)
	}

	path := cfg.MusicPath
	tracks, err := ScanDirectory(path)
	check(err)

	ti := trackIndex{tracks, db}
	ti.IndexTracks(tracks)

	return ti
}

func (ti trackIndex) IndexTracks(tracks []track) error {
	var query string
	tx, err := ti.db.Begin()
	check(err)

	artistSet := make(map[string]bool)
	query = `INSERT OR IGNORE INTO artist(name) VALUES (?)`
	artistStmt, err := tx.Prepare(query)
	check(err)

	albumSet := make(map[string]bool)
	query = `
	INSERT OR IGNORE INTO album (artist_id, name)
	VALUES ((SELECT artist_id FROM artist WHERE name = ?), ?);
	`
	albumStmt, err := tx.Prepare(query)
	check(err)

	query = `
	INSERT OR IGNORE INTO track (title, uri, artist_id, album_id) 
	VALUES (?, ?, (SELECT artist_id FROM artist WHERE name = ?), (SELECT album_id FROM album WHERE name = ?));
	`
	trackStmt, err := tx.Prepare(query)
	check(err)

	for _, t := range tracks {
		if _, ok := artistSet[t.Artist]; !ok {
			artistStmt.Exec(t.Artist)
			artistSet[t.Artist] = true
		}

		key := t.Artist + ":" + t.Album
		if _, ok := albumSet[key]; !ok {
			albumStmt.Exec(t.Artist, t.Album)
			albumSet[key] = true
		}

		trackStmt.Exec(t.Title, t.TrackPath, t.Artist, t.Album)
	}

	tx.Commit()
	return nil
}

func ScanDirectory(path string) ([]track, error) {
	tracks := []track{}
	filepath.WalkDir(path, func(path string, d fs.DirEntry, err error) error {
		check(err)
		if !d.IsDir() && strings.HasSuffix(d.Name(), ".flac") {
			t := TrackFromFlac(path)
			tracks = append(tracks, t)
		}
		return nil
	})

	return tracks, nil
}
