package repository

import (
	"context"

	"github.com/tocoteron/omigoto/backend/module/youtube/model"
)

type YouTubeRepository interface {
	GetChannel(
		ctx context.Context,
		channelID model.YouTubeChannelID,
	) (*model.YouTubeChannel, error)

	GetUploadsPlaylist(
		ctx context.Context,
		channelID model.YouTubeChannelID,
	) (*model.YouTubePlaylist, error)

	ListPlaylists(
		ctx context.Context,
		channelID model.YouTubeChannelID,
		pageToken *YouTubePageToken,
	) ([]*model.YouTubePlaylist, int64, *YouTubePageToken, error)

	ListVideoIDs(
		ctx context.Context,
		playlistID model.YouTubePlaylistID,
		pageToken *YouTubePageToken,
	) ([]model.YouTubeVideoID, int64, *YouTubePageToken, error)

	ListVideos(
		ctx context.Context,
		videoIDs []model.YouTubeVideoID,
		pageToken *YouTubePageToken,
	) ([]*model.YouTubeVideo, int64, *YouTubePageToken, error)
}

type YouTubePageToken string
