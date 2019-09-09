package twitter

import (
	"fmt"
	"net/url"
	"sagiri-bot/logger"
	"strconv"
	"sync"

	"github.com/ChimeraCoder/anaconda"
	"github.com/bwmarrin/discordgo"
	"github.com/spf13/viper"
)

var (
	defaultApi *anaconda.TwitterApi
	closeOnce  sync.Once
)

func getTwitterApi() *anaconda.TwitterApi {
	if defaultApi != nil {
		return defaultApi
	}
	consumerKey := viper.GetString("twitter-consumerkey")
	consumerSecret := viper.GetString("twitter-consumersecret")
	accesstoken := viper.GetString("twitter-accesstoken")
	accesstokenSecret := viper.GetString("twitter-accesstokensecret")

	anaconda.SetConsumerKey(consumerKey)
	anaconda.SetConsumerSecret(consumerSecret)
	defaultApi := anaconda.NewTwitterApi(accesstoken, accesstokenSecret)
	defaultApi.Log = anaconda.BasicLogger
	return defaultApi
}

func StartShiftcodeFollower(s *discordgo.Session, stop <-chan int, wg *sync.WaitGroup) {
	api := getTwitterApi()

	go func() {
		userID := "906234810"
		wg.Add(1)
		stream := api.PublicStreamFilter(url.Values{
			"follow": []string{userID},
		})
		defer stream.Stop()
		baseURL := "http://twitter.com/statuses/"
		channel := "453855461028921365"
		uid, err := strconv.ParseInt(userID, 10, 64)
		if err != nil {
			logger.PrintError("cannot parse "+userID+":", err)
			return
		}
		u, err := api.GetUsersShowById(uid, nil)
		if err != nil {
			logger.PrintError("cannot get user "+userID+":", err)
			return
		}
		_, err = s.ChannelMessageSend(channel, fmt.Sprintf("listening to twitter user %v(@%v)", u.Name, u.ScreenName))
		if err != nil {
			logger.PrintError("cannot send message to channel "+channel+":", err)
			return
		}
		for {
			select {
			case v := <-stream.C:
				t, ok := v.(anaconda.Tweet)
				if !ok {
					continue
				}
				fmt.Println(t.IdStr)

				if t.User.IdStr != userID {
					continue
				}

				if t.InReplyToStatusIdStr != "" {
					continue
				}
				//				if t.RetweetedStatus != nil {
				//					continue
				//				}

				url := baseURL + t.IdStr
				s.ChannelMessageSend(channel, url)
			case <-stop:
				closeOnce.Do(api.Close)
				wg.Done()
				break
			}
		}
	}()

}
