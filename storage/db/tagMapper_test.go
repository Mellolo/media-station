package db

import (
	"github.com/Mellolo/media-station/models/do/tagDO"
	"github.com/beego/beego/v2/client/orm"
	"reflect"
	"testing"
)

func init() {
	_ = orm.RegisterDataBase("default", "mysql", "media:media@tcp(localhost:3306)/media?charset=utf8mb4&parseTime=True&loc=Local")
}

func TestTagMapperImpl_DeleteArt(t *testing.T) {
	type args struct {
		artType string
		artId   int64
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
				artType: "video",
				artId:   1,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			impl := &TagMapperImpl{}
			if err := impl.DeleteArt(tt.args.artType, tt.args.artId, tt.args.tx...); (err != nil) != tt.wantErr {
				t.Errorf("DeleteArt() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestTagMapperImpl_DeleteTag(t *testing.T) {
	type args struct {
		tag string
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
				tag: "tag1",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			impl := &TagMapperImpl{}
			if err := impl.DeleteTag(tt.args.tag, tt.args.tx...); (err != nil) != tt.wantErr {
				t.Errorf("DeleteTag() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestTagMapperImpl_InsertOrUpdateTagsOfArt(t *testing.T) {
	type args struct {
		artType string
		artId   int64
		items   []tagDO.TagDO
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
				artType: "video",
				artId:   1,
				items: []tagDO.TagDO{
					{
						Tag: "tag1",
					},
					{
						Tag: "tag2",
					},
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			impl := &TagMapperImpl{}
			if err := impl.InsertOrUpdateTagsOfArt(tt.args.artType, tt.args.artId, tt.args.items, tt.args.tx...); (err != nil) != tt.wantErr {
				t.Errorf("InsertOrUpdateTagsOfArt() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestTagMapperImpl_SelectArtByTag(t *testing.T) {
	type args struct {
		artType string
		tag     string
		tx      []orm.TxOrmer
	}
	tests := []struct {
		name    string
		args    args
		want    []tagDO.TagDO
		wantErr bool
	}{
		{
			name: "case1",
			args: args{
				artType: "video",
				tag:     "tag1",
			},
			want: []tagDO.TagDO{
				{
					ArtType: "video",
					ArtId:   1,
					Tag:     "tag1",
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			impl := &TagMapperImpl{}
			got, err := impl.SelectArtByTag(tt.args.artType, tt.args.tag, tt.args.tx...)
			if (err != nil) != tt.wantErr {
				t.Errorf("SelectArtByTag() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("SelectArtByTag() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestTagMapperImpl_SelectTagByArt(t *testing.T) {
	type args struct {
		artType string
		artId   int64
		tx      []orm.TxOrmer
	}
	tests := []struct {
		name    string
		args    args
		want    []tagDO.TagDO
		wantErr bool
	}{
		{
			name: "case1",
			args: args{
				artType: "video",
				artId:   1,
			},
			want: []tagDO.TagDO{
				{
					ArtType: "video",
					ArtId:   1,
					Tag:     "tag1",
				},
				{
					ArtType: "video",
					ArtId:   1,
					Tag:     "tag2",
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			impl := &TagMapperImpl{}
			got, err := impl.SelectTagByArt(tt.args.artType, tt.args.artId, tt.args.tx...)
			if (err != nil) != tt.wantErr {
				t.Errorf("SelectTagByArt() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("SelectTagByArt() got = %v, want %v", got, tt.want)
			}
		})
	}
}
