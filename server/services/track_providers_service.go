package services

import (
	"errors"
	"fmt"

	dspclients "github.com/and-fm/whodistrod/dsp_clients"
	"github.com/and-fm/whodistrod/models"
	"go.uber.org/dig"
)

type trackProvidersService struct {
	dig.In

	SpotifyClient dspclients.SpotifyClient
	TidalClient   dspclients.TidalClient
}

type TrackProvidersService interface {
	GetTrackProviders(trackUrl string) (models.ProviderApiResponse, error)
}

func NewTrackProvidersService(s trackProvidersService) TrackProvidersService {
	return &s
}

func (s *trackProvidersService) GetTrackProviders(trackUrl string) (res models.ProviderApiResponse, err error) {
	trackService, err := getServiceFromTrackURL(trackUrl)

	if err != nil {
		return res, err
	}

	switch trackService {
	case "spotify":
		trackId, err := parseSpotifyTrackId(trackUrl)
		if err != nil {
			return res, err
		}

		track, err := s.GetSpotifyTrackById(trackId)
		if err != nil {
			return res, fmt.Errorf("get spotify track by id: %w", err)
		}

		if track.ExternalIDs.ISRC == "" {
			return res, errors.New("could not find ISRC for spotify track")
		}
		tidalTrack, err := s.GetTidalTrackIdByIsrc(track.ExternalIDs.ISRC)
		if err != nil {
			return res, fmt.Errorf("get tidal track by isrc: %w", err)
		}

		if len(tidalTrack.Included) == 0 {
			return res, errors.New("could not find tidal track providers")
		}

		return models.ProviderApiResponse{Provider: tidalTrack.Included[0].Attributes.Name}, nil
	case "tidal":
		trackId, err := parseTidalTrackId(trackUrl)
		if err != nil {
			return res, err
		}

		track, err := s.GetTidalTrackProvidersById(trackId)
		if err != nil {
			return res, err
		}

		if len(track.Included) == 0 {
			return res, errors.New("could not find tidal track providers")
		}

		return models.ProviderApiResponse{Provider: track.Included[0].Attributes.Name}, nil
	default:
		return models.ProviderApiResponse{}, errors.New("unsupported track service" + string(trackService))
	}
}

func (s *trackProvidersService) GetTidalTrackProvidersById(trackId string) (models.TidalTrackWithProviders, error) {
	trackProviderUrl := "https://openapi.tidal.com/v2/tracks/" + trackId + "/relationships/providers?countryCode=US&include=providers"
	_, track, err := dspclients.Do[models.TidalTrackWithProviders](s.TidalClient, "GET", trackProviderUrl, nil)
	return track, err
}

func (s *trackProvidersService) GetTidalTrackIdByIsrc(isrc string) (models.TidalTrackWithProviders, error) {
	trackProviderUrl := "https://openapi.tidal.com/v2/tracks?countryCode=US&include=providers&filter[isrc]=" + isrc
	_, track, err := dspclients.Do[models.TidalTrackWithProviders](s.TidalClient, "GET", trackProviderUrl, nil)
	return track, err
}

func (s *trackProvidersService) GetSpotifyTrackById(trackId string) (models.SpotifyTrack, error) {
	trackUrl := "https://api.spotify.com/v1/tracks/" + trackId + "?market=US"
	_, track, err := dspclients.Do[models.SpotifyTrack](s.SpotifyClient, "GET", trackUrl, nil)
	return track, err
}
