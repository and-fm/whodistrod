package models

type TidalTrack struct {
	Data  []TidalTrackData `json:"data"`
	Links TidalLinks       `json:"links"`
}

type TidalTrackWithProviders struct {
	TidalTrack
	Included []TidalProvider `json:"included"`
}

type TidalProvider struct {
	ID         string                  `json:"id"`
	Type       string                  `json:"type"`
	Attributes TidalProviderAttributes `json:"attributes"`
}

type TidalTrackData struct {
	ID            string                  `json:"id"`
	Type          string                  `json:"type"`
	Attributes    TidalTrackAttributes    `json:"attributes"`
	Relationships TidalTrackRelationships `json:"relationships"`
}

type TidalTrackAttributes struct {
	Title         string              `json:"title"`
	ISRC          string              `json:"isrc"`
	Duration      string              `json:"duration"`
	Copyright     TidalCopyright      `json:"copyright"`
	Explicit      bool                `json:"explicit"`
	Popularity    float64             `json:"popularity"`
	AccessType    string              `json:"accessType"`
	Availability  []string            `json:"availability"`
	MediaTags     []string            `json:"mediaTags"`
	ExternalLinks []TidalExternalLink `json:"externalLinks"`
	Spotlighted   bool                `json:"spotlighted"`
	CreatedAt     string              `json:"createdAt"`
}

type TidalCopyright struct {
	Text string `json:"text"`
}

type TidalExternalLink struct {
	Href string        `json:"href"`
	Meta TidalLinkMeta `json:"meta"`
}

type TidalLinkMeta struct {
	Type string `json:"type"`
}

type TidalTrackRelationships struct {
	Albums          TidalRelationshipData  `json:"albums"`
	TrackStatistics TidalRelationshipLinks `json:"trackStatistics"`
	Artists         TidalRelationshipLinks `json:"artists"`
	Genres          TidalRelationshipLinks `json:"genres"`
	SimilarTracks   TidalRelationshipLinks `json:"similarTracks"`
	Owners          TidalRelationshipLinks `json:"owners"`
	Lyrics          TidalRelationshipLinks `json:"lyrics"`
	Providers       TidalRelationshipLinks `json:"providers"`
	Radio           TidalRelationshipLinks `json:"radio"`
}

type TidalRelationshipData struct {
	Data  []TidalIDType `json:"data"`
	Links TidalLinks    `json:"links"`
}

type TidalRelationshipLinks struct {
	Links TidalLinks `json:"links"`
}

type TidalIDType struct {
	ID   string `json:"id"`
	Type string `json:"type"`
}

type TidalLinks struct {
	Self string `json:"self"`
}

type TidalAlbum struct {
	ID            string                  `json:"id"`
	Type          string                  `json:"type"`
	Attributes    TidalAlbumAttributes    `json:"attributes"`
	Relationships TidalAlbumRelationships `json:"relationships"`
}

type TidalAlbumAttributes struct {
	Title           string              `json:"title"`
	BarcodeID       string              `json:"barcodeId"`
	NumberOfVolumes int                 `json:"numberOfVolumes"`
	NumberOfItems   int                 `json:"numberOfItems"`
	Duration        string              `json:"duration"`
	Explicit        bool                `json:"explicit"`
	ReleaseDate     string              `json:"releaseDate"`
	Copyright       string              `json:"copyright"`
	Popularity      float64             `json:"popularity"`
	Availability    []string            `json:"availability"`
	MediaTags       []string            `json:"mediaTags"`
	ExternalLinks   []TidalExternalLink `json:"externalLinks"`
	Type            string              `json:"type"`
}

type TidalAlbumRelationships struct {
	SimilarAlbums TidalRelationshipLinks `json:"similarAlbums"`
	Artists       TidalRelationshipLinks `json:"artists"`
	Genres        TidalRelationshipLinks `json:"genres"`
	Owners        TidalRelationshipLinks `json:"owners"`
	CoverArt      TidalRelationshipLinks `json:"coverArt"`
	Items         TidalRelationshipLinks `json:"items"`
	Providers     TidalRelationshipLinks `json:"providers"`
}

type TidalProviderAttributes struct {
	Name string `json:"name"`
}
