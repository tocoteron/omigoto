package repository

import (
	"context"

	"github.com/tocoteron/omigoto/backend/module/youtube/model"
)

type YouTubeRepository interface {
	// Channel operations
	GetChannel(ctx context.Context, channelID model.YouTubeChannelID) (*model.YouTubeChannel, error)

	// Playlist operations
	GetUploadsPlaylist(ctx context.Context, channelID model.YouTubeChannelID) (*model.YouTubePlaylist, error)
	ListPlaylists(ctx context.Context, channelID model.YouTubeChannelID, pageToken *YouTubePageToken) ([]*model.YouTubePlaylist, int64, *YouTubePageToken, error)

	// Video operations
	ListVideos(ctx context.Context, videoIDs []model.YouTubeVideoID, pageToken *YouTubePageToken) ([]*model.YouTubeVideo, int64, *YouTubePageToken, error)
	ListVideoIDsByPlaylist(ctx context.Context, playlistID model.YouTubePlaylistID, pageToken *YouTubePageToken) ([]model.YouTubeVideoID, int64, *YouTubePageToken, error)
}

type YouTubeDBRepository interface {
	// Channel operations
	CreateChannel(ctx context.Context, channel *model.YouTubeChannel) error
	GetChannel(ctx context.Context, channelID model.YouTubeChannelID) (*model.YouTubeChannel, error)
	GetChannelByHandle(ctx context.Context, handle model.YouTubeChannelHandle) (*model.YouTubeChannel, error)

	// Playlist operations
	CreatePlaylist(ctx context.Context, channelID model.YouTubeChannelID, playlist *model.YouTubePlaylist) error
	GetPlaylist(ctx context.Context, playlistID model.YouTubePlaylistID) (*model.YouTubePlaylist, error)
	ListPlaylists(ctx context.Context, playlistIDs []model.YouTubePlaylistID) ([]*model.YouTubePlaylist, error)
	ListPlaylistIDsByChannel(ctx context.Context, channelID model.YouTubeChannelID) ([]model.YouTubePlaylistID, error)

	// Video operations
	CreateVideo(ctx context.Context, video *model.YouTubeVideo) error
	GetVideo(ctx context.Context, videoID model.YouTubeVideoID) (*model.YouTubeVideo, error)
	ListVideos(ctx context.Context, videoIDs []model.YouTubeVideoID) ([]*model.YouTubeVideo, error)

	// Playlist-Video relationship operations
	CreatePlaylistVideo(ctx context.Context, playlistID model.YouTubePlaylistID, videoID model.YouTubeVideoID) error
	ListVideoIDsByPlaylist(ctx context.Context, playlistID model.YouTubePlaylistID) ([]model.YouTubeVideoID, error)

	// Live streaming details operations
	CreateVideoLiveStreamingDetails(ctx context.Context, videoID model.YouTubeVideoID, details *model.YouTubeVideoLiveStreamingDetails) error
	GetVideoLiveStreamingDetails(ctx context.Context, videoID model.YouTubeVideoID) (*model.YouTubeVideoLiveStreamingDetails, error)
}

type YouTubePageToken string
