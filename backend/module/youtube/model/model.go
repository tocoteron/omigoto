package model

import (
	"net/url"
	"time"
)

type YouTubeChannelIdentity struct {
	ID     YouTubeChannelID
	Handle YouTubeChannelHandle
}

type YouTubeChannel struct {
	YouTubeChannelIdentity

	UploadsPlaylistID YouTubePlaylistID
}

type YouTubePlaylist struct {
	ID        YouTubePlaylistID
	IsUploads bool
	Title     *string // nil if IsUploads is true
}

type YouTubeVideo struct {
	ID                   YouTubeVideoID
	Title                string
	Description          string
	Duration             time.Duration
	Thumbnails           YouTubeVideoThumbnails
	LiveStreamingDetails *YouTubeVideoLiveStreamingDetails // nil if not live streaming
	PublishedAt          time.Time
}

type YouTubeVideoThumbnails struct {
	Default  *url.URL
	Medium   *url.URL
	High     *url.URL
	Standard *url.URL
	Maxres   *url.URL
}

type YouTubeVideoLiveStreamingDetails struct {
	ActualStartTime time.Time
	ActualEndTime   time.Time
	ScheduledStart  time.Time
}
