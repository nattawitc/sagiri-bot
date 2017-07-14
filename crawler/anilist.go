package crawler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"

	"github.com/spf13/viper"
)

const (
	anilistTokenURL       = "https://anilist.co/api/auth/access_token"
	anilistAnimeSearchURL = "https://anilist.co/api/anime/search/"
)

var (
	anilistClientID     string
	anilistClientSecret string
	AnilistClient       AnilistClientType
)

func init() {
	anilistClientID = viper.GetString("anilist-clientid")
	anilistClientSecret = viper.GetString("anilist-clientsecret")

	AnilistClient = AnilistClientType{}
	if err := AnilistClient.getToken(); err != nil {
		panic(err)
	}
}

type AnilistSeriesModel struct {
	ID             int      `json:"id"`
	SeriesType     string   `json:"series_type"`
	TitleRomaji    string   `json:"title_romaji"`
	TitleEnglish   string   `json:"title_english"`
	TitleJapanese  string   `json:"title_japanese"`
	Type           string   `json:"type"`
	StartDateFuzzy int      `json:"start_date_fuzzy"`
	EndDateFuzzy   int      `json:"end_date_fuzzy"`
	Season         int      `json:"season"`
	Description    string   `json:"description"`
	Synonyms       []string `json:"synonyms"`
	Genres         []string `json:"genres"`
	Adult          bool     `json:"adult"`
	AverageScore   float64  `json:"average_score"`
	Popularity     int      `json:"popularity"`
	ImageUrlSml    string   `json:"image_url_sml"`
	ImageUrlMed    string   `json:"image_url_med"`
	ImageUrlLge    string   `json:"image_url_lge"`
	ImageUrlBanner string   `json:"image_url_banner"`
	TotalEpisodes  int      `json:"total_episodes"`
	Duration       int      `json:"duration"`
	AiringStatus   string   `json:"airing_status"`
	Source         string   `json:"source"`
}

type AnilistClientType struct {
	token string
}

func (a *AnilistClientType) getToken() error {
	form := url.Values{}
	form.Set("grant_type", "client_credentials")
	form.Set("client_id", anilistClientID)
	form.Set("client_secret", anilistClientSecret)
	resp, err := http.PostForm(anilistTokenURL, form)
	if err != nil {
		return err
	}
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("anilist request token fail with status: %v", resp.Status)
	}
	defer resp.Body.Close()
	data := &struct {
		Token string `json:"access_token"`
	}{}
	err = json.NewDecoder(resp.Body).Decode(data)
	if err != nil {
		return err
	}
	a.token = data.Token
	if a.token == "" {
		return fmt.Errorf("can't find token in token response")
	}
	return nil
}

func (a *AnilistClientType) AnimeSearch(q string) ([]AnilistSeriesModel, error) {
	url := anilistAnimeSearchURL + q
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", a.token)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode == http.StatusUnauthorized {
		if err := a.getToken(); err != nil {
			return nil, err
		}
		req.Header.Set("Authorization", "Bearer "+a.token)
		resp, err = http.DefaultClient.Do(req)
		if err != nil {
			return nil, err
		}
		defer resp.Body.Close()
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("anime search with anilist failed with status: %v", resp.Status)
	}
	data := make([]AnilistSeriesModel, 0)
	err = json.NewDecoder(resp.Body).Decode(&data)
	if err != nil {
		return nil, err
	}
	return data, nil
}
