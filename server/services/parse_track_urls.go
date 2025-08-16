package services

import (
	"fmt"
	"net/url"
	"regexp"
	"strings"
)

type ServiceType string

const (
	ServiceTidal        ServiceType = "tidal"
	ServiceSpotify      ServiceType = "spotify"
	ServiceAppleMusic   ServiceType = "apple_music"
	ServiceDeezer       ServiceType = "deezer"
	ServiceYoutubeMusic ServiceType = "youtube_music"
	ServiceAmazonMusic  ServiceType = "amazon_music"
	ServiceUnknown      ServiceType = "unknown"
)

func getServiceFromTrackURL(trackURL string) (ServiceType, error) {
	parsedURL, err := url.Parse(trackURL)
	if err != nil {
		return ServiceUnknown, fmt.Errorf("failed to parse URL: %w", err)
	}

	switch {
	case strings.Contains(parsedURL.Host, "tidal.com"):
		return ServiceTidal, nil
	case strings.Contains(parsedURL.Host, "spotify.com"):
		return ServiceSpotify, nil
	case strings.Contains(parsedURL.Host, "music.apple.com"):
		return ServiceAppleMusic, nil
	case strings.Contains(parsedURL.Host, "deezer.com"):
		return ServiceDeezer, nil
	case strings.Contains(parsedURL.Host, "music.youtube.com"):
		return ServiceYoutubeMusic, nil
	case strings.Contains(parsedURL.Host, "music.amazon.com"):
		return ServiceAmazonMusic, nil
	default:
		return ServiceUnknown, nil
	}
}

func parseTidalTrackId(trackURL string) (trackId string, err error) {
	parsedURL, err := url.Parse(trackURL)
	if err != nil {
		return "", fmt.Errorf("failed to parse URL: %w", err)
	}

	patterns := []string{
		`/browse/track/(\d+)`,
		`/track/(\d+)`,
		`/album/.*/track/(\d+)`,
	}

	for _, pattern := range patterns {
		re := regexp.MustCompile(pattern)
		matches := re.FindStringSubmatch(parsedURL.Path)
		if len(matches) > 0 {
			return matches[1], nil
		}
	}

	return "", nil
}

func parseSpotifyTrackId(trackURL string) (trackId string, err error) {
	parsedURL, err := url.Parse(trackURL)
	if err != nil {
		return "", fmt.Errorf("failed to parse URL: %w", err)
	}

	re := regexp.MustCompile(`/track/([a-zA-Z0-9]{22})`)
	matches := re.FindStringSubmatch(parsedURL.Path)
	if len(matches) > 1 {
		return matches[1], nil
	}

	return parsedURL.Path, nil
}
