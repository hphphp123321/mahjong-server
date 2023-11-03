package robot

import (
	"context"
	"fmt"
	"github.com/hphphp123321/mahjong-go/mahjong"
	gprcRobot "github.com/hphphp123321/mahjong-server/app/api/v1/robot"
	model "github.com/hphphp123321/mahjong-server/app/service/robot/model"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type GrpcRobotType int

const (
	RobotEvents GrpcRobotType = iota
	RobotBoardState
)

type GrpcRobot struct {
	Name string
	Type GrpcRobotType

	client gprcRobot.GrpcRobotClient
}

func NewGrpcRobot(name string, robotType GrpcRobotType, ipAddr string, port int32) (*GrpcRobot, error) {
	tcpAddr := fmt.Sprintf("%s:%d", ipAddr, port)
	conn, err := grpc.Dial(tcpAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}
	client := gprcRobot.NewGrpcRobotClient(conn)
	return &GrpcRobot{
		Name:   name,
		Type:   robotType,
		client: client,
	}, nil
}

func (g GrpcRobot) GetRobotType() string {
	return g.Name
}

func (g GrpcRobot) ChooseAction(events mahjong.Events, validActions mahjong.Calls) (actionIdx int) {
	if g.client == nil {
		fmt.Println("g.client is nil")
		return 0
	}

	var err error
	switch g.Type {
	case RobotEvents:
		actionIdx, err = g.chooseActionByEvents(events, validActions)
	case RobotBoardState:
		actionIdx, err = g.chooseActionByBoardState(events, validActions)
	}

	if err != nil {
		fmt.Println(err)
		return 0
	}

	return actionIdx
}

func (g GrpcRobot) chooseActionByEvents(events mahjong.Events, actions mahjong.Calls) (int, error) {
	var es = model.ToPbEvents(events)
	var validActions = model.ToPbCalls(actions)
	var req = &gprcRobot.ChooseActionRequest{
		Events:       es,
		ValidActions: validActions,
	}
	var resp, err = g.client.ChooseAction(context.Background(), req)
	if err != nil {
		return 0, err
	}
	return int(resp.ActionIdx), nil
}

func (g GrpcRobot) chooseActionByBoardState(events mahjong.Events, actions mahjong.Calls) (int, error) {
	var boardState = mahjong.NewBoardState()
	boardState.DecodeEvents(events)
	boardState.ValidActions = actions
	var bs = model.ToPbBoardState(boardState)
	var res, err = g.client.ChooseActionByBoardState(context.Background(), bs)
	if err != nil {
		return 0, err
	}
	return int(res.ActionIdx), nil
}
