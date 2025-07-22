package db

import (
	"github.com/beego/beego/v2/client/orm"
	_ "github.com/go-sql-driver/mysql"
	"media-station/models/do/actorDO"
	"reflect"
	"testing"
)

func init() {
	_ = orm.RegisterDataBase("default", "mysql", "media:media@tcp(localhost:3306)/media?charset=utf8mb4&parseTime=True&loc=Local")
}

func TestActorMapperImpl_DeleteById(t *testing.T) {
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
			impl := &ActorMapperImpl{}
			if err := impl.DeleteById(tt.args.id, tt.args.tx...); (err != nil) != tt.wantErr {
				t.Errorf("DeleteById() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestActorMapperImpl_Insert(t *testing.T) {
	type args struct {
		actor *actorDO.ActorDO
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
				actor: &actorDO.ActorDO{
					Name:        "a",
					Description: "a",
					Creator:     "a",
					CoverUrl:    "a",
					Details: actorDO.ActorDetailsDO{
						VideoIds:   nil,
						GalleryIds: []int64{1},
					},
				},
				tx: nil,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			impl := &ActorMapperImpl{}
			id, err := impl.Insert(tt.args.actor, tt.args.tx...)
			if (err != nil) != tt.wantErr {
				t.Errorf("Insert() error = %v, wantErr %v", err, tt.wantErr)
			}
			t.Logf("id: %v", id)
		})
	}
}

func TestActorMapperImpl_SelectById(t *testing.T) {
	type args struct {
		id int64
		tx []orm.TxOrmer
	}
	tests := []struct {
		name    string
		args    args
		want    *actorDO.ActorDO
		wantErr bool
	}{
		{
			name: "case1",
			args: args{
				id: 1,
				tx: nil,
			},
			want: &actorDO.ActorDO{
				Id:          1,
				Name:        "a",
				Description: "a",
				Creator:     "a",
				CoverUrl:    "a",
				Details: actorDO.ActorDetailsDO{
					VideoIds:   nil,
					GalleryIds: []int64{1},
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			impl := &ActorMapperImpl{}
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

func TestActorMapperImpl_Update(t *testing.T) {
	type args struct {
		id    int64
		actor *actorDO.ActorDO
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
				id: 1,
				actor: &actorDO.ActorDO{
					Name:        "b",
					Description: "b",
					Creator:     "b",
					CoverUrl:    "b",
					Details: actorDO.ActorDetailsDO{
						VideoIds:   nil,
						GalleryIds: nil,
					},
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			impl := &ActorMapperImpl{}
			if err := impl.Update(tt.args.id, tt.args.actor, tt.args.tx...); (err != nil) != tt.wantErr {
				t.Errorf("Update() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
