package twitterscraper

import "time"

type (
	// Media type
	Media interface{}

	// MediaPhoto type
	MediaPhoto struct {
		Url string
		Alt string
	}

	// MediaVideo type
	MediaVideo struct {
		IsAnimatedGif bool
		Preview       string
		Url           string
		Alt           string
	}

	// Tweet type.
	Tweet struct {
		Hashtags         []string
		HTML             string
		ID               string
		InReplyToStatus  *Tweet
		IsQuoted         bool
		IsPin            bool
		IsReply          bool
		IsRetweet        bool
		IsRecommended    bool
		Likes            int
		Mentions         []string
		PermanentURL     string
		Place            *Place
		QuotedStatus     *Tweet
		Replies          int
		Retweets         int
		RetweetedStatus  *Tweet
		Text             string
		TimeParsed       time.Time
		Timestamp        int64
		URLs             []string
		UserID           string
		Username         string
		SensitiveContent bool
		Medias           []Media
	}

	// ProfileResult of scrapping.
	ProfileResult struct {
		Profile
		Error error
	}

	// TweetResult of scrapping.
	TweetResult struct {
		Tweet
		Error error
	}

	legacyUser struct {
		CreatedAt   string `json:"created_at"`
		Description string `json:"description"`
		Entities    struct {
			Description struct {
				Urls []struct {
					ExpandedURL string `json:"expanded_url"`
					URL         string `json:"url"`
				} `json:"urls"`
			} `json:"description"`
			URL struct {
				Urls []struct {
					ExpandedURL string `json:"expanded_url"`
				} `json:"urls"`
			} `json:"url"`
		} `json:"entities"`
		FavouritesCount      int      `json:"favourites_count"`
		FollowersCount       int      `json:"followers_count"`
		Following            bool     `json:"following"`
		FriendsCount         int      `json:"friends_count"`
		IDStr                string   `json:"id_str"`
		ListedCount          int      `json:"listed_count"`
		Name                 string   `json:"name"`
		Location             string   `json:"location"`
		PinnedTweetIdsStr    []string `json:"pinned_tweet_ids_str"`
		ProfileBannerURL     string   `json:"profile_banner_url"`
		ProfileImageURLHTTPS string   `json:"profile_image_url_https"`
		Protected            bool     `json:"protected"`
		ScreenName           string   `json:"screen_name"`
		StatusesCount        int      `json:"statuses_count"`
		Verified             bool     `json:"verified"`
	}

	Place struct {
		ID          string `json:"id"`
		PlaceType   string `json:"place_type"`
		Name        string `json:"name"`
		FullName    string `json:"full_name"`
		CountryCode string `json:"country_code"`
		Country     string `json:"country"`
		BoundingBox struct {
			Type        string        `json:"type"`
			Coordinates [][][]float64 `json:"coordinates"`
		} `json:"bounding_box"`
	}

	fetchProfileFunc func(query string, maxProfilesNbr int, cursor string) ([]*Profile, string, error)
	fetchTweetFunc   func(query string, maxTweetsNbr int, cursor string) ([]*Tweet, string, error)
)
