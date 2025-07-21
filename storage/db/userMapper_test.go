package db

import (
	"github.com/beego/beego/v2/client/orm"
	_ "github.com/go-sql-driver/mysql"
	"media-station/models/do/userDO"
	"reflect"
	"testing"
)

func init() {
	_ = orm.RegisterDataBase("default", "mysql", "media:media@tcp(localhost:3306)/media?charset=utf8mb4&parseTime=True&loc=Local")
}

func TestUserMapperImpl_DeleteById(t *testing.T) {
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
			impl := &UserMapperImpl{}
			if err := impl.DeleteById(tt.args.id, tt.args.tx...); (err != nil) != tt.wantErr {
				t.Errorf("DeleteById() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestUserMapperImpl_Insert(t *testing.T) {
	type args struct {
		user *userDO.UserDO
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
				user: &userDO.UserDO{
					Username:    "test_user",
					Password:    "password123",
					PhoneNumber: "1234567890",
					WechatId:    "wechat123",
					Details: userDO.UserDetails{
						VideoIds:   []int{1},
						GalleryIds: nil,
					},
				},
				tx: nil,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			impl := &UserMapperImpl{}
			if _, err := impl.Insert(tt.args.user, tt.args.tx...); (err != nil) != tt.wantErr {
				t.Errorf("Insert() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestUserMapperImpl_SelectById(t *testing.T) {
	type args struct {
		id int64
		tx []orm.TxOrmer
	}
	tests := []struct {
		name    string
		args    args
		want    *userDO.UserDO
		wantErr bool
	}{
		{
			name: "case1",
			args: args{
				id: 1,
				tx: nil,
			},
			want: &userDO.UserDO{
				Id:          1,
				Username:    "test_user",
				Password:    "password123",
				PhoneNumber: "1234567890",
				WechatId:    "wechat123",
				Details: userDO.UserDetails{
					VideoIds:   []int{1},
					GalleryIds: nil,
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			impl := &UserMapperImpl{}
			got, err := impl.SelectById(tt.args.id, tt.args.tx...)
			if (err != nil) != tt.wantErr {
				t.Errorf("SelectById() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("SelectById() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestUserMapperImpl_Update(t *testing.T) {
	type args struct {
		id   int64
		user *userDO.UserDO
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
				id: 1,
				user: &userDO.UserDO{
					Username:    "updated_user",
					Password:    "new_password",
					PhoneNumber: "0987654321",
					WechatId:    "wechat456",
					Details: userDO.UserDetails{
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
			impl := &UserMapperImpl{}
			if err := impl.Update(tt.args.id, tt.args.user, tt.args.tx...); (err != nil) != tt.wantErr {
				t.Errorf("Update() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
