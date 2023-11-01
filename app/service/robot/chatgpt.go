package robot

import (
	"context"
	"fmt"
	"math/rand"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"github.com/hphphp123321/mahjong-go/mahjong"
	"github.com/sashabaranov/go-openai"
)

// 确保实现了Robot接口
var _ Robot = (*ChatGPT)(nil)

type ChatGPT struct {
	Key      string
	Model    string
	Lang     string
	ProxyUrl string
}

func (r *ChatGPT) GetRobotType() string {
	return "chatgpt"
}

func (r *ChatGPT) ChooseAction(events mahjong.Events, validActions mahjong.Calls) (actionIdx int) {
	if len(validActions) == 1 {
		return 0
	}
	var b = mahjong.NewBoardState()
	b.DecodeEvents(events)
	b.ValidActions = validActions
	return r.chatgptChooseAction(b)
}

func (r *ChatGPT) chatgptChooseAction(boardState *mahjong.BoardState) (actionIdx int) {
	var conf = openai.DefaultConfig(r.Key)

	var transport *http.Transport = http.DefaultTransport.(*http.Transport)

	if r.ProxyUrl == "env" {
		transport.Proxy = http.ProxyFromEnvironment // 设置环境代理
	} else if r.ProxyUrl != "" {
		var proxy, err = url.Parse(r.ProxyUrl)
		if err != nil {
			fmt.Println(err)
			transport.Proxy = nil // 不设置代理
		} else {
			transport.Proxy = http.ProxyURL(proxy)
		}
	} else {
		transport.Proxy = nil // 不设置代理
	}

	conf.HTTPClient.Transport = transport
	conf.HTTPClient.Timeout = 10 * time.Second
	client := openai.NewClientWithConfig(conf)

	var prompt = r.getEnPrompt(boardState)

	fmt.Println()
	fmt.Println("start chatgpt request")
	fmt.Println(prompt)
	resp, err := client.CreateChatCompletion(
		context.Background(),
		openai.ChatCompletionRequest{
			Model:       r.Model,
			Temperature: 0,
			Messages: []openai.ChatCompletionMessage{
				{
					Role:    openai.ChatMessageRoleUser,
					Content: prompt,
				},
			},
		},
	)
	if err != nil {
		return r.chatgptReChooseAction(client, prompt, boardState)
	}
	var content = resp.Choices[0].Message.Content
	var choice int
	var choiceS string
	for _, b := range content {
		if b != ':' {
			choiceS += string(b)
		} else {
			break
		}
	}
	choice, err = strconv.Atoi(choiceS)
	if err != nil {
		fmt.Println(err)
		return r.chatgptReChooseAction(client, prompt, boardState)
	}
	if (choice < 0) || (choice >= len(boardState.ValidActions)) {
		return r.chatgptReChooseAction(client, prompt, boardState)
	}

	fmt.Println("Choice: " + boardState.ValidActions[choice].String())
	fmt.Println("end chatgpt request")
	fmt.Println()

	return choice
}

func (r *ChatGPT) chatgptReChooseAction(client *openai.Client, prompt string, boardState *mahjong.BoardState) (actionIdx int) {

	fmt.Println()
	fmt.Println("Restart chatgpt request")
	fmt.Println(prompt)
	resp, err := client.CreateChatCompletion(
		context.Background(),
		openai.ChatCompletionRequest{
			Model:       r.Model,
			Temperature: 0,
			Messages: []openai.ChatCompletionMessage{
				{
					Role:    openai.ChatMessageRoleUser,
					Content: prompt,
				},
			},
		},
	)
	if err != nil {
		fmt.Println(err)
		return rand.Intn(len(boardState.ValidActions))
	}
	var content = resp.Choices[0].Message.Content
	var choice int
	var choiceS string
	for _, b := range content {
		if b != ':' {
			choiceS += string(b)
		} else {
			break
		}
	}
	choice, err = strconv.Atoi(choiceS)
	if err != nil {
		fmt.Println(err)
		return rand.Intn(len(boardState.ValidActions))
	}
	if (choice < 0) || (choice >= len(boardState.ValidActions)) {
		return rand.Intn(len(boardState.ValidActions))
	}
	fmt.Println("Choice: " + boardState.ValidActions[choice].String())
	fmt.Println("End chatgpt request")
	fmt.Println()

	return choice
}

func (r *ChatGPT) getEnPrompt(boardState *mahjong.BoardState) string {
	var background = "Background: You're a Japanese Riichi mahjong pro, and you're playing one game. " +
		"The format of tiles are like \"Man1\" which means 1m or 一万, \"Sou3\" means 3s or 三索, \"Ton, Nan, Shaa, Pei\" means 东, 南, 西, 北, \"Haku, Hatsu, Chun\" means 白, 发, 中\n"
	var motivation = "Goal: You need to choose the optimal play based on the current situation\n"
	var board = fmt.Sprintf("Situation: Your self wind is %s, the wind round is %s(East2 means \"东二局\"), "+
		"the num of honba is %d, the num of riichi sticks is %d, the remaining tiles in the wall is %d, "+
		"the dora indicators are [%s], the turn position(if position is not your self wind means you can call chi,pon,gang etc) is %s, "+
		"your hand tiles are [%s];\n"+
		"the dealer/east's discarded tiles are [%s], points is %d, riichi status is %t, melds are [%s] (none means he hasnot call chi pon or kan)\n"+
		"the south's discarded tiles are [%s], points is %d, riichi status is %t, melds are [%s]\n"+
		"the west's discarded tiles are [%s], points is %d, riichi status is %t, melds are [%s]\n"+
		"the north's discarded tiles are [%s], points is %d, riichi status is %t, melds are [%s]\n\n"+
		"now it's your turn to choose one of the valid actions below:\n",

		boardState.PlayerWind.String(), boardState.WindRound.String(),
		boardState.NumHonba, boardState.NumRiichi, boardState.NumRemainTiles,
		boardState.DoraIndicators.Classes().String(), boardState.Position.String(),
		boardState.HandTiles.Classes().String(),
		boardState.PlayerStates[mahjong.East].DiscardTiles.Classes().String(),
		boardState.PlayerStates[mahjong.East].Points, boardState.PlayerStates[mahjong.East].IsRiichi,
		boardState.PlayerStates[mahjong.East].Melds.String(),
		boardState.PlayerStates[mahjong.South].DiscardTiles.Classes().String(),
		boardState.PlayerStates[mahjong.South].Points, boardState.PlayerStates[mahjong.South].IsRiichi,
		boardState.PlayerStates[mahjong.South].Melds.String(),
		boardState.PlayerStates[mahjong.West].DiscardTiles.Classes().String(),
		boardState.PlayerStates[mahjong.West].Points, boardState.PlayerStates[mahjong.West].IsRiichi,
		boardState.PlayerStates[mahjong.West].Melds.String(),
		boardState.PlayerStates[mahjong.North].DiscardTiles.Classes().String(),
		boardState.PlayerStates[mahjong.North].Points, boardState.PlayerStates[mahjong.North].IsRiichi,
		boardState.PlayerStates[mahjong.North].Melds.String(),
	)

	var validActions = ""
	for i, call := range boardState.ValidActions {
		validActions += fmt.Sprintf("%d: %s\n", i, call.String())
	}
	var require = "Requirements: according to the legal action number, output the corresponding number, only reply one number of all choices, such as \"2\", DO NOT ANSWER OTHER WORDS EXCEPT NUMBER!"

	var prompt = background + motivation + board + validActions + require
	return prompt
}
