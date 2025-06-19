package omikun

import "github.com/tocoteron/omigoto/backend/model"

type YouTubeChannelIdentity struct {
	ID     model.YouTubeChannelID
	Handle model.YouTubeChannelHandle
}

var YouTubeChannel = YouTubeChannelIdentity{
	ID:     model.YouTubeChannelID("UC1cnByKe24JjTv38tH_7BYw"),
	Handle: model.YouTubeChannelHandle("@izuho_omi"),
}
