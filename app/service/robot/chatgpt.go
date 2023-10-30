package robot

import (
	"context"
	"fmt"
	"github.com/hphphp123321/mahjong-go/mahjong"
	"github.com/sashabaranov/go-openai"
	"math/rand"
	"net/http"
	"net/url"
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
	var b = mahjong.NewBoardState()
	b.DecodeEvents(events)
	b.ValidActions = validActions
	return r.chatgptChooseAction(b)
}

func (r *ChatGPT) chatgptChooseAction(boardState *mahjong.BoardState) (actionIdx int) {
	var conf = openai.DefaultConfig(r.Key)

	var transport *http.Transport = nil

	if r.ProxyUrl == "env" {
		transport = &http.Transport{
			// 设置环境代理
			Proxy: http.ProxyFromEnvironment,
		}
	} else if r.ProxyUrl != "" {
		var proxy, err = url.Parse(r.ProxyUrl)
		if err != nil {
			fmt.Println(err)
			return rand.Intn(len(boardState.ValidActions))
		}
		transport = &http.Transport{
			// 设置代理
			Proxy: http.ProxyURL(proxy),
		}
	}

	conf.HTTPClient.Transport = transport
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
		fmt.Println(err)
		return rand.Intn(len(boardState.ValidActions))
	}
	fmt.Println(resp.Choices[0].Message.Content)
	fmt.Println("end chatgpt request")
	fmt.Println()

	return rand.Intn(len(boardState.ValidActions))
}

func (r *ChatGPT) getEnPrompt(boardState *mahjong.BoardState) string {
	var background = "Background: You're a Japanese Riichi mahjong pro, and you're playing one game.\n"
	var motivation = "Goal: You need to choose the optimal play based on the current situation\n"
	var board = fmt.Sprintf("Situation: Your self wind is %s, the wind round is %s(East2 means \"东二局\"), "+
		"the num of honba is %d, the num of riichi sticks is %d, the remaining tiles in the wall is %d, "+
		"the dora indicators are %s, the turn position(if position is not your self wind means you can call chi,pon,gang etc) is %s, "+
		"your hand tiles are %s;\n"+
		"the dealer/east's discarded tiles are %s, points is %d, riichi status is %t, melds are %s (none means he hasnot call chi pon or kan)\n"+
		"the south's discarded tiles are %s, points is %d, riichi status is %t, melds are %s\n"+
		"the west's discarded tiles are %s, points is %d, riichi status is %t, melds are %s\n"+
		"the north's discarded tiles are %s, points is %d, riichi status is %t, melds are %s\n\n"+
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
	var require = "Requirements: according to the legal action number, output the corresponding number, only reply one number of all choices, such as \"2\", DO NOT ANSWER OTHER EXCEPT NUMBER;"

	var prompt = background + motivation + board + validActions + require
	return prompt
}
