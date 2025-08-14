package db

import (
	"github.com/beego/beego/v2/client/orm"
	"media-station/models/do/performDO"
	"reflect"
	"testing"
)

func init() {
	_ = orm.RegisterDataBase("default", "mysql", "media:media@tcp(localhost:3306)/media?charset=utf8mb4&parseTime=True&loc=Local")
}

func TestPerformMapperImpl_DeleteActor(t *testing.T) {
	type args struct {
		actorId int64
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
				actorId: 1,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			impl := PerformMapperImpl{}
			if err := impl.DeleteActor(tt.args.actorId, tt.args.tx...); (err != nil) != tt.wantErr {
				t.Errorf("DeleteActor() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestPerformMapperImpl_DeleteArt(t *testing.T) {
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
			impl := PerformMapperImpl{}
			if err := impl.DeleteArt(tt.args.artType, tt.args.artId, tt.args.tx...); (err != nil) != tt.wantErr {
				t.Errorf("DeleteArt() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestPerformMapperImpl_InsertOrUpdateActorsOfArt(t *testing.T) {
	type args struct {
		artType string
		artId   int64
		items   []performDO.PerformDO
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
				items: []performDO.PerformDO{
					{
						ActorId: 1,
					},
					{
						ActorId: 2,
					},
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			impl := PerformMapperImpl{}
			if err := impl.InsertOrUpdateActorsOfArt(tt.args.artType, tt.args.artId, tt.args.items, tt.args.tx...); (err != nil) != tt.wantErr {
				t.Errorf("InsertOrUpdateActorsOfArt() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestPerformMapperImpl_SelectActorByArt(t *testing.T) {
	type args struct {
		artType string
		artId   int64
		tx      []orm.TxOrmer
	}
	tests := []struct {
		name    string
		args    args
		want    []performDO.PerformDO
		wantErr bool
	}{
		{
			name: "case1",
			args: args{
				artType: "video",
				artId:   1,
			},
			want: []performDO.PerformDO{
				{
					ArtType: "video",
					ActorId: 1,
					ArtId:   1,
				},
				{
					ArtType: "video",
					ActorId: 2,
					ArtId:   1,
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			impl := PerformMapperImpl{}
			got, err := impl.SelectActorByArt(tt.args.artType, tt.args.artId, tt.args.tx...)
			if (err != nil) != tt.wantErr {
				t.Errorf("SelectActorByArt() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("SelectActorByArt() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPerformMapperImpl_SelectArtByActor(t *testing.T) {
	type args struct {
		artType string
		actorId int64
		tx      []orm.TxOrmer
	}
	tests := []struct {
		name    string
		args    args
		want    []performDO.PerformDO
		wantErr bool
	}{
		{
			name: "case1",
			args: args{
				artType: "video",
				actorId: 1,
			},
			want: []performDO.PerformDO{
				{
					ArtType: "video",
					ActorId: 1,
					ArtId:   1,
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			impl := PerformMapperImpl{}
			got, err := impl.SelectArtByActor(tt.args.artType, tt.args.actorId, tt.args.tx...)
			if (err != nil) != tt.wantErr {
				t.Errorf("SelectArtByActor() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("SelectArtByActor() got = %v, want %v", got, tt.want)
			}
		})
	}
}
