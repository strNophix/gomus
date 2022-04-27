package gomus

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/faiface/beep"
	bflac "github.com/faiface/beep/flac"
	"github.com/mewkiz/flac"
	"github.com/mewkiz/flac/meta"
)

type track struct {
	Title       string
	Album       string
	Artist      string
	TrackPath   string
	TrackTotal  uint8
	TrackNumber uint8
}

func (t track) FilterValue() string { return "" }
func (t track) fullName() string    { return fmt.Sprintf("%s - %s", t.Artist, t.Title) }

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
	path, err := filepath.Abs(path)
	check(err)

	s, err := flac.ParseFile(path)
	check(err)

	t := track{TrackPath: path}
	for _, block := range s.Blocks {
		if block.Header.Type == meta.TypeVorbisComment {
			c := block.Body.(*meta.VorbisComment)
			for _, tagTuple := range c.Tags {
				tag, val := tagTuple[0], tagTuple[1]
				switch tag {
				case "TITLE":
					t.Title = val
				case "ARTIST":
					t.Artist = val
				case "ALBUM":
					t.Album = val
				case "TRACKNUMBER":
					trackNum, err := strconv.ParseUint(val, 10, 8)
					check(err)
					t.TrackNumber = uint8(trackNum)
				case "TRACKTOTAL":
					trackTotal, err := strconv.ParseUint(val, 10, 8)
					check(err)
					t.TrackTotal = uint8(trackTotal)
				}
			}
		}
	}
	return t
}
