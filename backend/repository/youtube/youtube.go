package youtube

import (
	"context"
	"fmt"

	"google.golang.org/api/option"
	"google.golang.org/api/youtube/v3"
)

type YouTubeRepository struct {
	service *youtube.Service
}

func NewYouTubeRepository(ctx context.Context, apiKey string) (*YouTubeRepository, error) {
	service, err := youtube.NewService(ctx, option.WithAPIKey(apiKey))
	if err != nil {
		return nil, fmt.Errorf("failed to create youtube service: %w", err)
	}

	return &YouTubeRepository{
		service: service,
	}, nil
}
