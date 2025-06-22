package adapter

import (
	"context"
	"errors"
	"fmt"
	"net/url"

	"github.com/jackc/pgx/v5"
	"github.com/tocoteron/omigoto/backend/gen/db"
	"github.com/tocoteron/omigoto/backend/module/youtube/model"
	"github.com/tocoteron/omigoto/backend/module/youtube/repository"
)

var _ repository.YouTubeDBRepository = &youtubeDBRepository{}

type youtubeDBRepository struct {
	q db.Querier
}

func NewYouTubeDBRepository(q db.Querier) repository.YouTubeDBRepository {
	return &youtubeDBRepository{
		q: q,
	}
}

// ----- Channel operations -----

func (r *youtubeDBRepository) CreateChannel(ctx context.Context, channel *model.YouTubeChannel) error {
	err := r.q.CreateYouTubeChannel(ctx, db.CreateYouTubeChannelParams{
		ChannelID:         string(channel.ID),
		Handle:            string(channel.Handle),
		UploadsPlaylistID: string(channel.UploadsPlaylistID),
	})
	if err != nil {
		return fmt.Errorf("failed to create channel: %w", err)
	}
	return nil
}

func (r *youtubeDBRepository) GetChannel(ctx context.Context, channelID model.YouTubeChannelID) (*model.YouTubeChannel, error) {
	dbChannel, err := r.q.GetYouTubeChannel(ctx, string(channelID))
	if err != nil {
		return nil, fmt.Errorf("failed to get channel: %w", err)
	}

	return &model.YouTubeChannel{
		YouTubeChannelIdentity: model.YouTubeChannelIdentity{
			ID:     model.YouTubeChannelID(dbChannel.ChannelID),
			Handle: model.YouTubeChannelHandle(dbChannel.Handle),
		},
		UploadsPlaylistID: model.YouTubePlaylistID(dbChannel.UploadsPlaylistID),
	}, nil
}

func (r *youtubeDBRepository) GetChannelByHandle(ctx context.Context, handle model.YouTubeChannelHandle) (*model.YouTubeChannel, error) {
	dbChannel, err := r.q.GetYouTubeChannelByHandle(ctx, string(handle))
	if err != nil {
		return nil, fmt.Errorf("failed to get channel by handle: %w", err)
	}

	return &model.YouTubeChannel{
		YouTubeChannelIdentity: model.YouTubeChannelIdentity{
			ID:     model.YouTubeChannelID(dbChannel.ChannelID),
			Handle: model.YouTubeChannelHandle(dbChannel.Handle),
		},
		UploadsPlaylistID: model.YouTubePlaylistID(dbChannel.UploadsPlaylistID),
	}, nil
}

// ----- Playlist operations -----

func (r *youtubeDBRepository) CreatePlaylist(ctx context.Context, channelID model.YouTubeChannelID, playlist *model.YouTubePlaylist) error {
	err := r.q.CreateYouTubePlaylist(ctx, db.CreateYouTubePlaylistParams{
		PlaylistID: string(playlist.ID),
		ChannelID:  string(channelID),
		Title:      playlist.Title,
	})
	if err != nil {
		return fmt.Errorf("failed to create playlist: %w", err)
	}
	return nil
}

func (r *youtubeDBRepository) GetPlaylist(ctx context.Context, playlistID model.YouTubePlaylistID) (*model.YouTubePlaylist, error) {
	dbPlaylist, err := r.q.GetYouTubePlaylist(ctx, string(playlistID))
	if err != nil {
		return nil, fmt.Errorf("failed to get playlist: %w", err)
	}

	return &model.YouTubePlaylist{
		ID:    model.YouTubePlaylistID(dbPlaylist.PlaylistID),
		Title: dbPlaylist.Title,
	}, nil
}

func (r *youtubeDBRepository) ListPlaylists(ctx context.Context, playlistIDs []model.YouTubePlaylistID) ([]*model.YouTubePlaylist, error) {
	ids := make([]string, len(playlistIDs))
	for i, id := range playlistIDs {
		ids[i] = string(id)
	}

	dbPlaylists, err := r.q.ListPlaylists(ctx, ids)
	if err != nil {
		return nil, fmt.Errorf("failed to list playlists: %w", err)
	}

	playlists := make([]*model.YouTubePlaylist, len(dbPlaylists))
	for i, dbPlaylist := range dbPlaylists {
		playlists[i] = &model.YouTubePlaylist{
			ID:    model.YouTubePlaylistID(dbPlaylist.PlaylistID),
			Title: dbPlaylist.Title,
		}
	}

	return playlists, nil
}

