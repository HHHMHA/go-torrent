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
		// List with Dictionary inside
		{
			name:    "List with Dict Decode",
			args:    args{data: []byte("ld3:bar4:spam3:fooi42eee")},
			want:    []interface{}{map[string]interface{}{"bar": []byte("spam"), "foo": int64(42)}},
			wantErr: nil,
		},
		{
			name:    "List with Nested Dict Decode",
			args:    args{data: []byte("ld3:bar4:spam3:food3:bari42eeee")},
			want:    []interface{}{map[string]interface{}{"bar": []byte("spam"), "foo": map[string]interface{}{"bar": int64(42)}}},
			wantErr: nil,
		},
		{
			name:    "List with Dict Decode with Invalid Key Format",
			args:    args{data: []byte("ld3:bar4:spami42ee")},
			want:    nil,
			wantErr: fmt.Errorf("list element format invalid"),
		},
		{
			name:    "List with Dict Decode with Invalid Value Format",
			args:    args{data: []byte("ld3:bar4:spami42e")},
			want:    nil,
			wantErr: fmt.Errorf("list element format invalid"),
		},
		{
			name:    "List with Dict Decode with Unsorted Keys",
			args:    args{data: []byte("ld3:foo4:spam3:bar4:eggs3:foo4:testee")},
			want:    nil,
			wantErr: fmt.Errorf("list element format invalid"),
		},
		// Dictionary decoding test cases
		{
			name:    "Empty Dict Decode",
			args:    args{data: []byte("de")},
			want:    map[string]interface{}{},
			wantErr: nil,
		},
		{
			name:    "Dict with Integers and Strings Decode",
			args:    args{data: []byte("d3:bar4:spam3:fooi42ee")},
			want:    map[string]interface{}{"bar": []byte("spam"), "foo": int64(42)},
			wantErr: nil,
		},
		{
			name:    "Nested Dict Decode",
			args:    args{data: []byte("d3:bar4:spam3:food3:bari42eee")},
			want:    map[string]interface{}{"bar": []byte("spam"), "foo": map[string]interface{}{"bar": int64(42)}},
			wantErr: nil,
		},
		{
			name:    "Dict Decode with Invalid Key Format",
			args:    args{data: []byte("d3:bar4:spami42ee")},
			want:    nil,
			wantErr: fmt.Errorf("invalid dictionary format"),
		},
		{
			name:    "Dict Decode with Invalid Value Format",
			args:    args{data: []byte("d3:bar4:spami42")},
			want:    nil,
			wantErr: fmt.Errorf("invalid dictionary format"),
		},
		{
			name:    "Dict Decode with repeated Keys",
			args:    args{data: []byte("d3:foo4:spam3:bar4:eggs3:foo4:teste")},
			want:    nil,
			wantErr: fmt.Errorf("invalid dictionary format"),
		},
		{
			name:    "Dict with List Decode",
			args:    args{data: []byte("d3:barli1ei2ei3ee3:foo4:spam4:spaml4:eggsee")},
			want:    map[string]interface{}{"bar": []interface{}{int64(1), int64(2), int64(3)}, "foo": []byte("spam"), "spam": []interface{}{[]byte("eggs")}},
			wantErr: nil,
		},
		{
			name:    "Dict with List Decode Invalid End",
			args:    args{data: []byte("d3:barli1ei2ei3ee3:foo4:spam4:spaml4:eggse")},
			want:    nil,
			wantErr: fmt.Errorf("invalid dictionary format"),
		},
		{
			name:    "Dict with Nested List Decode",
			args:    args{data: []byte("d3:barlli1ei2eei3ee3:foo4:spam4:spaml4:eggsee")},
			want:    map[string]interface{}{"bar": []interface{}{[]interface{}{int64(1), int64(2)}, int64(3)}, "foo": []byte("spam"), "spam": []interface{}{[]byte("eggs")}},
			wantErr: nil,
		},
		{
			name:    "Dict with List Decode with Invalid Format",
			args:    args{data: []byte("d3:barli1ei2ei3e3:foo4:spam3:spaml4:eggse")},
			want:    nil,
			wantErr: fmt.Errorf("invalid dictionary format"),
		},
		{
			name:    "Dict with List Decode with Invalid Elements",
			args:    args{data: []byte("d3:barli1e4:eggs3:foo4:spam3:spaml4:eggse")},
			want:    nil,
			wantErr: fmt.Errorf("invalid dictionary format"),
		},
		{
			name: "Full Torrent File Decode",
			args: args{data: []byte("d8:announce14:http://tracker3:foo5:hello4:infod5:filesld6:lengthi12345e4:pathl8:filenameeee4:name8:testfile12:piece lengthi16384e6:pieces20:12345678901234567890ee")},
			want: map[string]interface{}{
				"announce": []byte("http://tracker"),
				"foo":      []byte("hello"),
				"info": map[string]interface{}{
					"files": []interface{}{
						map[string]interface{}{
							"length": int64(12345),
							"path":   []interface{}{[]byte("filename")},
						},
					},
					"name":         []byte("testfile"),
					"piece length": int64(16384),
					"pieces":       []byte("12345678901234567890"),
				},
			},
			wantErr: nil,
		},
		{
			name: "Full Torrent File Decode 2",
			args: args{data: []byte{
				100, 49, 51, 58, 99, 114, 101, 97, 116, 105, 111, 110, 32, 100, 97, 116, 101, 105, 49, 52, 53, 50, 52, 54, 56, 55, 50, 53, 48, 57, 49, 101, 56, 58, 101, 110, 99, 111, 100, 105, 110, 103, 53, 58, 85, 84, 70, 45, 56, 52, 58, 105, 110, 102, 111, 100, 54, 58, 108, 101, 110, 103, 116, 104, 105, 49, 54, 51, 55, 56, 51, 101, 52, 58, 110, 97, 109, 101, 57, 58, 97, 108, 105, 99, 101, 46, 116, 120, 116, 49, 50, 58, 112, 105, 101, 99, 101, 32, 108, 101, 110, 103, 116, 104, 105, 49, 54, 51, 56, 52, 101, 54, 58, 112, 105, 101, 99, 101, 115, 50, 48, 48, 58, 36, 192, 99, 82, 184, 241, 141, 203, 196, 131, 20, 34, 77, 108, 162, 38, 14, 24, 242, 191, 210, 203, 185, 139, 225, 63, 229, 126, 97, 253, 2, 36, 169, 2, 24, 60, 125, 94, 174, 101, 65, 191, 31, 23, 187, 228, 99, 219, 57, 27, 105, 129, 220, 175, 47, 249, 67, 66, 88, 219, 90, 69, 8, 190, 16, 91, 237, 212, 48, 81, 204, 248, 77, 212, 226, 202, 22, 118, 93, 234, 188, 70, 204, 161, 101, 0, 254, 14, 115, 49, 160, 146, 35, 157, 73, 49, 209, 157, 253, 65, 108, 71, 131, 71, 193, 148, 236, 27, 225, 45, 208, 104, 88, 119, 25, 194, 42, 248, 106, 155, 141, 75, 83, 107, 165, 237, 219, 6, 68, 242, 128, 7, 123, 191, 216, 99, 203, 123, 152, 96, 234, 210, 60, 79, 60, 124, 15, 71, 156, 53, 40, 2, 159, 159, 239, 184, 150, 117, 135, 129, 171, 163, 218, 137, 252, 11, 185, 71, 71, 168, 84, 170, 129, 181, 158, 238, 69, 34, 2, 103, 217, 14, 2, 89, 218, 191, 146, 13, 129, 88, 40, 232, 215, 93, 177, 130, 205, 43, 248, 100, 101, 101,
			}},
			/// {
			//   "creation date": 1452468725091,
			//   "encoding": "UTF-8",
			//   "info": {
			//      "length": 163783,
			//      "name": "alice.txt",
			//      "piece length": 16384,
			//      "pieces": "<hex>24 C0 63 52 B8 F1 8D CB C4 83 14 22 4D 6C A2 26 0E 18 F2 BF D2 CB B9 8B E1 3F E5 7E 61 FD 02 24 A9 02 18 3C 7D 5E AE 65 41 BF 1F 17 BB E4 63 DB 39 1B 69 81 DC AF 2F F9 43 42 58 DB 5A 45 08 BE 10 5B ED D4 30 51 CC F8 4D D4 E2 CA 16 76 5D EA BC 46 CC A1 65 00 FE 0E 73 31 A0 92 23 9D 49 31 D1 9D FD 41 6C 47 83 47 C1 94 EC 1B E1 2D D0 68 58 77 19 C2 2A F8 6A 9B 8D 4B 53 6B A5 ED DB 06 44 F2 80 07 7B BF D8 63 CB 7B 98 60 EA D2 3C 4F 3C 7C 0F 47 9C 35 28 02 9F 9F EF B8 96 75 87 81 AB A3 DA 89 FC 0B B9 47 47 A8 54 AA 81 B5 9E EE 45 22 02 67 D9 0E 02 59 DA BF 92 0D 81 58 28 E8 D7 5D B1 82 CD 2B F8 64</hex>"
			//   }
			//}
			want: map[string]interface{}{
				"creation date": int64(1452468725091),
				"encoding":      []byte("UTF-8"),
				"info": map[string]interface{}{
					"name":         []byte("alice.txt"),
					"length":       int64(163783),
					"piece length": int64(16384),
					"pieces":       []byte{36, 192, 99, 82, 184, 241, 141, 203, 196, 131, 20, 34, 77, 108, 162, 38, 14, 24, 242, 191, 210, 203, 185, 139, 225, 63, 229, 126, 97, 253, 2, 36, 169, 2, 24, 60, 125, 94, 174, 101, 65, 191, 31, 23, 187, 228, 99, 219, 57, 27, 105, 129, 220, 175, 47, 249, 67, 66, 88, 219, 90, 69, 8, 190, 16, 91, 237, 212, 48, 81, 204, 248, 77, 212, 226, 202, 22, 118, 93, 234, 188, 70, 204, 161, 101, 0, 254, 14, 115, 49, 160, 146, 35, 157, 73, 49, 209, 157, 253, 65, 108, 71, 131, 71, 193, 148, 236, 27, 225, 45, 208, 104, 88, 119, 25, 194, 42, 248, 106, 155, 141, 75, 83, 107, 165, 237, 219, 6, 68, 242, 128, 7, 123, 191, 216, 99, 203, 123, 152, 96, 234, 210, 60, 79, 60, 124, 15, 71, 156, 53, 40, 2, 159, 159, 239, 184, 150, 117, 135, 129, 171, 163, 218, 137, 252, 11, 185, 71, 71, 168, 84, 170, 129, 181, 158, 238, 69, 34, 2, 103, 217, 14, 2, 89, 218, 191, 146, 13, 129, 88, 40, 232, 215, 93, 177, 130, 205, 43, 248, 100},
				},
			},
			wantErr: nil,
		},
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
		{
			name:    "Empty Encode",
			args:    args{data: nil},
			want:    nil,
			wantErr: true,
		},
		{
			name:    "Integer Encode",
			args:    args{data: int64(128)},
			want:    []byte("i128e"),
			wantErr: false,
		},
		{
			name:    "String Encode",
			args:    args{data: "spam"},
			want:    []byte("4:spam"),
			wantErr: false,
		},
		{
			name:    "String Encode Array",
			args:    args{data: []byte("spam")},
			want:    []byte("4:spam"),
			wantErr: false,
		},
		{
			name:    "Empty String",
			args:    args{data: []byte("")},
			want:    []byte("0:"),
			wantErr: false,
		},
		{
			name:    "Empty List Encode",
			args:    args{data: []interface{}{}},
			want:    []byte("le"),
			wantErr: false,
		},
		{
			name:    "List of Integers Encode",
			args:    args{data: []interface{}{int64(1), int64(2), int64(3)}},
			want:    []byte("li1ei2ei3ee"),
			wantErr: false,
		},
		{
			name:    "List of Strings Encode",
			args:    args{data: []interface{}{[]byte("spam"), []byte("eggs")}},
			want:    []byte("l4:spam4:eggse"),
			wantErr: false,
		},
		{
			name: "Nested List Encode",
			args: args{data: []interface{}{
				[]interface{}{[]byte("spam"), []byte("eggs")},
				int64(1), []byte("spam"), int64(2),
			}},
			want:    []byte("ll4:spam4:eggsei1e4:spami2ee"),
			wantErr: false,
		},
		{
			name:    "List with Dict Encode",
			args:    args{data: []interface{}{map[string]interface{}{"bar": []byte("spam"), "foo": int64(42)}}},
			want:    []byte("ld3:bar4:spam3:fooi42eee"),
			wantErr: false,
		},
		{
			name:    "Empty Dict Encode",
			args:    args{data: map[string]interface{}{}},
			want:    []byte("de"),
			wantErr: false,
		},
		{
			name:    "Dict Encode",
			args:    args{data: map[string]interface{}{"key": "value"}},
			want:    []byte("d3:key5:valuee"),
			wantErr: false,
		},
		{
			name:    "Nested Dict Encode",
			args:    args{data: map[string]interface{}{"bar": []byte("spam"), "foo": map[string]interface{}{"bar": int64(42)}}},
			want:    []byte("d3:bar4:spam3:food3:bari42eee"),
			wantErr: false,
		},
		{
			name:    "Dict with List Encode",
			args:    args{data: map[string]interface{}{"bar": []interface{}{int64(1), int64(2), int64(3)}, "foo": []byte("spam"), "spam": []interface{}{[]byte("eggs")}}},
			want:    []byte("d3:barli1ei2ei3ee3:foo4:spam4:spaml4:eggsee"),
			wantErr: false,
		},
		{
			name: "Full Torrent File Encode",
			args: args{data: map[string]interface{}{
				"announce": []byte("http://tracker"),
				"foo":      []byte("hello"),
				"info": map[string]interface{}{
					"files": []interface{}{
						map[string]interface{}{
							"length": int64(12345),
							"path":   []interface{}{[]byte("filename")},
						},
					},
					"name":         []byte("testfile"),
					"piece length": int64(16384),
					"pieces":       []byte("12345678901234567890"),
				},
			}},
			want:    []byte("d8:announce14:http://tracker3:foo5:hello4:infod5:filesld6:lengthi12345e4:pathl8:filenameeee4:name8:testfile12:piece lengthi16384e6:pieces20:12345678901234567890ee"),
			wantErr: false,
		},
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

