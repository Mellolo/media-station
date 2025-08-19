package actorVO

import (
	"github.com/Mellolo/media-station/models/vo/galleryVO"
	"github.com/Mellolo/media-station/models/vo/videoVO"
)

// +k8s:deepcopy-gen=true
type ActorPageVO struct {
	Id          int64  `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Creator     string `json:"creator"`

	Videos    []videoVO.VideoItemVO     `json:"videos"`
	Galleries []galleryVO.GalleryItemVO `json:"galleries"`
}
