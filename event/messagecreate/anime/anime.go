package anime

import (
	"fmt"
	"sagiri-bot/cmdnav"
	"sagiri-bot/crawler"
	"sagiri-bot/emoji"
	"sagiri-bot/event/messagecreate"
	"sagiri-bot/logger"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
)

var (
	seasonName  = []string{"Winter", "Spring", "Summer", "Fall"}
	itemPerPage = 5
)

func init() {
	messagecreate.AddCommand("anime", "search for anime info", Anime)
}

func Anime(s *discordgo.Session, m *discordgo.MessageCreate) {
	query := strings.TrimSpace(strings.TrimPrefix(m.Content, "!anime"))
	animes, err := crawler.AnilistClient.AnimeSearch(query)
	if err != nil {
		logger.PrintError(err)
		return
	}

	index := []string{emoji.One, emoji.Two, emoji.Three, emoji.Four, emoji.Five}
	menu := &discordgo.MessageSend{}
	menu.Embed = &discordgo.MessageEmbed{}
	menuPage := cmdnav.NewPageWithMessage(menu)
	firstpage := menuPage
	counter := 0
	totalPage := len(animes)/itemPerPage + 1
	currentPage := 1
	menu.Embed.Description += fmt.Sprintf("**%d/%d** *navigate using button below*\n\n", currentPage, totalPage)
	for _, anime := range animes {
		if counter == 5 {
			menu = &discordgo.MessageSend{}
			menu.Embed = &discordgo.MessageEmbed{}
			currentPage += 1
			menu.Embed.Description += fmt.Sprintf("**%d/%d** *navigate using button below*\n\n", currentPage, totalPage)
			newMenuPage := cmdnav.NewPageWithMessage(menu)
			menuPage.AddReaction(emoji.ArrowRight, newMenuPage)
			newMenuPage.AddReaction(emoji.ArrowLeft, menuPage)
			menuPage = newMenuPage
			counter = 0
		}
		menu.Embed.Description += fmt.Sprintf("%v **%v(%v)**\n", index[counter], anime.TitleEnglish, anime.TitleRomaji)
		animePage := cmdnav.NewPageWithMessage(generateAnimeMessage(anime))
		menuPage.AddReaction(index[counter], animePage)
		animePage.AddReaction(emoji.LeftwardsArrowWithHook, menuPage)
		counter += 1
	}

	msg, err := cmdnav.WritePage(s, m.ChannelID, firstpage)
	if err != nil {
		logger.PrintError(err)
		return
	}

	cmdnav.RegisterCommand(firstpage, m.Author.ID, msg.ID)
}
func generateAnimeMessage(anime crawler.AnilistSeriesModel) *discordgo.MessageSend {
	animeUrl := fmt.Sprintf("https://anilist.co/anime/%v", anime.ID)
	embed := &discordgo.MessageEmbed{}
	embed.Description += fmt.Sprintf("*navigate using button below*\n")
	embed.Thumbnail = &discordgo.MessageEmbedThumbnail{URL: anime.ImageUrlMed}
	embed.Fields = make([]*discordgo.MessageEmbedField, 0) //{title, description, Type, status, start, season, source, duration, genre}

	title := &discordgo.MessageEmbedField{
		Name:   "Title",
		Inline: false,
	}
	title.Value = "[**" + anime.TitleEnglish
	if anime.TitleEnglish != anime.TitleRomaji {
		title.Value += "\n" + anime.TitleRomaji
	}
	if anime.TitleEnglish != anime.TitleJapanese {
		title.Value += "\n" + anime.TitleJapanese
	}
	title.Value += "**](" + animeUrl + ")"
	embed.Fields = append(embed.Fields, title)

	if anime.Description != "" {
		description := &discordgo.MessageEmbedField{
			Name:   "Description",
			Value:  anime.Description,
			Inline: false,
		}
		embed.Fields = append(embed.Fields, description)
	}

	if anime.Type != "" {
		Type := &discordgo.MessageEmbedField{
			Name:   "Type",
			Value:  anime.Type,
			Inline: true,
		}
		embed.Fields = append(embed.Fields, Type)
	}

	if anime.AiringStatus != "" {
		status := &discordgo.MessageEmbedField{
			Name:   "Status",
			Value:  anime.AiringStatus,
			Inline: true,
		}
		embed.Fields = append(embed.Fields, status)
	}

	if anime.StartDateFuzzy != 0 {
		year := anime.StartDateFuzzy / 10000
		month := time.Month((anime.StartDateFuzzy % 10000) / 100)
		day := anime.StartDateFuzzy % 100
		start := &discordgo.MessageEmbedField{
			Name:   "Start",
			Value:  fmt.Sprintf("%v %v %v", day, month, year),
			Inline: true,
		}
		embed.Fields = append(embed.Fields, start)
	}

	if anime.Season != 0 {
		syear := anime.Season/10 + 2000
		ss := seasonName[anime.Season%10-1]
		season := &discordgo.MessageEmbedField{
			Name:   "Season",
			Value:  fmt.Sprintf("%v %v", ss, syear),
			Inline: true,
		}
		embed.Fields = append(embed.Fields, season)
	}

	if anime.Source != "" {
		source := &discordgo.MessageEmbedField{
			Name:   "Source",
			Value:  anime.Source,
			Inline: true,
		}
		embed.Fields = append(embed.Fields, source)
	}

	if anime.Duration != 0 {
		duration := &discordgo.MessageEmbedField{
			Name:   "Duration",
			Value:  fmt.Sprintf("%v min", anime.Duration),
			Inline: true,
		}
		embed.Fields = append(embed.Fields, duration)
	}

	if len(anime.Genres) != 0 {
		genre := &discordgo.MessageEmbedField{
			Name:   "Genre",
			Value:  strings.Join(anime.Genres, ","),
			Inline: false,
		}
		embed.Fields = append(embed.Fields, genre)
	}

	return &discordgo.MessageSend{Embed: embed}
}