func TestSimpleBencoder_Unmarshal(t *testing.T) {
	type args struct {
		data   []byte
		target interface{}
	}

	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Unmarshal int",
			args: args{
				data: []byte("d3:fooi123ee"), // Bencode for {"foo": 123}
				target: new(struct {
					Foo int64 `bencode:"foo"`
				}),
			},
			wantErr: false,
		},
		{
			name: "Unmarshal string",
			args: args{
				data: []byte("d3:foo4:spamee"), // Bencode for {"foo": "spam"}
				target: new(struct {
					Foo string `bencode:"foo"`
				}),
			},
			wantErr: false,
		},
		{
			name: "Unmarshal slice of strings",
			args: args{
				data: []byte("d6:valuesl4:spam4:eggs4:testee"), // Bencode for {"values": ["spam", "eggs", "test"]}
				target: new(struct {
					Values []string `bencode:"values"`
				}),
			},
			wantErr: false,
		},
		{
			name: "Unmarshal slice of slices",
			args: args{
				data: []byte("d6:valuesll4:spam4:eggseee"), // Bencode for {"values": [["spam", "eggs"], []]}
				target: new(struct {
					Values [][]string `bencode:"values"`
				}),
			},
			wantErr: false,
		},
		{
			name: "Unmarshal map of slices",
			args: args{
				data: []byte("d5:filesld4:name8:filenameed4:name7:filetwoeee"), // Bencode for {"files": [{"name": "filename"}, {"name": "filetwo"}]}
				target: new(struct {
					Files []struct {
						Name string `bencode:"name"`
					} `bencode:"files"`
				}),
			},
			wantErr: false,
		},
		{
			name: "Unmarshal nested structs",
			args: args{
				data: []byte("d4:userd4:name4:John3:agei30eee"), // Bencode for {"user": {"name": "John", "age": 30}}
				target: new(struct {
					User struct {
						Name string `bencode:"name"`
						Age  int64  `bencode:"age"`
					} `bencode:"user"`
				}),
			},
			wantErr: false,
		},
		{
			name: "Unmarshal invalid data",
			args: args{
				data: []byte("d3:foo4:bari3:100ee"), // Incorrect bencode
				target: new(struct {
					Foo string `bencode:"foo"`
				}),
			},
			wantErr: true, // Expecting an error due to incorrect bencode
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			bencoder := &SimpleBencoder{}
			if err := bencoder.Unmarshal(tt.args.data, tt.args.target); (err != nil) != tt.wantErr {
				t.Errorf("Unmarshal() error = %v, wantErr %v", err, tt.wantErr)
			}
			// Optionally, check the values of the target struct after unmarshalling
			if !tt.wantErr && !reflect.DeepEqual(tt.args.target, tt.args.target) {
				t.Errorf("Unmarshal() got = %v, want = %v", tt.args.target, tt.args.target)
			}
		})
	}
}

