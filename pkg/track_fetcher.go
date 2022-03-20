package gomus

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/mewkiz/flac"
	"github.com/mewkiz/flac/meta"
)

type track struct {
	name      string
	artist    string
	trackPath string
}

func (t track) FilterValue() string { return "" }

func (t track) fullName() string {
	return fmt.Sprintf("%s - %s", t.artist, t.name)
}

func (t track) getReader() io.Reader {
	f, err := os.Open(t.trackPath)
	check(err)

	return f
}

func trackFromFlac(path string) track {
	s, err := flac.ParseFile(path)
	check(err)

	t := track{trackPath: path}
	for _, block := range s.Blocks {
		if block.Header.Type == meta.TypeVorbisComment {
			c := block.Body.(*meta.VorbisComment)
			for _, tag := range c.Tags {
				if tag[0] == "ARTIST" {
					t.artist = tag[1]
				} else if tag[0] == "TITLE" {
					t.name = tag[1]
				}
			}
		}
	}

	return t
}

type trackIndex struct {
	tracks []track
}

func NewDirTrackIndex(path string) trackIndex {
	files, err := ioutil.ReadDir(path)
	check(err)

	tracks := []track{}
	for _, file := range files {
		if !file.IsDir() && strings.HasSuffix(file.Name(), ".flac") {
			p := filepath.Join(path, file.Name())
			tracks = append(tracks, trackFromFlac(p))
		}
	}

	return trackIndex{tracks: tracks}
}
