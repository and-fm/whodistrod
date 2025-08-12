package models

type Track struct {
	Data     []TrackData `json:"data"`
	Links    Links     `json:"links"`
}

type TrackWithProviders struct {
	Track
	Included []Provider `json:"included"`
}

type Provider struct {
	ID            string             `json:"id"`
	Type          string             `json:"type"`
	Attributes    ProviderAttributes    `json:"attributes"`
}

type TrackData struct {
	ID            string             `json:"id"`
	Type          string             `json:"type"`
	Attributes    TrackAttributes    `json:"attributes"`
	Relationships TrackRelationships `json:"relationships"`
}

type TrackAttributes struct {
	Title         string         `json:"title"`
	ISRC          string         `json:"isrc"`
	Duration      string         `json:"duration"`
	Copyright     string         `json:"copyright"`
	Explicit      bool           `json:"explicit"`
	Popularity    float64        `json:"popularity"`
	AccessType    string         `json:"accessType"`
	Availability  []string       `json:"availability"`
	MediaTags     []string       `json:"mediaTags"`
	ExternalLinks []ExternalLink `json:"externalLinks"`
	Spotlighted   bool           `json:"spotlighted"`
	CreatedAt     string         `json:"createdAt"`
}

type ExternalLink struct {
	Href string   `json:"href"`
	Meta LinkMeta `json:"meta"`
}

type LinkMeta struct {
	Type string `json:"type"`
}

type TrackRelationships struct {
	Albums          RelationshipData  `json:"albums"`
	TrackStatistics RelationshipLinks `json:"trackStatistics"`
	Artists         RelationshipLinks `json:"artists"`
	Genres          RelationshipLinks `json:"genres"`
	SimilarTracks   RelationshipLinks `json:"similarTracks"`
	Owners          RelationshipLinks `json:"owners"`
	Lyrics          RelationshipLinks `json:"lyrics"`
	Providers       RelationshipLinks `json:"providers"`
	Radio           RelationshipLinks `json:"radio"`
}

type RelationshipData struct {
	Data  []IDType `json:"data"`
	Links Links    `json:"links"`
}

type RelationshipLinks struct {
	Links Links `json:"links"`
}

type IDType struct {
	ID   string `json:"id"`
	Type string `json:"type"`
}

type Links struct {
	Self string `json:"self"`
}

type Album struct {
	ID            string             `json:"id"`
	Type          string             `json:"type"`
	Attributes    AlbumAttributes    `json:"attributes"`
	Relationships AlbumRelationships `json:"relationships"`
}

type AlbumAttributes struct {
	Title           string         `json:"title"`
	BarcodeID       string         `json:"barcodeId"`
	NumberOfVolumes int            `json:"numberOfVolumes"`
	NumberOfItems   int            `json:"numberOfItems"`
	Duration        string         `json:"duration"`
	Explicit        bool           `json:"explicit"`
	ReleaseDate     string         `json:"releaseDate"`
	Copyright       string         `json:"copyright"`
	Popularity      float64        `json:"popularity"`
	Availability    []string       `json:"availability"`
	MediaTags       []string       `json:"mediaTags"`
	ExternalLinks   []ExternalLink `json:"externalLinks"`
	Type            string         `json:"type"`
}

type AlbumRelationships struct {
	SimilarAlbums RelationshipLinks `json:"similarAlbums"`
	Artists       RelationshipLinks `json:"artists"`
	Genres        RelationshipLinks `json:"genres"`
	Owners        RelationshipLinks `json:"owners"`
	CoverArt      RelationshipLinks `json:"coverArt"`
	Items         RelationshipLinks `json:"items"`
	Providers     RelationshipLinks `json:"providers"`
}

type ProviderAttributes struct {
	Name string `json:"name"`
}