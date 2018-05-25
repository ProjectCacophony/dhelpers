package dhelpers

import (
	"strconv"
	"time"

	"context"

	"github.com/Seklfreak/lastfm-go/lastfm"
	"github.com/opentracing/opentracing-go"
	"gitlab.com/Cacophony/dhelpers/cache"
)

// LastFmPeriod is a type for periods used for Last.FM requests
type LastFmPeriod string

// defines possible LastFM periods
const (
	LastFmPeriodOverall LastFmPeriod = "overall"
	LastFmPeriod7day                 = "7day"
	LastFmPeriod1month               = "1month"
	LastFmPeriod3month               = "3month"
	LastFmPeriod6month               = "6month"
	LastFmPeriod12month              = "12month"
)

const (
	lastFmTargetImageSize = "extralarge"
)

// LastfmUserData contains information about an User on LastFM
type LastfmUserData struct {
	Username        string
	Name            string
	Icon            string
	Scrobbles       int
	Country         string
	AccountCreation time.Time
}

// LastfmTrackData contains information about a Track on LastFM
type LastfmTrackData struct {
	Name           string
	URL            string
	ImageURL       string
	Artist         string
	ArtistURL      string
	ArtistImageURL string
	Album          string
	Time           time.Time
	Loved          bool
	NowPlaying     bool
	Scrobbles      int
	// used for guild stats
	Users int
}

// LastfmArtistData contains information about an Artist on LastFM
type LastfmArtistData struct {
	Name      string
	URL       string
	ImageURL  string
	Scrobbles int
}

// LastfmAlbumData contains information about an Album on LastFM
type LastfmAlbumData struct {
	Name      string
	URL       string
	ImageURL  string
	Artist    string
	ArtistURL string
	Scrobbles int
}

// LastFmGuildTopTracks contains the top tracks for a guild, it is built by the Worker and stored in redis
type LastFmGuildTopTracks struct {
	GuildID       string
	NumberOfUsers int
	Period        LastFmPeriod
	Tracks        []LastfmTrackData
	CachedAt      time.Time
}

// LastFmGuildTopTracksKey returns the redis key for LastFmGuildTopTracks
func LastFmGuildTopTracksKey(guildID string, period LastFmPeriod) (key string) {
	return "project-d:lastfm:guild-top-tracks:" + guildID + ":" + string(period)
}

// LastFmGetUserinfo returns information about a LastFM user
func LastFmGetUserinfo(ctx context.Context, lastfmUsername string) (userData LastfmUserData, err error) {
	// start tracing span
	var span opentracing.Span
	span, _ = opentracing.StartSpanFromContext(ctx, "dhelpers.LastFmGetUserinfo")
	defer span.Finish()

	// request data
	var lastfmUser lastfm.UserGetInfo
	lastfmUser, err = cache.GetLastFm().User.GetInfo(lastfm.P{"user": lastfmUsername})
	if err != nil {
		return userData, err
	}
	// parse fields into lastfmUserData
	userData.Username = lastfmUser.Name
	userData.Name = lastfmUser.RealName
	userData.Country = lastfmUser.Country
	if lastfmUser.PlayCount != "" {
		userData.Scrobbles, _ = strconv.Atoi(lastfmUser.PlayCount) // nolint: errcheck, gas
	}

	if len(lastfmUser.Images) > 0 {
		for _, image := range lastfmUser.Images {
			if image.Size == lastFmTargetImageSize {
				userData.Icon = image.Url
			}
		}
	}

	if lastfmUser.Registered.Unixtime != "" {
		timeI, err := strconv.ParseInt(lastfmUser.Registered.Unixtime, 10, 64)
		if err == nil {
			userData.AccountCreation = time.Unix(timeI, 0)
		}
	}

	return userData, nil
}

// LastFmGetRecentTracks returns recent tracks listened to by an user
func LastFmGetRecentTracks(ctx context.Context, lastfmUsername string, limit int) (tracksData []LastfmTrackData, err error) {
	// start tracing span
	var span opentracing.Span
	span, _ = opentracing.StartSpanFromContext(ctx, "dhelpers.LastFmGetRecentTracks")
	defer span.Finish()

	// request data
	var lastfmRecentTracks lastfm.UserGetRecentTracksExtended
	lastfmRecentTracks, err = cache.GetLastFm().User.GetRecentTracksExtended(lastfm.P{
		"limit": limit + 1, // in case nowplaying + already scrobbled
		"user":  lastfmUsername,
	})
	if err != nil {
		return nil, err
	}

	// parse fields
	if lastfmRecentTracks.Total > 0 {
		for i, track := range lastfmRecentTracks.Tracks {
			if i == 1 {
				// prevent nowplaying + already scrobbled
				if lastfmRecentTracks.Tracks[0].Url == track.Url {
					continue
				}
			}
			lastTrack := LastfmTrackData{
				Name:      track.Name,
				URL:       track.Url,
				Artist:    track.Artist.Name,
				ArtistURL: track.Artist.Url,
				Album:     track.Album.Name,
				Loved:     false,
			}
			for _, image := range track.Images {
				if image.Size == lastFmTargetImageSize {
					lastTrack.ImageURL = image.Url
				}
			}
			for _, image := range track.Artist.Image {
				if image.Size == lastFmTargetImageSize {
					lastTrack.ArtistImageURL = image.Url
				}
			}
			if track.Loved == "1" || track.Loved == "true" {
				lastTrack.Loved = true
			}
			if track.NowPlaying == "1" || track.NowPlaying == "true" {
				lastTrack.NowPlaying = true
			}

			timestamp, err := strconv.Atoi(track.Date.Uts)
			if err == nil {
				lastTrack.Time = time.Unix(int64(timestamp), 0)
			}

			tracksData = append(tracksData, lastTrack)
			if len(tracksData) >= limit {
				break
			}
		}
	}

	return tracksData, nil
}

