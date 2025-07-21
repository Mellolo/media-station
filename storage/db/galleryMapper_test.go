package db

import (
	"github.com/beego/beego/v2/client/orm"
	_ "github.com/go-sql-driver/mysql"
	"media-station/models/do/galleryDO"
	"reflect"
	"testing"
)

func init() {
	_ = orm.RegisterDataBase("default", "mysql", "media:media@tcp(localhost:3306)/media?charset=utf8mb4&parseTime=True&loc=Local")
}

func TestGalleryMapperImpl_DeleteById(t *testing.T) {
	type args struct {
		id int64
		tx []orm.TxOrmer
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "case1",
			args: args{
				id: 1,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			impl := &GalleryMapperImpl{}
			if err := impl.DeleteById(tt.args.id, tt.args.tx...); (err != nil) != tt.wantErr {
				t.Errorf("DeleteById() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestGalleryMapperImpl_Insert(t *testing.T) {
	type args struct {
		gallery *galleryDO.GalleryDO
		tx      []orm.TxOrmer
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "case1",
			args: args{
				gallery: &galleryDO.GalleryDO{
					Name:            "test",
					Description:     "test",
					PageCount:       30,
					Actors:          nil,
					Tags:            nil,
					Uploader:        "",
					CoverUrl:        "",
					GalleryUrl:      "",
					PermissionLevel: "",
				},
				tx: nil,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			impl := &GalleryMapperImpl{}
			if _, err := impl.Insert(tt.args.gallery, tt.args.tx...); (err != nil) != tt.wantErr {
				t.Errorf("Insert() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestGalleryMapperImpl_SelectById(t *testing.T) {
	type args struct {
		id int64
		tx []orm.TxOrmer
	}
	tests := []struct {
		name    string
		args    args
		want    *galleryDO.GalleryDO
		wantErr bool
	}{
		{
			name: "case1",
			args: args{
				id: 1,
			},
			want: &galleryDO.GalleryDO{
				Id:              1,
				Name:            "1",
				Description:     "1",
				PageCount:       1,
				Actors:          []int{1},
				Tags:            []string{"1"},
				Uploader:        "1",
				CoverUrl:        "1",
				GalleryUrl:      "1",
				PermissionLevel: "1",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			impl := &GalleryMapperImpl{}
			got, err := impl.SelectById(tt.args.id, tt.args.tx...)
			if (err != nil) != tt.wantErr {
				t.Errorf("SelectById() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			got.CreateAt = orm.DateTimeField{}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("SelectById() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGalleryMapperImpl_Update(t *testing.T) {
	type args struct {
		id      int64
		gallery *galleryDO.GalleryDO
		tx      []orm.TxOrmer
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "case1",
			args: args{
				id: 2,
				gallery: &galleryDO.GalleryDO{
					Name:            "1",
					Description:     "1",
					PageCount:       1,
					Actors:          []int{1},
					Tags:            []string{"1"},
					Uploader:        "1",
					CoverUrl:        "1",
					GalleryUrl:      "1",
					PermissionLevel: "1",
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			impl := &GalleryMapperImpl{}
			if err := impl.Update(tt.args.id, tt.args.gallery, tt.args.tx...); (err != nil) != tt.wantErr {
				t.Errorf("Update() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
