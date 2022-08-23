package main

type songwhipLinks struct {
	Qobuz        bool
	Tidal        bool
	Amazon       bool
	Deezer       bool
	Itunes       bool
	Napster      bool
	Pandora      bool
	Spotify      bool
	Youtube      bool
	Audiomack    bool
	LineMusic    bool
	AmazonMusic  bool
	ItunesStore  bool
	YoutubeMusic bool
}

type songwhipArtist struct {
	Type           string
	Id             int
	Path           string
	Name           string
	SourceUrl      string
	SourceCountry  string
	Url            string
	Image          string
	CreatedAt      string
	UpdatedAt      string
	RefreshedAt    string
	LinksCountries []string
	Links          songwhipArtistLinks
	Description    string
	ServiceIds     songwhipServices
	OrchardId      string
	SpotifyId      string
}

type songwhipArtistLink struct {
	Link      string
	Countries []string
}

type songwhipArtistLinks struct {
	Tidal        []songwhipArtistLink
	Amazon       []songwhipArtistLink
	Deezer       []songwhipArtistLink
	Itunes       []songwhipArtistLink
	Yandex       []songwhipArtistLink
	Discogs      []songwhipArtistLink
	Napster      []songwhipArtistLink
	Pandora      []songwhipArtistLink
	Spotify      []songwhipArtistLink
	Twitter      []songwhipArtistLink
	Youtube      []songwhipArtistLink
	Facebook     []songwhipArtistLink
	Instagram    []songwhipArtistLink
	Wikipedia    []songwhipArtistLink
	Soundcloud   []songwhipArtistLink
	AmazonMusic  []songwhipArtistLink
	ItunesStore  []songwhipArtistLink
	MusicBrainz  []songwhipArtistLink
	YoutubeMusic []songwhipArtistLink
}

type songwhipServices struct {
	Tidal       string
	Amazon      string
	Deezer      string
	Itunes      string
	Discogs     string
	Napster     string
	Pandora     string
	Spotify     string
	Googleplay  string
	Soundcloud  string
	MusicBrainz string
}

type songwhipResponse struct {
	Type           string
	Id             int
	Path           string
	Name           string
	Url            string
	SourceUrl      string
	SourceCountry  string
	ReleaseDate    string
	CreatedAt      string
	UpdatedAt      string
	RefreshedAt    string
	Image          string
	Isrc           string
	Links          songwhipLinks
	LinksCountries []string
	Artists        []songwhipArtist
}