func (r *youtubeDBRepository) ListPlaylistIDsByChannel(ctx context.Context, channelID model.YouTubeChannelID) ([]model.YouTubePlaylistID, error) {
	ids, err := r.q.ListPlaylistIDsByChannel(ctx, string(channelID))
	if err != nil {
		return nil, fmt.Errorf("failed to list playlist IDs by channel: %w", err)
	}

	playlistIDs := make([]model.YouTubePlaylistID, len(ids))
	for i, id := range ids {
		playlistIDs[i] = model.YouTubePlaylistID(id)
	}

	return playlistIDs, nil
}

// ----- Video operations -----

func (r *youtubeDBRepository) CreateVideo(ctx context.Context, video *model.YouTubeVideo) error {
	err := r.q.CreateYouTubeVideo(ctx, db.CreateYouTubeVideoParams{
		VideoID:              string(video.ID),
		Title:                video.Title,
		Description:          video.Description,
		Duration:             video.Duration,
		ThumbnailDefaultUrl:  urlToString(video.Thumbnails.Default),
		ThumbnailMediumUrl:   urlToString(video.Thumbnails.Medium),
		ThumbnailHighUrl:     urlToString(video.Thumbnails.High),
		ThumbnailStandardUrl: urlToString(video.Thumbnails.Standard),
		ThumbnailMaxresUrl:   urlToString(video.Thumbnails.Maxres),
		PublishedAt:          video.PublishedAt,
	})
	if err != nil {
		return fmt.Errorf("failed to create video: %w", err)
	}

	// Create live streaming details if available
	if video.LiveStreamingDetails != nil {
		err = r.CreateVideoLiveStreamingDetails(ctx, video.ID, video.LiveStreamingDetails)
		if err != nil {
			return fmt.Errorf("failed to create video live streaming details: %w", err)
		}
	}

	return nil
}

func (r *youtubeDBRepository) GetVideo(ctx context.Context, videoID model.YouTubeVideoID) (*model.YouTubeVideo, error) {
	dbVideo, err := r.q.GetYouTubeVideo(ctx, string(videoID))
	if err != nil {
		return nil, fmt.Errorf("failed to get video: %w", err)
	}

	video, err := convertYouTubeVideo(dbVideo)
	if err != nil {
		return nil, fmt.Errorf("failed to convert video: %w", err)
	}

	// Get live streaming details if available
	dbLiveDetails, err := r.q.GetYouTubeVideoLiveStreamingDetails(ctx, string(videoID))
	if err == nil {
		video.LiveStreamingDetails = convertYouTubeVideoLiveStreamingDetails(dbLiveDetails)
	} else if !errors.Is(err, pgx.ErrNoRows) {
		return nil, fmt.Errorf("failed to get video live streaming details: %w", err)
	}

	return video, nil
}

func (r *youtubeDBRepository) ListVideos(ctx context.Context, videoIDs []model.YouTubeVideoID) ([]*model.YouTubeVideo, error) {
	ids := make([]string, len(videoIDs))
	for i, id := range videoIDs {
		ids[i] = string(id)
	}

	dbVideos, err := r.q.ListYouTubeVideos(ctx, ids)
	if err != nil {
		return nil, fmt.Errorf("failed to list videos: %w", err)
	}

	videos := make([]*model.YouTubeVideo, len(dbVideos))
	for i, dbVideo := range dbVideos {
		video, err := convertYouTubeVideo(dbVideo)
		if err != nil {
			return nil, fmt.Errorf("failed to convert video: %w", err)
		}

		// Get live streaming details if available
		dbLiveDetails, err := r.q.GetYouTubeVideoLiveStreamingDetails(ctx, dbVideo.VideoID)
		if err == nil {
			video.LiveStreamingDetails = convertYouTubeVideoLiveStreamingDetails(dbLiveDetails)
		} else if !errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("failed to get video live streaming details: %w", err)
		}

		videos[i] = video
	}

	return videos, nil
}

// ----- Playlist-Video relationship operations -----

func (r *youtubeDBRepository) CreatePlaylistVideo(ctx context.Context, playlistID model.YouTubePlaylistID, videoID model.YouTubeVideoID) error {
	err := r.q.CreateYouTubePlaylistVideo(ctx, db.CreateYouTubePlaylistVideoParams{
		PlaylistID: string(playlistID),
		VideoID:    string(videoID),
	})
	if err != nil {
		return fmt.Errorf("failed to create playlist video: %w", err)
	}
	return nil
}

func (r *youtubeDBRepository) ListVideoIDsByPlaylist(ctx context.Context, playlistID model.YouTubePlaylistID) ([]model.YouTubeVideoID, error) {
	ids, err := r.q.ListYouTubePlaylistVideoIDs(ctx, string(playlistID))
	if err != nil {
		return nil, fmt.Errorf("failed to list video IDs by playlist: %w", err)
	}

	videoIDs := make([]model.YouTubeVideoID, len(ids))
	for i, id := range ids {
		videoIDs[i] = model.YouTubeVideoID(id)
	}

	return videoIDs, nil
}

