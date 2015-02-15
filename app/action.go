package main

import (
	"./db"
	"fmt"
	"math/rand"
	"net/url"
	"time"
)

const (
	_FOLLOW   = iota
	_RETWEET  = iota
	_FAVORITE = iota
	_TWEET    = iota
)

type Action struct {
	name   int
	weight int
}

func performAction() {
	actions := make([]Action, 0, 4)

	actions = append(actions, Action{name: _FOLLOW, weight: ACTION_FOLLOW_WEIGHT * rand.Intn(100)})
	actions = append(actions, Action{name: _RETWEET, weight: ACTION_RETWEET_WEIGHT * rand.Intn(100)})
	actions = append(actions, Action{name: _FAVORITE, weight: ACTION_FAVORITE_WEIGHT * rand.Intn(100)})
	actions = append(actions, Action{name: _TWEET, weight: ACTION_TWEET_WEIGHT * rand.Intn(100)})

	selectedAction := Action{name: -1, weight: -1}

	for _, action := range actions {
		if action.weight > selectedAction.weight {
			selectedAction = action
		}
	}

	switch selectedAction.name {
	case _FOLLOW:
		actionFollow()
		break
	case _RETWEET:
		actionRetweet()
		break
	case _FAVORITE:
		actionFavorite()
		break
	case _TWEET:
		actionTweet()
		break
	}
}

func actionFollow() {
	fmt.Println("Action follow")

	searchResult, err := api.GetSearch(KEYWORDS[rand.Intn(len(KEYWORDS))], nil)
	if err != nil {
		fmt.Println("Error while querying twitter API", err)
		return
	}

	for _, tweet := range searchResult.Statuses {

		follow, err := db.AlreadyFollow(tweet.User.Id)
		if err == nil && !follow {

			err := db.Follow{UserId: tweet.User.Id, UserName: tweet.User.ScreenName, Status: tweet.Text, FollowDate: time.Now()}.Persist()
			if err != nil {
				fmt.Println("Error while persisting follow", err)
				return
			}

			_, err = api.FollowUser(tweet.User.ScreenName)
			if err != nil {
				fmt.Println("Error while following user "+tweet.User.ScreenName+" : ", err)
			}

			fmt.Println("Now follow ", tweet.User.ScreenName)
			return
		}

	}
}

func actionRetweet() {
	fmt.Println("Action retweet")
}

func actionFavorite() {
	fmt.Println("Action fav")
}

func actionTweet() {
	fmt.Println("Action tweet")

	content, err := generateTweetContent()
	if err != nil {
		fmt.Println("Error while getting tweet content : ", err)
		return
	}

	tweetText := content.text + " " + content.url

	err = db.Tweet{Content: tweetText, Date: time.Now()}.Persist()
	if err != nil {
		fmt.Println("Error while persisting tweet", err)
		return
	}

	tweet, err := api.PostTweet(tweetText, url.Values{})
	if err != nil {
		fmt.Println("Error while posting tweet", err)
		return
	}

	fmt.Println("Tweet posted : ", tweet.Text)
}
