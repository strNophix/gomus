package gomus

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/faiface/beep"
	bflac "github.com/faiface/beep/flac"
	"github.com/mewkiz/flac"
	"github.com/mewkiz/flac/meta"
)

type track struct {
	Name      string
	Artist    string
	TrackPath string
}

func (t track) FilterValue() string { return "" }
func (t track) fullName() string    { return fmt.Sprintf("%s - %s", t.Artist, t.Name) }

func (t track) GetStream() (beep.StreamSeekCloser, beep.Format, error) {
	f, err := os.Open(t.TrackPath)
	check(err)

	if strings.HasSuffix(t.TrackPath, ".flac") {
		streamer, format, err := bflac.Decode(f)
		check(err)
		return streamer, format, nil
	}

	return nil, beep.Format{}, fmt.Errorf("Could not parse track")
}

func TrackFromFlac(path string) track {
	s, err := flac.ParseFile(path)
	check(err)

	t := track{TrackPath: path}
	for _, block := range s.Blocks {
		if block.Header.Type == meta.TypeVorbisComment {
			c := block.Body.(*meta.VorbisComment)
			for _, tag := range c.Tags {
				if tag[0] == "ARTIST" {
					t.Artist = tag[1]
				} else if tag[0] == "TITLE" {
					t.Name = tag[1]
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
			tracks = append(tracks, TrackFromFlac(p))
		}
	}

	return trackIndex{tracks: tracks}
}
