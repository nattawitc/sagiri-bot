package smash

import (
	"image"
	"image/color"
	"image/draw"
	"io"
	"net/http"
	"os"

	"sagiri-bot/event/messagecreate"
	"sagiri-bot/logger"

	"github.com/bwmarrin/discordgo"
	"github.com/disintegration/gift"
	"github.com/fogleman/gg"
)

var (
	throwTemplate draw.Image
)

func init() {
	loadTemplate()
	messagecreate.AddCommand("smash", "smash avatar of user who's mentioned", Smash)
}

func loadTemplate() {
	f, err := os.Open("resource/slam_template.png")
	if err != nil {
		panic(err)
	}
	img, _, err := image.Decode(f)
	if err != nil {
		panic(err)
	}
	width := img.Bounds().Dx()
	g := gift.New(
		gift.Resize(width*2, 0, gift.LanczosResampling),
	)
	throwTemplate = image.NewRGBA(g.Bounds(img.Bounds()))
	g.Draw(throwTemplate, img)
}

func Smash(s *discordgo.Session, m *discordgo.MessageCreate) {
	//  filename := "58724194_p02.jpg"
	for _, user := range m.Mentions {
		url := "https://cdn.discordapp.com/avatars/" + user.ID + "/" + user.Avatar + ".png"
		resp, err := http.Get(url)
		if err != nil {
			logger.PrintError(err)
			return
		}
		if resp.StatusCode != http.StatusOK {
			s.ChannelMessageSendComplex(m.ChannelID, &discordgo.MessageSend{
				Content: "I can't get the avatar!!",
				Embed: &discordgo.MessageEmbed{
					Image: &discordgo.MessageEmbedImage{
						URL: "https://media.giphy.com/media/l0IyhP3coXnaEi9os/giphy.gif",
					},
				},
			})
			logger.PrintError(resp.Status)
			return
		}
		defer resp.Body.Close()

		avatar, _, err := image.Decode(resp.Body)
		g := gift.New(
			gift.Resize(0, 150, gift.LanczosResampling),
			gift.Rotate(55, color.Transparent, gift.CubicInterpolation),
		)
		dst := image.NewRGBA(g.Bounds(avatar.Bounds()))
		g.Draw(dst, avatar)

		ctx := gg.NewContextForImage(throwTemplate)
		ctx.DrawImageAnchored(dst, 140, 330, 0.5, 0.5)

		r, w := io.Pipe()
		go func() {
			ctx.EncodePNG(w)
			w.Close()
		}()
		data := &discordgo.MessageSend{
			File: &discordgo.File{
				Name:   user.Avatar + ".png",
				Reader: r,
			},
		}
		//      s.ChannelMessageSend(channelID, "Avatar of "+user.Username)
		s.ChannelMessageSendComplex(m.ChannelID, data)
		if err != nil {
			logger.PrintError(err)
		}
	}
}
