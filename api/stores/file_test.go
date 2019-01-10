package stores

import (
	"couchsport/api/types"
	"github.com/golang/leveldb/db"
	"github.com/golang/leveldb/memfs"
	"io"
	"os"
	"strings"
	"testing"
)

type memFS struct {
	_os db.FileSystem
}

func (mem memFS) OpenFile(name string) (io.WriteCloser, error) {
	f, err := mem._os.Create(name)
	return f, err
}

func (mem memFS) Stat(name string) (os.FileInfo, error) {
	return mem._os.Stat(name)
}

func (mem memFS) MkdirAll(path string) error {
	return mem._os.MkdirAll(path, 0700)
}

func (mem memFS) IsNotExist(err error) bool {
	if err != nil {
		return true
	}
	return false
}

const DefaultPublicPath = "public/"

func TestFileStore_Save(t *testing.T) {
	type fields struct {
		FileSystem    types.FileSystem
		PublicPath    string
		ImageBasePath string
		FilePrefix    string
	}
	type args struct {
		UserID   uint
		prefix   string
		filename string
		buf      io.Reader
	}

	memos := memFS{_os: memfs.New()}

	tests := []struct {
		name    string
		fields  fields
		args    args
		want    string
		wantErr bool
	}{
		{
			name: "should return an error",
			fields: fields{
				FileSystem:    memos,
				PublicPath:    DefaultPublicPath,
				ImageBasePath: "static/img/",
				FilePrefix:    "isupload.",
			},
			args: args{
				UserID:   35,
				prefix:   "user-",
				filename: "test-file.jpg",
				buf:      strings.NewReader(``),
			},
			want:    "",
			wantErr: true,
		},
		{
			name: "should return correct filename",
			fields: fields{
				FileSystem:    memos,
				PublicPath:    DefaultPublicPath,
				ImageBasePath: "static/img/",
				FilePrefix:    "isupload.",
			},
			args: args{
				UserID:   3,
				prefix:   "user-",
				filename: "test-file-1.jpg",
				buf:      strings.NewReader(`tototototo`),
			},
			want:    "/static/img/user-3/isupload.test-file-1.jpg",
			wantErr: false,
		},
		{
			name: "empty prefix",
			fields: fields{
				FileSystem:    memos,
				PublicPath:    DefaultPublicPath,
				ImageBasePath: "static/img/",
				FilePrefix:    "isupload.",
			},
			args: args{
				UserID:   3,
				prefix:   "",
				filename: "test-file-1.jpg",
				buf:      strings.NewReader(`tototototo`),
			},
			want:    "/static/img/3/isupload.test-file-1.jpg",
			wantErr: false,
		},
		{
			name: "empty file prefix",
			fields: fields{
				FileSystem:    memos,
				PublicPath:    DefaultPublicPath,
				ImageBasePath: "static/img/",
				FilePrefix:    "",
			},
			args: args{
				UserID:   3,
				prefix:   "user-",
				filename: "test-file-1.jpg",
				buf:      strings.NewReader(`tototototo`),
			},
			want:    "/static/img/user-3/test-file-1.jpg",
			wantErr: false,
		},
		{
			name: "empty filename",
			fields: fields{
				FileSystem:    memos,
				PublicPath:    DefaultPublicPath,
				ImageBasePath: "static/img/",
				FilePrefix:    "isupload.",
			},
			args: args{
				UserID:   3,
				prefix:   "user-",
				filename: "",
				buf:      strings.NewReader(`tototototo`),
			},
			want:    "",
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			app := FileStore{
				FileSystem:    tt.fields.FileSystem,
				PublicPath:    tt.fields.PublicPath,
				ImageBasePath: tt.fields.ImageBasePath,
				FilePrefix:    tt.fields.FilePrefix,
			}
			got, err := app.Save(tt.args.UserID, tt.args.prefix, tt.args.filename, tt.args.buf)
			if (err != nil) != tt.wantErr {
				t.Errorf("FileStore.Save() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if got != tt.want {
				t.Errorf("FileStore.Save() = have %v, want %v", got, tt.want)
				return
			}

			if tt.wantErr {
				return
			}

			tmp, err := memos.Stat(DefaultPublicPath + got)
			if err != nil {
				t.Errorf("couldn stat file %s", got)
				return
			}
			t.Logf("file correctly stat %s", tmp)
		})
	}
}
