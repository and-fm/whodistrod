package models

type SpotifyImage struct {
	URL    string `json:"url"`
	Width  int    `json:"width"`
	Height int    `json:"height"`
}

type SpotifyArtist struct {
	ExternalURLs map[string]string `json:"external_urls"`
	Href         string            `json:"href"`
	ID           string            `json:"id"`
	Name         string            `json:"name"`
	Type         string            `json:"type"`
	URI          string            `json:"uri"`
}

type SpotifyAlbum struct {
	AlbumType            string            `json:"album_type"`
	Artists              []SpotifyArtist   `json:"artists"`
	ExternalURLs         map[string]string `json:"external_urls"`
	Href                 string            `json:"href"`
	ID                   string            `json:"id"`
	Images               []SpotifyImage    `json:"images"`
	IsPlayable           bool              `json:"is_playable"`
	Name                 string            `json:"name"`
	ReleaseDate          string            `json:"release_date"`
	ReleaseDatePrecision string            `json:"release_date_precision"`
	TotalTracks          int               `json:"total_tracks"`
	Type                 string            `json:"type"`
	URI                  string            `json:"uri"`
}

type SpotifyTrack struct {
	Album        SpotifyAlbum      `json:"album"`
	Artists      []SpotifyArtist   `json:"artists"`
	DiscNumber   int               `json:"disc_number"`
	DurationMS   int               `json:"duration_ms"`
	Explicit     bool              `json:"explicit"`
	ExternalIDs  SpotifyExternalID `json:"external_ids"`
	ExternalURLs map[string]string `json:"external_urls"`
	Href         string            `json:"href"`
	ID           string            `json:"id"`
	IsLocal      bool              `json:"is_local"`
	IsPlayable   bool              `json:"is_playable"`
	Name         string            `json:"name"`
	Popularity   int               `json:"popularity"`
	PreviewURL   *string           `json:"preview_url"`
	TrackNumber  int               `json:"track_number"`
	Type         string            `json:"type"`
	URI          string            `json:"uri"`
}

type SpotifyExternalID struct {
	ISRC string `json:"isrc"`
}
