package db

import (
	"github.com/beego/beego/v2/client/orm"
	_ "github.com/go-sql-driver/mysql"
	"media-station/models/do/tagDO"
	"reflect"
	"testing"
)

func init() {
	_ = orm.RegisterDataBase("default", "mysql", "media:media@tcp(localhost:3306)/media?charset=utf8mb4&parseTime=True&loc=Local")
}

func TestTagMapperImpl_DeleteById(t *testing.T) {
	type args struct {
		name string
		tx   []orm.TxOrmer
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "case1",
			args: args{
				name: "c",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			impl := &TagMapperImpl{}
			if err := impl.DeleteByName(tt.args.name, tt.args.tx...); (err != nil) != tt.wantErr {
				t.Errorf("DeleteById() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestTagMapperImpl_SelectById(t *testing.T) {
	type args struct {
		name string
		tx   []orm.TxOrmer
	}
	tests := []struct {
		name    string
		args    args
		want    *tagDO.TagDO
		wantErr bool
	}{
		{
			name: "case1",
			args: args{
				name: "c",
				tx:   nil,
			},
			want: &tagDO.TagDO{
				Id:      9,
				Name:    "c",
				Creator: "coc",
				Details: tagDO.TagDetails{
					VideoIds:   []int64{2},
					GalleryIds: []int64{2},
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			impl := &TagMapperImpl{}
			got, err := impl.SelectByName(tt.args.name, tt.args.tx...)
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

func TestTagMapperImpl_InsertOrUpdate(t *testing.T) {
	type args struct {
		id  int
		tag *tagDO.TagDO
		tx  []orm.TxOrmer
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "case1",
			args: args{
				tag: &tagDO.TagDO{
					Name:    "b",
					Creator: "bob",
					Details: tagDO.TagDetails{
						VideoIds:   []int64{2},
						GalleryIds: []int64{2},
					},
				},
			},
			wantErr: false,
		},
		{
			name: "case2",
			args: args{
				tag: &tagDO.TagDO{
					Name:    "c",
					Creator: "coc",
					Details: tagDO.TagDetails{
						VideoIds:   []int64{2},
						GalleryIds: []int64{2},
					},
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			impl := &TagMapperImpl{}
			if err := impl.InsertOrUpdate(tt.args.tag, tt.args.tx...); (err != nil) != tt.wantErr {
				t.Errorf("Update() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
