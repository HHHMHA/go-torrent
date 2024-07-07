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
		// List decoding test cases
		{
			name:    "Empty List Decode",
			args:    args{data: []byte("le")},
			want:    []interface{}{},
			wantErr: nil,
		},
		{
			name:    "List of Integers Decode",
			args:    args{data: []byte("li1ei2ei3ee")},
			want:    []interface{}{int64(1), int64(2), int64(3)},
			wantErr: nil,
		},
		{
			name:    "List of Strings Decode",
			args:    args{data: []byte("l4:spam4:eggse")},
			want:    []interface{}{[]byte("spam"), []byte("eggs")},
			wantErr: nil,
		},
		{
			name: "Nested List Decode",
			args: args{data: []byte("ll4:spam4:eggsei1e4:spami2ee")},
			want: []interface{}{
				[]interface{}{[]byte("spam"), []byte("eggs")},
				int64(1), []byte("spam"), int64(2),
			},
			wantErr: nil,
		},
		{
			name:    "List Decode with Invalid Format",
			args:    args{data: []byte("l4:spam4:eggsi2e")},
			want:    nil,
			wantErr: fmt.Errorf("list element format invalid"),
		},
		{
			name:    "List Decode with Invalid elements",
			args:    args{data: []byte("l4:spa4:eggsi2e")},
			want:    nil,
			wantErr: fmt.Errorf("list element format invalid"),
		},
		{
			name:    "List Decode with Invalid elements",
			args:    args{data: []byte("l4:spammmm4:spame")},
			want:    nil,
			wantErr: fmt.Errorf("list element format invalid"),
		},
		// TODO: list with dict inside
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			bencoder := NewSimpleBencoder()
			got, err := bencoder.Decode(tt.args.data)
			if err != nil && tt.wantErr != nil && err.Error() != tt.wantErr.Error() {
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
