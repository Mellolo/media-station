package db

import (
	"github.com/Mellolo/media-station/enum"
	"github.com/Mellolo/media-station/models/do/videoDO"
	"github.com/beego/beego/v2/client/orm"
	_ "github.com/go-sql-driver/mysql"
	"reflect"
	"testing"
)

func init() {
	_ = orm.RegisterDataBase("default", "mysql", "media:media@tcp(localhost:3306)/media?charset=utf8mb4&parseTime=True&loc=Local")
}

func TestVideoMapperImpl_DeleteById(t *testing.T) {
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
				id: 5,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			impl := NewVideoMapper()
			if err := impl.DeleteById(tt.args.id, tt.args.tx...); (err != nil) != tt.wantErr {
				t.Errorf("DeleteById() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestVideoMapperImpl_Insert(t *testing.T) {
	type args struct {
		video videoDO.VideoDO
		tx    []orm.TxOrmer
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "case1",
			args: args{
				video: videoDO.VideoDO{
					Name:            "test",
					Description:     "test",
					Uploader:        "test",
					VideoUrl:        "test",
					CoverUrl:        "test",
					PermissionLevel: enum.PermissionPublic,
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			impl := &VideoMapperImpl{}
			if _, err := impl.Insert(tt.args.video, tt.args.tx...); (err != nil) != tt.wantErr {
				t.Errorf("Insert() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestVideoMapperImpl_Update(t *testing.T) {
	type args struct {
		id    int64
		video videoDO.VideoDO
		tx    []orm.TxOrmer
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "case1",
			args: args{
				id: 5,
				video: videoDO.VideoDO{
					Id:              5,
					Name:            "test",
					Description:     "test",
					Uploader:        "test1",
					VideoUrl:        "test1",
					CoverUrl:        "test1",
					PermissionLevel: enum.PermissionPublic,
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			impl := NewVideoMapper()
			if err := impl.Update(tt.args.id, tt.args.video, tt.args.tx...); (err != nil) != tt.wantErr {
				t.Errorf("Update() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestVideoMapperImpl_SelectById(t *testing.T) {
	type args struct {
		id int64
		tx []orm.TxOrmer
	}
	tests := []struct {
		name    string
		args    args
		want    videoDO.VideoDO
		wantErr bool
	}{
		{
			name: "case1",
			args: args{
				id: 4,
			},
			want: videoDO.VideoDO{
				Id:              4,
				Name:            "test",
				Description:     "test",
				Uploader:        "test1",
				VideoUrl:        "test1",
				CoverUrl:        "test1",
				PermissionLevel: enum.PermissionPublic,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			impl := &VideoMapperImpl{}
			got, err := impl.SelectById(tt.args.id, tt.args.tx...)
			if (err != nil) != tt.wantErr {
				t.Errorf("SelectById() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			got.CreateAt = ""
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("SelectById() got = %v, want %v", got, tt.want)
			}
		})
	}
}
