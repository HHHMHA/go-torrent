package bencoder

import (
	"errors"
	"fmt"
	"reflect"
	"testing"
)

func TestSimpleBencoder_Decode(t *testing.T) {
	type args struct {
		data []byte
	}
	tests := []struct {
		name    string
		args    args
		want    interface{}
		wantErr error
	}{
		{
			name:    "Empty Decode",
			args:    args{data: []byte("")},
			want:    nil,
			wantErr: errors.New("empty data"),
		},
		{
			name:    "Integer Decode",
			args:    args{data: []byte("i128e")},
			want:    int64(128),
			wantErr: nil,
		},
		{
			name:    "Integer Decode Error",
			args:    args{data: []byte("i128")},
			want:    nil,
			wantErr: fmt.Errorf("invalid integer format: %s", string([]byte("i128"))),
		},
		{
			name:    "String Decode",
			args:    args{data: []byte("4:spam")},
			want:    []byte("spam"),
			wantErr: nil,
		},
		{
			name:    "String Decode Fail On Longer",
			args:    args{data: []byte("3:spam")},
			want:    nil,
			wantErr: fmt.Errorf("mismatch of length and byte string"),
		},
		{
			name:    "String Decode Fail On Shorter",
			args:    args{data: []byte("5:spam")},
			want:    nil,
			wantErr: fmt.Errorf("mismatch of length and byte string"),
		},
		{
			name:    "String Decode no length",
			args:    args{data: []byte("spam")},
			want:    nil,
			wantErr: fmt.Errorf("unknown format"),
		},
		{
			name:    "String Decode Fake Length",
			args:    args{data: []byte("s:spam")},
			want:    nil,
			wantErr: fmt.Errorf("length of string is not correct"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			bencoder := NewSimpleBencoder()
			got, err := bencoder.Decode(tt.args.data)
			if err != nil && err.Error() != tt.wantErr.Error() {
				t.Errorf("Decode() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Decode() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSimpleBencoder_Encode(t *testing.T) {
	type args struct {
		data interface{}
	}
	tests := []struct {
		name    string
		args    args
		want    []byte
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			bencoder := NewSimpleBencoder()
			got, err := bencoder.Encode(tt.args.data)
			if (err != nil) != tt.wantErr {
				t.Errorf("Encode() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Encode() got = %v, want %v", got, tt.want)
			}
		})
	}
}
