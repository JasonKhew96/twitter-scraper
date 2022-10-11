package twitterscraper_test

import (
	"context"
	"testing"
	"time"

	twitterscraper "github.com/JasonKhew96/twitter-scraper"
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
)

var cmpOptions = cmp.Options{
	cmpopts.IgnoreFields(twitterscraper.Tweet{}, "Likes"),
	cmpopts.IgnoreFields(twitterscraper.Tweet{}, "Replies"),
	cmpopts.IgnoreFields(twitterscraper.Tweet{}, "Retweets"),
}

func TestGetTweets(t *testing.T) {
	count := 0
	maxTweetsNbr := 200
	dupcheck := make(map[string]bool)
	scraper := twitterscraper.New()
	for tweet := range scraper.GetTweets(context.Background(), "Twitter", maxTweetsNbr) {
		if tweet.Error != nil {
			t.Error(tweet.Error)
		} else {
			count++
			if tweet.ID == "" {
				t.Error("Expected tweet ID is empty")
			} else {
				if dupcheck[tweet.ID] {
					t.Errorf("Detect duplicated tweet ID: %s", tweet.ID)
				} else {
					dupcheck[tweet.ID] = true
				}
			}
			if tweet.UserID == "" {
				t.Error("Expected tweet UserID is empty")
			}
			if tweet.Username == "" {
				t.Error("Expected tweet Username is empty")
			}
			if tweet.PermanentURL == "" {
				t.Error("Expected tweet PermanentURL is empty")
			}
			if tweet.Text == "" {
				t.Error("Expected tweet Text is empty")
			}
			if tweet.TimeParsed.IsZero() {
				t.Error("Expected tweet TimeParsed is zero")
			}
			if tweet.Timestamp == 0 {
				t.Error("Expected tweet Timestamp is greater than zero")
			}
			for _, media := range tweet.Medias {
				switch v := media.(type) {
				case twitterscraper.MediaVideo:
					if v.IsAnimatedGif {
						t.Error("Expected video is not animated gif")
					}
					if v.Preview == "" {
						t.Error("Expected video preview is empty")
					}
					if v.Url == "" {
						t.Error("Expected video URL is empty")
					}
				}
			}
		}
	}
	if count != maxTweetsNbr {
		t.Errorf("Expected tweets count=%v, got: %v", maxTweetsNbr, count)
	}
}

func TestGetTweet(t *testing.T) {
	sample := twitterscraper.Tweet{
		HTML:         "whoa, it works<br><br>now everyone can mix GIFs, videos, and images in one Tweet, available on iOS and Android <br><a href=\"https://t.co/LVVolAQPZi\"><img src=\"https://pbs.twimg.com/tweet_video_thumb/FeU5fh1XkA0vDAE.jpg\"/></a><br><img src=\"https://pbs.twimg.com/media/FeU5fhPXkCoZXZB.jpg\"/>",
		ID:           "1577730467436138524",
		PermanentURL: "https://twitter.com/Twitter/status/1577730467436138524",
		Text:         "whoa, it works\n\nnow everyone can mix GIFs, videos, and images in one Tweet, available on iOS and Android",
		TimeParsed:   time.Date(2022, 10, 05, 18, 40, 30, 0, time.FixedZone("UTC", 0)),
		Timestamp:    1664995230,
		UserID:       "783214",
		Username:     "Twitter",
		Medias: []twitterscraper.Media{
			twitterscraper.MediaVideo{
				IsAnimatedGif: true,
				Preview:       "https://pbs.twimg.com/tweet_video_thumb/FeU5fh1XkA0vDAE.jpg",
				Url:           "https://video.twimg.com/tweet_video/FeU5fh1XkA0vDAE.mp4",
			},
			twitterscraper.MediaPhoto{
				Url: "https://pbs.twimg.com/media/FeU5fhPXkCoZXZB.jpg",
				Alt: "picture of Kermit doing a one legged stand on a bicycle seat riding through the park",
			},
		},
	}
	scraper := twitterscraper.New()
	tweet, err := scraper.GetTweet("1577730467436138524")
	if err != nil {
		t.Error(err)
	} else {
		if diff := cmp.Diff(sample, *tweet, cmpOptions...); diff != "" {
			t.Error("Resulting tweet does not match the sample", diff)
		}
	}
}

