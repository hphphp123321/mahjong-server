package remote

import (
	"context"
	"github.com/hphphp123321/mahjong-go/mahjong"
	grpcRobot "github.com/hphphp123321/mahjong-server/app/api/v1/robot"
	model "github.com/hphphp123321/mahjong-server/app/service/robot/model"
)

type ActionChooser interface {
	ChooseAction(grpcClient grpcRobot.GrpcRobotClient, events mahjong.Events, validActions mahjong.Calls) (actionIdx int, err error)
}

type EventsChooser struct{}
type BoardStateChooser struct{}
type JSONChooser struct{}

func (c EventsChooser) ChooseAction(grpcClient grpcRobot.GrpcRobotClient, events mahjong.Events, actions mahjong.Calls) (int, error) {
	var es = model.ToPbEvents(events)
	var validActions = model.ToPbCalls(actions)
	var req = &grpcRobot.ChooseActionRequest{
		Events:       es,
		ValidActions: validActions,
	}
	var resp, err = grpcClient.ChooseAction(context.Background(), req)
	if err != nil {
		return 0, err
	}
	return int(resp.ActionIdx), nil
}

func (c BoardStateChooser) ChooseAction(grpcClient grpcRobot.GrpcRobotClient, events mahjong.Events, actions mahjong.Calls) (int, error) {
	var boardState = mahjong.NewBoardState()
	boardState.DecodeEvents(events)
	boardState.ValidActions = actions
	var bs = model.ToPbBoardState(boardState)
	var resp, err = grpcClient.ChooseActionByBoardState(context.Background(), bs)
	if err != nil {
		return 0, err
	}
	return int(resp.ActionIdx), nil
}

func (c JSONChooser) ChooseAction(grpcClient grpcRobot.GrpcRobotClient, events mahjong.Events, actions mahjong.Calls) (int, error) {
	var boardState = mahjong.NewBoardState()
	boardState.DecodeEvents(events)
	boardState.ValidActions = actions
	var j, err = boardState.MarshalJSON()
	if err != nil {
		return 0, err
	}
	var req = &grpcRobot.JsonMessage{
		Message: string(j),
	}
	resp, err := grpcClient.ChooseActionByJSON(context.Background(), req)
	if err != nil {
		return 0, err
	}
	return int(resp.ActionIdx), nil
}
