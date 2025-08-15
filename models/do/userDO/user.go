package userDO

// +k8s:deepcopy-gen=true
type UserDO struct {
	Id       int64
	Username string
	Password string
}
