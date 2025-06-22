package adapter

import (
	"context"
	"fmt"
	"net/url"
	"strings"
	"time"

	"github.com/tocoteron/omigoto/backend/module/youtube/model"
	"github.com/tocoteron/omigoto/backend/module/youtube/repository"
	"google.golang.org/api/option"
	"google.golang.org/api/youtube/v3"
)

const YouTubeMaxResults = 50

var _ repository.YouTubeRepository = &youtubeRepository{}

type youtubeRepository struct {
	service *youtube.Service
}

func NewYouTubeRepository(ctx context.Context, apiKey string) (*youtubeRepository, error) {
	service, err := youtube.NewService(ctx, option.WithAPIKey(apiKey))
	if err != nil {
		return nil, fmt.Errorf("failed to create youtube service: %w", err)
	}

	return &youtubeRepository{
		service: service,
	}, nil
}

func (r *youtubeRepository) GetChannel(
	ctx context.Context,
	channelID model.YouTubeChannelID,
) (*model.YouTubeChannel, error) {
	call := r.service.Channels.List([]string{"contentDetails", "snippet"}).
		Id(string(channelID)).
		MaxResults(1)

	response, err := call.Do()
	if err != nil {
		return nil, fmt.Errorf("failed to get channel: %w", err)
	}

	if len(response.Items) == 0 {
		return nil, fmt.Errorf("channel not found")
	}
	if len(response.Items) > 1 {
		return nil, fmt.Errorf("multiple channels found")
	}

	return &model.YouTubeChannel{
		YouTubeChannelIdentity: model.YouTubeChannelIdentity{
			ID:     channelID,
			Handle: model.YouTubeChannelHandle(response.Items[0].Snippet.CustomUrl),
		},
		UploadsPlaylistID: model.YouTubePlaylistID(response.Items[0].ContentDetails.RelatedPlaylists.Uploads),
	}, nil
}

func (r *youtubeRepository) GetPlaylist(
	ctx context.Context,
	playlistID model.YouTubePlaylistID,
) (*model.YouTubePlaylist, error) {
	call := r.service.Playlists.List([]string{"snippet"}).
		Id(string(playlistID)).
		MaxResults(1)

	response, err := call.Do()
	if err != nil {
		return nil, fmt.Errorf("failed to get playlist: %w", err)
	}

	if len(response.Items) == 0 {
		return nil, fmt.Errorf("playlist not found")
	}
	if len(response.Items) > 1 {
		return nil, fmt.Errorf("multiple playlists found")
	}

	return &model.YouTubePlaylist{
		ID:    model.YouTubePlaylistID(response.Items[0].Id),
		Title: response.Items[0].Snippet.Title,
	}, nil
}

func (r *youtubeRepository) ListPlaylists(
	ctx context.Context,
	channelID model.YouTubeChannelID,
	pageToken *repository.YouTubePageToken,
) ([]*model.YouTubePlaylist, int64, *repository.YouTubePageToken, error) {
	call := r.service.Playlists.List([]string{"snippet"}).
		ChannelId(string(channelID)).
		MaxResults(50) // max value is 50

	if pageToken != nil {
		call.PageToken(string(*pageToken))
	}

	response, err := call.Do()
	if err != nil {
		return nil, 0, nil, fmt.Errorf("failed to list playlists: %w", err)
	}

	playlists := make([]*model.YouTubePlaylist, 0, len(response.Items))
	for _, item := range response.Items {
		playlists = append(playlists, &model.YouTubePlaylist{
			ID:    model.YouTubePlaylistID(item.Id),
			Title: item.Snippet.Title,
		})
	}

	nextPageToken := pageTokenFromString(response.NextPageToken)

	return playlists, response.PageInfo.TotalResults, nextPageToken, nil
}

func (r *youtubeRepository) ListVideos(
	ctx context.Context,
	videoIDs []model.YouTubeVideoID,
	pageToken *repository.YouTubePageToken,
) ([]*model.YouTubeVideo, int64, *repository.YouTubePageToken, error) {
	if len(videoIDs) > YouTubeMaxResults {
		return nil, 0, nil, fmt.Errorf("videoIDs length must be less than or equal to %d", YouTubeMaxResults)
	}

	ids := make([]string, 0, len(videoIDs))
	for _, id := range videoIDs {
		ids = append(ids, string(id))
	}

	call := r.service.Videos.List([]string{"contentDetails", "snippet", "liveStreamingDetails"}).
		Id(ids...).
		MaxResults(YouTubeMaxResults)

	if pageToken != nil {
		call.PageToken(string(*pageToken))
	}

	response, err := call.Do()
	if err != nil {
		return nil, 0, nil, fmt.Errorf("failed to list videos: %w", err)
	}

	videos := make([]*model.YouTubeVideo, 0, len(response.Items))
	for _, item := range response.Items {
		video, err := videoFromYouTubeVideo(item)
		if err != nil {
			return nil, 0, nil, fmt.Errorf("failed to parse video: %w", err)
		}

		videos = append(videos, video)
	}

	nextPageToken := pageTokenFromString(response.NextPageToken)

	return videos, response.PageInfo.TotalResults, nextPageToken, nil
}

