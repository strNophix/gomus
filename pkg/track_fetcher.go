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
	CREATE TABLE IF NOT EXISTS track (
		title TEXT NOT NULL,
		album TEXT,
		artist TEXT NOT NULL,
		uri TEXT NOT NULL,
		trackNumber INTEGER,
		trackTotal INTEGER
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
	tx, err := ti.db.Begin()
	check(err)

	stmt, err := tx.Prepare("insert into track(title, album, artist, uri, trackNumber, trackTotal) values (?, ?, ?, ?, ?, ?)")
	check(err)

	for _, t := range tracks {
		stmt.Exec(t.Title, t.Album, t.Artist, t.TrackPath, t.TrackNumber, t.TrackTotal)
	}

	err = tx.Commit()
	check(err)

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
