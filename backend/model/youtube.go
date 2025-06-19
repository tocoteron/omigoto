package model

import (
	"context"
	"net/url"
	"time"
)

// Repository

type YouTubeRepository interface {
	GetChannel(
		ctx context.Context,
		channelID YouTubeChannelID,
	) (*YouTubeChannel, error)

	GetUploadsPlaylist(
		ctx context.Context,
		channelID YouTubeChannelID,
	) (*YouTubePlaylist, error)

	ListPlaylists(
		ctx context.Context,
		channelID YouTubeChannelID,
		pageToken *YouTubePageToken,
	) ([]*YouTubePlaylist, int64, *YouTubePageToken, error)

	ListVideoIDs(
		ctx context.Context,
		playlistID YouTubePlaylistID,
		pageToken *YouTubePageToken,
	) ([]YouTubeVideoID, int64, *YouTubePageToken, error)

	ListVideos(
		ctx context.Context,
		videoIDs []YouTubeVideoID,
		pageToken *YouTubePageToken,
	) ([]*YouTubeVideo, int64, *YouTubePageToken, error)
}

// IDs

type YouTubeChannelID string

type YouTubeChannelHandle string

type YouTubePlaylistID string

type YouTubeVideoID string

// Resources

type YouTubeChannel struct {
	ID                YouTubeChannelID
	Handle            YouTubeChannelHandle
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

// Utils

type YouTubePageToken string