func (r *youtubeRepository) ListVideoIDsByPlaylist(
	ctx context.Context,
	playlistID model.YouTubePlaylistID,
	pageToken *repository.YouTubePageToken,
) ([]model.YouTubeVideoID, int64, *repository.YouTubePageToken, error) {
	call := r.service.PlaylistItems.List([]string{"snippet"}).
		PlaylistId(string(playlistID)).
		MaxResults(YouTubeMaxResults)

	if pageToken != nil {
		call.PageToken(string(*pageToken))
	}

	response, err := call.Do()
	if err != nil {
		return nil, 0, nil, fmt.Errorf("failed to list playlist items: %w", err)
	}

	videoIDs := make([]model.YouTubeVideoID, 0, len(response.Items))
	for _, item := range response.Items {
		videoIDs = append(videoIDs, model.YouTubeVideoID(item.Snippet.ResourceId.VideoId))
	}

	nextPageToken := pageTokenFromString(response.NextPageToken)

	return videoIDs, response.PageInfo.TotalResults, nextPageToken, nil
}

func pageTokenFromString(token string) *repository.YouTubePageToken {
	if token == "" {
		return nil
	}

	t := repository.YouTubePageToken(token)

	return &t
}

func videoFromYouTubeVideo(video *youtube.Video) (*model.YouTubeVideo, error) {
	var duration time.Duration
	if video.ContentDetails.Duration != "P0D" {
		dur, err := time.ParseDuration(strings.ToLower(strings.TrimPrefix(video.ContentDetails.Duration, "PT")))
		if err != nil {
			return nil, fmt.Errorf("failed to parse duration: %w", err)
		}
		duration = dur
	}

	var thumbnails model.YouTubeVideoThumbnails
	{
		defaultURL, err := thumbnailURLFromYouTubeThumbnail(video.Snippet.Thumbnails.Default)
		if err != nil {
			return nil, fmt.Errorf("failed to parse default thumbnail URL: %w", err)
		}

		mediumURL, err := thumbnailURLFromYouTubeThumbnail(video.Snippet.Thumbnails.Medium)
		if err != nil {
			return nil, fmt.Errorf("failed to parse medium thumbnail URL: %w", err)
		}

		highURL, err := thumbnailURLFromYouTubeThumbnail(video.Snippet.Thumbnails.High)
		if err != nil {
			return nil, fmt.Errorf("failed to parse high thumbnail URL: %w", err)
		}

		standardURL, err := thumbnailURLFromYouTubeThumbnail(video.Snippet.Thumbnails.Standard)
		if err != nil {
			return nil, fmt.Errorf("failed to parse standard thumbnail URL: %w", err)
		}

		maxresURL, err := thumbnailURLFromYouTubeThumbnail(video.Snippet.Thumbnails.Maxres)
		if err != nil {
			return nil, fmt.Errorf("failed to parse maxres thumbnail URL: %w", err)
		}

		thumbnails = model.YouTubeVideoThumbnails{
			Default:  defaultURL,
			Medium:   mediumURL,
			High:     highURL,
			Standard: standardURL,
			Maxres:   maxresURL,
		}
	}

	var liveStreamingDetails *model.YouTubeVideoLiveStreamingDetails
	if video.LiveStreamingDetails != nil && video.LiveStreamingDetails.ActualEndTime != "" {
		actualStartTime, err := time.Parse(time.RFC3339, video.LiveStreamingDetails.ActualStartTime)
		if err != nil {
			return nil, fmt.Errorf("failed to parse actual start time: %w", err)
		}

		actualEndTime, err := time.Parse(time.RFC3339, video.LiveStreamingDetails.ActualEndTime)
		if err != nil {
			return nil, fmt.Errorf("failed to parse actual end time: %w", err)
		}

		scheduledStartTime, err := time.Parse(time.RFC3339, video.LiveStreamingDetails.ScheduledStartTime)
		if err != nil {
			return nil, fmt.Errorf("failed to parse scheduled start time: %w", err)
		}

		liveStreamingDetails = &model.YouTubeVideoLiveStreamingDetails{
			ActualStartTime: actualStartTime,
			ActualEndTime:   actualEndTime,
			ScheduledStart:  scheduledStartTime,
		}
	}

	publishedAt, err := time.Parse(time.RFC3339, video.Snippet.PublishedAt)
	if err != nil {
		return nil, fmt.Errorf("failed to parse published at: %w", err)
	}

	return &model.YouTubeVideo{
		ID:                   model.YouTubeVideoID(video.Id),
		Title:                video.Snippet.Title,
		Description:          video.Snippet.Description,
		Duration:             duration,
		Thumbnails:           thumbnails,
		LiveStreamingDetails: liveStreamingDetails,
		PublishedAt:          publishedAt,
	}, nil
}

func thumbnailURLFromYouTubeThumbnail(thumbnail *youtube.Thumbnail) (*url.URL, error) {
	if thumbnail == nil {
		return nil, nil
	}

	thumbnailURL, err := url.Parse(thumbnail.Url)
	if err != nil {
		return nil, fmt.Errorf("failed to parse thumbnail URL: %w", err)
	}

	return thumbnailURL, nil
}
