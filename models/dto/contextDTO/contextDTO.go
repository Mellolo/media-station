package contextDTO

import "media-station/models/dto/userDTO"

// +k8s:deepcopy-gen=true
type ContextDTO struct {
	UserClaim userDTO.UserClaimDTO `json:"userClaim"`
}
