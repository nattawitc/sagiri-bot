package cmdnav

import (
	"fmt"
	"net/http"
	"sync"
	"time"

	"sagiri-bot/logger"
	"sagiri-bot/utils"

	"github.com/bwmarrin/discordgo"
	"github.com/spf13/viper"
)

var (
	mutex       = &sync.Mutex{}
	commandList = make(map[string]*CommandResult)
	currentPage = make(map[string]*Page)
	delay       = 5 * time.Minute
)

func init() {
}

type Page struct {
	Message       func() *discordgo.MessageSend
	Reaction      map[string]*Page
	ReactionOrder []string
	cmdID         string
}

func NewPage(m func() *discordgo.MessageSend) *Page {
	return &Page{
		Message:       m,
		Reaction:      make(map[string]*Page),
		ReactionOrder: make([]string, 0),
	}
}

func NewPageWithMessage(m *discordgo.MessageSend) *Page {
	return &Page{
		Message: func() *discordgo.MessageSend {
			return m
		},
		Reaction:      make(map[string]*Page),
		ReactionOrder: make([]string, 0),
	}
}

func (p *Page) AddReaction(emoji string, page *Page) {
	p.Reaction[emoji] = page
	p.ReactionOrder = append(p.ReactionOrder, emoji)
}

type CommandResult struct {
	root  *Page
	Owner string
	timer *time.Timer
}

func (c *CommandResult) NewDeleteTimer(del func()) {
	c.timer = time.AfterFunc(5*time.Minute, del)
}

func (c *CommandResult) ResetDeleteTimer() {
	c.timer.Reset(delay)
}

func WritePage(s *discordgo.Session, channelID string, page *Page) (*discordgo.Message, error) {
	msg, err := s.ChannelMessageSendComplex(channelID, page.Message())
	if err != nil {
		return nil, err
	}
	for _, react := range page.ReactionOrder {
		url := fmt.Sprintf("https://discordapp.com/api/channels/%v/messages/%v/reactions/%v/@me", channelID, msg.ID, react)
		req, err := http.NewRequest(http.MethodPut, url, nil)
		req.Header.Add("Authorization", "Bot "+viper.GetString("token"))
		if err != nil {
			logger.PrintError(err)
			continue
		}
		_, err = http.DefaultClient.Do(req)
		if err != nil {
			logger.PrintError(err)
		}
		time.Sleep(50 * time.Millisecond)
	}
	return msg, nil
}

func RegisterCommand(root *Page, owner, messageID string) {
	id := utils.RandomString(6)
	for _, ok := commandList[id]; ok; {
		id = utils.RandomString(6)
	}
	commandResult := &CommandResult{
		root:  root,
		Owner: owner,
		timer: time.AfterFunc(delay, func() {
			mutex.Lock()
			delete(commandList, id)
			mutex.Unlock()
		}),
	}
	mutex.Lock()
	commandList[id] = commandResult
	root.cmdID = id
	currentPage[messageID] = root
	mutex.Unlock()
	time.AfterFunc(delay, func() {
		mutex.Lock()
		delete(currentPage, messageID)
		mutex.Unlock()
	})
}

func GetNextPage(messageID, emoji, userID string) *Page {
	mutex.Lock()
	defer mutex.Unlock()
	page, ok := currentPage[messageID]
	if !ok {
		return nil
	}
	cmd, ok := commandList[page.cmdID]
	if !ok {
		return nil
	}
	if cmd.Owner != userID {
		return nil
	}
	cmd.ResetDeleteTimer()

	nextPage, ok := page.Reaction[emoji]
	if !ok {
		return nil
	}

	delete(currentPage, messageID)
	nextPage.cmdID = page.cmdID
	return nextPage
}

func AddPage(messageID string, page *Page) {
	mutex.Lock()
	defer mutex.Unlock()
	currentPage[messageID] = page
	time.AfterFunc(delay, func() {
		mutex.Lock()
		delete(currentPage, messageID)
		mutex.Unlock()
	})
}