func TestSimpleBencoder_Marshal(t *testing.T) {
	type args struct {
		target interface{}
	}
	tests := []struct {
		name    string
		args    args
		want    []byte
		wantErr bool
	}{
		{
			name: "Marshal int",
			args: args{
				target: struct {
					Foo int64 `bencode:"foo"`
				}{
					Foo: 123,
				},
			},
			want:    []byte("d3:fooi123ee"),
			wantErr: false,
		},
		{
			name: "Marshal string",
			args: args{
				target: struct {
					Foo string `bencode:"foo"`
				}{
					Foo: "spam",
				},
			},
			want:    []byte("d3:foo4:spame"),
			wantErr: false,
		},
		{
			name: "Marshal slice of strings",
			args: args{
				target: struct {
					Values []string `bencode:"values"`
				}{
					Values: []string{"spam", "eggs", "test"},
				},
			},
			want:    []byte("d6:valuesl4:spam4:eggs4:testee"),
			wantErr: false,
		},
		{
			name: "Marshal slice of slices",
			args: args{
				target: struct {
					Values [][]string `bencode:"values"`
				}{
					Values: [][]string{{"spam", "eggs"}, {}},
				},
			},
			want:    []byte("d6:valuesll4:spam4:eggseleee"),
			wantErr: false,
		},
		{
			name: "Marshal map of slices",
			args: args{
				target: struct {
					Files []struct {
						Name string `bencode:"name"`
					} `bencode:"files"`
				}{
					Files: []struct {
						Name string `bencode:"name"`
					}{
						{Name: "filename"},
						{Name: "filetwo"},
					},
				},
			},
			want:    []byte("d5:filesld4:name8:filenameed4:name7:filetwoeee"),
			wantErr: false,
		},
		{
			name: "Marshal nested structs",
			args: args{
				target: struct {
					User struct {
						Name string `bencode:"name"`
						Age  int64  `bencode:"age"`
					} `bencode:"user"`
				}{
					User: struct {
						Name string `bencode:"name"`
						Age  int64  `bencode:"age"`
					}{
						Name: "John",
						Age:  30,
					},
				},
			},
			want:    []byte("d4:userd3:agei30e4:name4:Johnee"),
			wantErr: false,
		},
		{
			name: "Marshal invalid target",
			args: args{
				target: 123,
			},
			want:    nil,
			wantErr: true, // Expecting an error due to invalid target
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			bencoder := &SimpleBencoder{}
			got, err := bencoder.Marshal(tt.args.target)
			if (err != nil) != tt.wantErr {
				t.Errorf("Marshal() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Marshal() got = %v, want %v", got, tt.want)
				t.Errorf("Marshal() got = %v, want %v", string(got), string(tt.want))
			}
		})
	}
}
