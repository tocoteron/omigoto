package main

import (
	"context"
	"fmt"
	"log"

	"github.com/caarlos0/env/v11"
	"github.com/tocoteron/omigoto/backend/model"
	"github.com/tocoteron/omigoto/backend/model/adapter"
	"github.com/tocoteron/omigoto/backend/omikun"
)

type config struct {
	YouTubeAPIKey string `env:"YOUTUBE_API_KEY,notEmpty"`
}

func main() {
	var cfg config
	if err := env.Parse(&cfg); err != nil {
		log.Fatalf("failed to parse config: %v", err)
	}

	ctx := context.Background()

	youtubeRepo, err := adapter.NewYouTubeRepository(ctx, cfg.YouTubeAPIKey)
	if err != nil {
		log.Fatalf("failed to create youtube repository: %v", err)
	}

	channel, err := getChannel(ctx, youtubeRepo, omikun.YouTubeChannel.ID)
	if err != nil {
		log.Fatalf("failed to get channel: %v", err)
	}
	fmt.Printf("channel: %+v\n", channel)

	playlists, err := listAllPlaylists(ctx, youtubeRepo, omikun.YouTubeChannel.ID)
	if err != nil {
		log.Fatalf("failed to list playlists: %v", err)
	}
	fmt.Printf("playlists: %+v\n", playlists)

	uploadsPlaylist, err := getUploadsPlaylist(ctx, youtubeRepo, omikun.YouTubeChannel.ID)
	if err != nil {
		log.Fatalf("failed to get uploads playlist: %v", err)
	}
	fmt.Printf("uploadsPlaylist: %+v\n", uploadsPlaylist)

	videoIDs, err := listAllVideoIDs(ctx, youtubeRepo, uploadsPlaylist.ID)
	if err != nil {
		log.Fatalf("failed to list video IDs: %v", err)
	}
	fmt.Printf("videoIDs: %+v\n", videoIDs)

	videos, err := listAllVideos(ctx, youtubeRepo, videoIDs[:50])
	if err != nil {
		log.Fatalf("failed to list videos: %v", err)
	}
	fmt.Printf("videos: %+v\n", videos)
}

func getChannel(
	ctx context.Context,
	youtubeRepo model.YouTubeRepository,
	channelID model.YouTubeChannelID,
) (*model.YouTubeChannel, error) {
	channel, err := youtubeRepo.GetChannel(ctx, channelID)
	if err != nil {
		return nil, fmt.Errorf("failed to get channel: %w", err)
	}

	return channel, nil
}

func getUploadsPlaylist(
	ctx context.Context,
	youtubeRepo model.YouTubeRepository,
	channelID model.YouTubeChannelID,
) (*model.YouTubePlaylist, error) {
	uploadsPlaylist, err := youtubeRepo.GetUploadsPlaylist(ctx, channelID)
	if err != nil {
		return nil, fmt.Errorf("failed to get uploads playlist: %w", err)
	}

	return uploadsPlaylist, nil
}

func listAllPlaylists(
	ctx context.Context,
	youtubeRepo model.YouTubeRepository,
	channelID model.YouTubeChannelID,
) ([]*model.YouTubePlaylist, error) {
	playlists := make([]*model.YouTubePlaylist, 0)

	var pageToken *model.YouTubePageToken
	for {
		pls, _, nextPageToken, err := youtubeRepo.ListPlaylists(ctx, channelID, pageToken)
		if err != nil {
			return nil, fmt.Errorf("failed to list playlists: %w", err)
		}

		playlists = append(playlists, pls...)

		if nextPageToken == nil {
			break
		}

		pageToken = nextPageToken
	}

	return playlists, nil
}

func listAllVideoIDs(
	ctx context.Context,
	youtubeRepo model.YouTubeRepository,
	playlistID model.YouTubePlaylistID,
) ([]model.YouTubeVideoID, error) {
	videoIDs := make([]model.YouTubeVideoID, 0)

	var pageToken *model.YouTubePageToken
	for {
		ids, _, nextPageToken, err := youtubeRepo.ListVideoIDs(ctx, playlistID, pageToken)
		if err != nil {
			return nil, fmt.Errorf("failed to list video IDs: %w", err)
		}

		videoIDs = append(videoIDs, ids...)

		if nextPageToken == nil {
			break
		}

		pageToken = nextPageToken
	}

	return videoIDs, nil
}

func listAllVideos(
	ctx context.Context,
	youtubeRepo model.YouTubeRepository,
	videoIDs []model.YouTubeVideoID,
) ([]*model.YouTubeVideo, error) {
	videos := make([]*model.YouTubeVideo, 0)

	var pageToken *model.YouTubePageToken
	for {
		vs, _, nextPageToken, err := youtubeRepo.ListVideos(ctx, videoIDs, pageToken)
		if err != nil {
			return nil, fmt.Errorf("failed to list videos: %w", err)
		}

		videos = append(videos, vs...)

		if nextPageToken == nil {
			break
		}

		pageToken = nextPageToken
	}

	return videos, nil
}
