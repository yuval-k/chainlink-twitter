package twitter

import (
	"os"

	"github.com/dghubble/go-twitter/twitter"
	"github.com/dghubble/oauth1"
)

type TwitterClient interface {
	GetTweetsFor(user string) ([]string, error)
}

type twitterClient struct {
	Client *twitter.Client
}

func NewTwitterClientFromEnv() TwitterClient {
	var tc twitterClient
	consumerKey := os.Getenv("TWITTER_API_KEY")
	consumerSecret := os.Getenv("TWITTER_API_KEY_SECRET")
	accessToken := os.Getenv("TWITTER_ACCESS_TOKEN")
	accessTokenSecret := os.Getenv("TWITTER_ACCESS_TOKEN_SECRET")
	config := oauth1.NewConfig(consumerKey, consumerSecret)
	token := oauth1.NewToken(accessToken, accessTokenSecret)
	httpClient := config.Client(oauth1.NoContext, token)

	// Twitter client
	tc.Client = twitter.NewClient(httpClient)

	return &tc
}

func (tc *twitterClient) GetTweetsFor(user string) ([]string, error) {
	t := true
	tweets, _, err := tc.Client.Timelines.UserTimeline(&twitter.UserTimelineParams{
		ScreenName:     user,
		Count:          100,
		ExcludeReplies: &t,
		TrimUser:       &t,
	})
	if err != nil {
		return nil, err
	}

	var ret []string
	for _, t := range tweets {
		ret = append(ret, t.Text)
	}
	return ret, nil
}
