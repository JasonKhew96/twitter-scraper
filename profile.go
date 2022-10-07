package twitterscraper

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"sync"
	"time"
)

// Global cache for user IDs
var cacheIDs sync.Map

// Profile of twitter user.
type Profile struct {
	Avatar         string
	Banner         string
	Biography      string
	Birthday       string
	FollowersCount int
	FollowingCount int
	FriendsCount   int
	IsFollowing    bool
	IsPrivate      bool
	IsVerified     bool
	Joined         *time.Time
	LikesCount     int
	ListedCount    int
	Location       string
	Name           string
	PinnedTweetIDs []string
	TweetsCount    int
	URL            string
	UserID         string
	Username       string
	Website        string
}

type user struct {
	Result struct {
		TypeName string     `json:"__typename"`
		RestID   string     `json:"rest_id"`
		Legacy   legacyUser `json:"legacy"`
		Reason   string     `json:"reason"`
	} `json:"result"`
}

type userByScreenNameResp struct {
	Data struct {
		User *user `json:"user"`
	} `json:"data"`
}

type GetProfileVariables struct {
	ScreenName                 string `json:"screen_name"`
	WithSafetyModeUserFields   bool   `json:"withSafetyModeUserFields"`
	WithSuperFollowsUserFields bool   `json:"withSuperFollowsUserFields"`
}

type GetProfileFeatures struct {
	VerifiedPhoneLabelEnabled                     bool `json:"verified_phone_label_enabled"`
	ResponsiveWebGraphqlTimelineNavigationEnabled bool `json:"responsive_web_graphql_timeline_navigation_enabled"`
}

// GetProfile return parsed user profile.
func (s *Scraper) GetProfile(variables *GetProfileVariables, features *GetProfileFeatures) (*Profile, error) {
	if variables == nil {
		return nil, fmt.Errorf("variables is nil")
	}

	if features == nil {
		features = &GetProfileFeatures{}
	}

	jsonVariables, err := json.Marshal(variables)
	if err != nil {
		return nil, err
	}

	jsonFeatures, err := json.Marshal(features)
	if err != nil {
		return nil, err
	}

	queries := url.Values{}
	queries.Add("variables", string(jsonVariables))
	queries.Add("features", string(jsonFeatures))

	reqUrl, err := url.Parse("https://twitter.com/i/api/graphql/HThKoC4xtXHcuMIok4O0HA/UserByScreenName")
	if err != nil {
		return nil, err
	}

	reqUrl.RawQuery = queries.Encode()

	req, err := http.NewRequest("GET", reqUrl.String(), nil)
	if err != nil {
		return nil, err
	}

	var jsn userByScreenNameResp
	err = s.RequestAPI(req, &jsn)
	if err != nil {
		return nil, err
	}

	if jsn.Data.User == nil {
		return nil, fmt.Errorf("user '%s' not found", variables.ScreenName)
	}

	if jsn.Data.User.Result.TypeName != "User" {
		return nil, fmt.Errorf("%s", jsn.Data.User.Result.Reason)
	}

	if jsn.Data.User.Result.RestID == "" {
		return nil, fmt.Errorf("rest_id not found")
	}
	jsn.Data.User.Result.Legacy.IDStr = jsn.Data.User.Result.RestID

	if jsn.Data.User.Result.Legacy.ScreenName == "" {
		return nil, fmt.Errorf("either @%s does not exist or is private", variables.ScreenName)
	}

	profile := parseProfile(jsn.Data.User.Result.Legacy)

	return &profile, nil
}

// Deprecated: GetProfile wrapper for default scraper
func GetProfile(username string) (*Profile, error) {
	return defaultScraper.GetProfile(&GetProfileVariables{ScreenName: username}, nil)
}

// GetUserIDByScreenName from API
func (s *Scraper) GetUserIDByScreenName(screenName string) (string, error) {
	id, ok := cacheIDs.Load(screenName)
	if ok {
		return id.(string), nil
	}

	profile, err := s.GetProfile(&GetProfileVariables{ScreenName: screenName}, nil)
	if err != nil {
		return "", err
	}

	cacheIDs.Store(screenName, profile.UserID)

	return profile.UserID, nil
}
