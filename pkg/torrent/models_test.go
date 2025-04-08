package torrent

import (
	"bytes"
	"fmt"
	"os"
	"reflect"
	"testing"
)

func TestNewTorrentFromBencode(t *testing.T) {
	type args struct {
		data []byte
	}
	tests := []struct {
		name    string
		args    args
		want    *TorrentFile
		wantErr bool
	}{
		{
			name: "Normal Torrent",
			args: args{
				data: []byte("d8:announce15:http://test.com13:announce-listll2:a12:a2ee13:creation datei100e7:comment7:comment10:created by4:J2mF4:infod12:piece lengthi1000e6:pieces5:\x01\x02\x03\x04\x054:name4:Test6:lengthi5e5:filesld6:lengthi1e4:pathl6:/home/eeeee"),
			},
			want: &TorrentFile{
				Announce:     "http://test.com",
				AnnounceList: [][]string{{"a1", "a2"}},
				CreationDate: 100,
				Comment:      "comment",
				CreatedBy:    "J2mF",
				Info: InfoDict{
					PieceLength: 1000,
					Pieces:      []byte{1, 2, 3, 4, 5},
					Name:        "Test",
					Length:      5,
					Files: []File{
						{
							Length: 1,
							Path:   []string{"/home/"},
						},
					},
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewTorrentFromBencode(tt.args.data)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewTorrentFromBencode() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewTorrentFromBencode() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func assertTorrentMatchesExpected(t *testing.T, torrent *TorrentFile) {
	expectedAnnounce := "udp://tracker.openbittorrent.com:80/announce"
	expectedAnnounceList := [][]string{
		{"udp://tracker.openbittorrent.com:80/announce"},
		{"udp://tracker.opentrackr.org:1337/announce"},
	}
	expectedComment := "test"
	expectedCreatedBy := "uTorrent/3.5.5"
	expectedCreationDate := int64(1744127373)
	expectedName := "sub_zip.py"
	expectedLength := int64(829)
	expectedPieceLength := int64(16384)
	expectedPieceHashHex := "a660238a9b904f33ed71fe450c7e0befb1bf4451"

	if torrent.Announce != expectedAnnounce {
		t.Errorf("Announce mismatch. Got %q, expected %q", torrent.Announce, expectedAnnounce)
	}

	if len(torrent.AnnounceList) != len(expectedAnnounceList) {
		t.Errorf("AnnounceList length mismatch. Got %d, expected %d", len(torrent.AnnounceList), len(expectedAnnounceList))
	} else {
		for i, list := range torrent.AnnounceList {
			for j, url := range list {
				if url != expectedAnnounceList[i][j] {
					t.Errorf("AnnounceList[%d][%d] mismatch. Got %q, expected %q", i, j, url, expectedAnnounceList[i][j])
				}
			}
		}
	}

	if torrent.Comment != expectedComment {
		t.Errorf("Comment mismatch. Got %q, expected %q", torrent.Comment, expectedComment)
	}

	if torrent.CreatedBy != expectedCreatedBy {
		t.Errorf("CreatedBy mismatch. Got %q, expected %q", torrent.CreatedBy, expectedCreatedBy)
	}

	if torrent.CreationDate != expectedCreationDate {
		t.Errorf("CreationDate mismatch. Got %d, expected %d", torrent.CreationDate, expectedCreationDate)
	}

	info := torrent.Info

	if info.Name != expectedName {
		t.Errorf("Info.Name mismatch. Got %q, expected %q", info.Name, expectedName)
	}

	if info.Length != expectedLength {
		t.Errorf("Info.Length mismatch. Got %d, expected %d", info.Length, expectedLength)
	}

	if info.PieceLength != expectedPieceLength {
		t.Errorf("Info.PieceLength mismatch. Got %d, expected %d", info.PieceLength, expectedPieceLength)
	}

	actualPieceHashHex := fmt.Sprintf("%x", info.Pieces)
	if actualPieceHashHex != expectedPieceHashHex {
		t.Errorf("Info.Pieces hash mismatch. Got %s, expected %s", actualPieceHashHex, expectedPieceHashHex)
	}
}

func TestNewTorrentFromFile(t *testing.T) {
	path := "./testdata/sub_zip.py.torrent"

	torrent, err := NewTorrentFromFile(path)
	if err != nil {
		t.Fatalf("Failed to read torrent file: %v", err)
	}

	assertTorrentMatchesExpected(t, torrent)
}

func TestNewTorrentFromReader(t *testing.T) {
	data, err := os.ReadFile("./testdata/sub_zip.py.torrent")
	if err != nil {
		t.Fatalf("Failed to load test torrent file: %v", err)
	}

	reader := bytes.NewReader(data)

	torrent, err := NewTorrentFromReader(reader)
	if err != nil {
		t.Fatalf("Failed to parse torrent from reader: %v", err)
	}

	assertTorrentMatchesExpected(t, torrent)
}