func TestQuotedAndReply(t *testing.T) {
	sample := &twitterscraper.Tweet{
		HTML:         "The Easiest Problem Everyone Gets Wrong <br><br>[new video] --&gt; <a href=\"https://youtu.be/ytfCdqWhmdg\">https://t.co/YdaeDYmPAU</a> <br><a href=\"https://t.co/iKu4Xs6o2V\"><img src=\"https://pbs.twimg.com/media/ESsZa9AXgAIAYnF.jpg\"/></a>",
		ID:           "1237110546383724547",
		Likes:        485,
		PermanentURL: "https://twitter.com/VsauceTwo/status/1237110546383724547",
		// Photos:       []string{"https://pbs.twimg.com/media/ESsZa9AXgAIAYnF.jpg"},
		Replies:    12,
		Retweets:   18,
		Text:       "The Easiest Problem Everyone Gets Wrong \n\n[new video] --> https://youtu.be/ytfCdqWhmdg",
		TimeParsed: time.Date(2020, 03, 9, 20, 18, 33, 0, time.FixedZone("UTC", 0)),
		Timestamp:  1583785113,
		URLs:       []string{"https://youtu.be/ytfCdqWhmdg"},
		UserID:     "978944851",
		Username:   "VsauceTwo",
		Medias: []twitterscraper.Media{twitterscraper.MediaPhoto{
			Url: "https://pbs.twimg.com/media/ESsZa9AXgAIAYnF.jpg",
		}},
	}
	scraper := twitterscraper.New()
	tweet, err := scraper.GetTweet("1237110897597976576")
	if err != nil {
		t.Error(err)
	} else {
		if !tweet.IsQuoted {
			t.Error("IsQuoted must be True")
		}
		if diff := cmp.Diff(sample, tweet.QuotedStatus, cmpOptions...); diff != "" {
			t.Error("Resulting quote does not match the sample", diff)
		}
	}
	tweet, err = scraper.GetTweet("1237111868445134850")
	if err != nil {
		t.Error(err)
	} else {
		if !tweet.IsReply {
			t.Error("IsReply must be True")
		}
		if diff := cmp.Diff(sample, tweet.InReplyToStatus, cmpOptions...); diff != "" {
			t.Error("Resulting reply does not match the sample", diff)
		}
	}

}
func TestRetweet(t *testing.T) {
	sample := &twitterscraper.Tweet{
		HTML:         "We’ve seen an increase in attacks against Asian communities and individuals around the world. It’s important to know that this isn’t new; throughout history, Asians have experienced violence and exclusion. However, their diverse lived experiences have largely been overlooked.",
		ID:           "1359151057872580612",
		Likes:        6683,
		PermanentURL: "https://twitter.com/TwitterTogether/status/1359151057872580612",
		Replies:      456,
		Retweets:     1495,
		Text:         "We’ve seen an increase in attacks against Asian communities and individuals around the world. It’s important to know that this isn’t new; throughout history, Asians have experienced violence and exclusion. However, their diverse lived experiences have largely been overlooked.",
		TimeParsed:   time.Date(2021, 02, 9, 14, 43, 58, 0, time.FixedZone("UTC", 0)),
		Timestamp:    1612881838,
		UserID:       "773578328498372608",
		Username:     "TwitterTogether",
	}
	scraper := twitterscraper.New()
	tweet, err := scraper.GetTweet("1362849141248974853")
	if err != nil {
		t.Error(err)
	} else {
		if !tweet.IsRetweet {
			t.Error("IsRetweet must be True")
		}
		if diff := cmp.Diff(sample, tweet.RetweetedStatus, cmpOptions...); diff != "" {
			t.Error("Resulting retweet does not match the sample", diff)
		}
	}
}
