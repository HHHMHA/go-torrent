package torrent

import (
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