// LastFmGetTopArtists returns the top artists of an user
func LastFmGetTopArtists(ctx context.Context, lastfmUsername string, limit int, period LastFmPeriod) (artistsData []LastfmArtistData, err error) {
	// start tracing span
	var span opentracing.Span
	span, _ = opentracing.StartSpanFromContext(ctx, "dhelpers.LastFmGetTopArtists")
	defer span.Finish()

	// request data
	var lastfmTopArtists lastfm.UserGetTopArtists
	lastfmTopArtists, err = cache.GetLastFm().User.GetTopArtists(lastfm.P{
		"limit":  limit,
		"user":   lastfmUsername,
		"period": string(period),
	})
	if err != nil {
		return nil, err
	}

	// parse fields
	if lastfmTopArtists.Total > 0 {
		for _, artist := range lastfmTopArtists.Artists {
			lastArtist := LastfmArtistData{
				Name: artist.Name,
				URL:  artist.Url,
			}
			for _, image := range artist.Images {
				if image.Size == lastFmTargetImageSize {
					lastArtist.ImageURL = image.Url
				}
			}
			lastArtist.Scrobbles, _ = strconv.Atoi(artist.PlayCount) // nolint: gas

			artistsData = append(artistsData, lastArtist)
			if len(artistsData) >= limit {
				break
			}
		}
	}

	return artistsData, nil
}

// LastFmGetTopTracks returns the top tracks of an user
func LastFmGetTopTracks(ctx context.Context, lastfmUsername string, limit int, period LastFmPeriod) (tracksData []LastfmTrackData, err error) {
	// start tracing span
	var span opentracing.Span
	span, _ = opentracing.StartSpanFromContext(ctx, "dhelpers.LastFmGetTopTracks")
	defer span.Finish()

	// request data
	var lastfmTopTracks lastfm.UserGetTopTracks
	lastfmTopTracks, err = cache.GetLastFm().User.GetTopTracks(lastfm.P{
		"limit":  limit,
		"user":   lastfmUsername,
		"period": string(period),
	})
	if err != nil {
		return nil, err
	}

	// parse fields
	if lastfmTopTracks.Total > 0 {
		for _, track := range lastfmTopTracks.Tracks {
			lastTrack := LastfmTrackData{
				Name:      track.Name,
				URL:       track.Url,
				Artist:    track.Artist.Name,
				ArtistURL: track.Artist.Url,
			}
			for _, image := range track.Images {
				if image.Size == lastFmTargetImageSize {
					lastTrack.ImageURL = image.Url
				}
			}
			lastTrack.Scrobbles, _ = strconv.Atoi(track.PlayCount) // nolint: gas

			tracksData = append(tracksData, lastTrack)
			if len(tracksData) >= limit {
				break
			}
		}
	}

	return tracksData, nil
}

// LastFmGetTopAlbums returns the top albums of an user
func LastFmGetTopAlbums(ctx context.Context, lastfmUsername string, limit int, period LastFmPeriod) (albumsData []LastfmAlbumData, err error) {
	// start tracing span
	var span opentracing.Span
	span, _ = opentracing.StartSpanFromContext(ctx, "dhelpers.LastFmGetTopAlbums")
	defer span.Finish()

	// request data
	var lastfmTopAlbums lastfm.UserGetTopAlbums
	lastfmTopAlbums, err = cache.GetLastFm().User.GetTopAlbums(lastfm.P{
		"limit":  limit,
		"user":   lastfmUsername,
		"period": string(period),
	})
	if err != nil {
		return nil, err
	}

	// parse fields
	if lastfmTopAlbums.Total > 0 {
		for _, album := range lastfmTopAlbums.Albums {
			lastAlbum := LastfmAlbumData{
				Name:      album.Name,
				URL:       album.Url,
				Artist:    album.Artist.Name,
				ArtistURL: album.Artist.Url,
			}
			for _, image := range album.Images {
				if image.Size == lastFmTargetImageSize {
					lastAlbum.ImageURL = image.Url
				}
			}
			lastAlbum.Scrobbles, _ = strconv.Atoi(album.PlayCount) // nolint: gas

			albumsData = append(albumsData, lastAlbum)
			if len(albumsData) >= limit {
				break
			}
		}
	}

	return albumsData, nil
}

// LastFmGetPeriodFromArgs parses args to figure out the correct period
func LastFmGetPeriodFromArgs(args []string) (period LastFmPeriod, newArgs []string) {
	for i, arg := range args {
		switch arg {
		case "7day", "7days", "week", "7", "seven":
			newArgs = append(args[:i], args[i+1:]...)
			return LastFmPeriod7day, newArgs
		case "1month", "month", "1", "one":
			newArgs = append(args[:i], args[i+1:]...)
			return LastFmPeriod1month, newArgs
		case "3month", "threemonths", "quarter", "3", "three":
			newArgs = append(args[:i], args[i+1:]...)
			return LastFmPeriod3month, newArgs
		case "6month", "halfyear", "half", "sixmonths", "6", "six":
			newArgs = append(args[:i], args[i+1:]...)
			return LastFmPeriod6month, newArgs
		case "12month", "year", "twelvemonths", "12", "twelve":
			newArgs = append(args[:i], args[i+1:]...)
			return LastFmPeriod12month, newArgs
		case "overall", "all", "alltime", "all-time": // nolint: misspell
			newArgs = append(args[:i], args[i+1:]...)
			return LastFmPeriodOverall, newArgs
		}
	}
	return LastFmPeriodOverall, args
}