// ----- Live streaming details operations -----

func (r *youtubeDBRepository) CreateVideoLiveStreamingDetails(ctx context.Context, videoID model.YouTubeVideoID, details *model.YouTubeVideoLiveStreamingDetails) error {
	err := r.q.CreateYouTubeVideoLiveStreamingDetails(ctx, db.CreateYouTubeVideoLiveStreamingDetailsParams{
		VideoID:            string(videoID),
		ActualStartTime:    details.ActualStartTime,
		ActualEndTime:      details.ActualEndTime,
		ScheduledStartTime: details.ScheduledStart,
	})
	if err != nil {
		return fmt.Errorf("failed to create video live streaming details: %w", err)
	}
	return nil
}

func (r *youtubeDBRepository) GetVideoLiveStreamingDetails(ctx context.Context, videoID model.YouTubeVideoID) (*model.YouTubeVideoLiveStreamingDetails, error) {
	dbDetails, err := r.q.GetYouTubeVideoLiveStreamingDetails(ctx, string(videoID))
	if err != nil {
		return nil, fmt.Errorf("failed to get video live streaming details: %w", err)
	}

	return &model.YouTubeVideoLiveStreamingDetails{
		ActualStartTime: dbDetails.ActualStartTime,
		ActualEndTime:   dbDetails.ActualEndTime,
		ScheduledStart:  dbDetails.ScheduledStartTime,
	}, nil
}

// ----- Converters -----

func convertYouTubeVideo(dbVideo db.YoutubeVideo) (*model.YouTubeVideo, error) {
	thumbnails, err := convertYouTubeVideoThumbnails(dbVideo)
	if err != nil {
		return nil, fmt.Errorf("failed to convert thumbnails: %w", err)
	}

	return &model.YouTubeVideo{
		ID:          model.YouTubeVideoID(dbVideo.VideoID),
		Title:       dbVideo.Title,
		Description: dbVideo.Description,
		Duration:    dbVideo.Duration,
		Thumbnails:  *thumbnails,
		PublishedAt: dbVideo.PublishedAt,
	}, nil
}

func convertYouTubeVideoThumbnails(dbVideo db.YoutubeVideo) (*model.YouTubeVideoThumbnails, error) {
	thumbnailDefaultURL, err := stringToURL(dbVideo.ThumbnailDefaultUrl)
	if err != nil {
		return nil, fmt.Errorf("failed to parse thumbnail default URL: %w", err)
	}

	thumbnailMediumURL, err := stringToURL(dbVideo.ThumbnailMediumUrl)
	if err != nil {
		return nil, fmt.Errorf("failed to parse thumbnail medium URL: %w", err)
	}

	thumbnailHighURL, err := stringToURL(dbVideo.ThumbnailHighUrl)
	if err != nil {
		return nil, fmt.Errorf("failed to parse thumbnail high URL: %w", err)
	}

	thumbnailStandardURL, err := stringToURL(dbVideo.ThumbnailStandardUrl)
	if err != nil {
		return nil, fmt.Errorf("failed to parse thumbnail standard URL: %w", err)
	}

	thumbnailMaxresURL, err := stringToURL(dbVideo.ThumbnailMaxresUrl)
	if err != nil {
		return nil, fmt.Errorf("failed to parse thumbnail maxres URL: %w", err)
	}

	return &model.YouTubeVideoThumbnails{
		Default:  thumbnailDefaultURL,
		Medium:   thumbnailMediumURL,
		High:     thumbnailHighURL,
		Standard: thumbnailStandardURL,
		Maxres:   thumbnailMaxresURL,
	}, nil
}

func convertYouTubeVideoLiveStreamingDetails(dbLiveDetails db.YoutubeVideoLiveStreamingDetail) *model.YouTubeVideoLiveStreamingDetails {
	return &model.YouTubeVideoLiveStreamingDetails{
		ActualStartTime: dbLiveDetails.ActualStartTime,
		ActualEndTime:   dbLiveDetails.ActualEndTime,
		ScheduledStart:  dbLiveDetails.ScheduledStartTime,
	}
}

// ----- Helper functions -----

func urlToString(u *url.URL) *string {
	if u == nil {
		return nil
	}

	s := u.String()

	return &s
}

func stringToURL(s *string) (*url.URL, error) {
	if s == nil {
		return nil, nil
	}

	u, err := url.Parse(*s)
	if err != nil {
		return nil, fmt.Errorf("failed to parse URL: %w", err)
	}

	return u, nil
}
