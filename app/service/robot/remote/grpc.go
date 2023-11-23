package remote

import (
	"errors"
	"github.com/hphphp123321/mahjong-go/mahjong"
	gprcRobot "github.com/hphphp123321/mahjong-server/app/api/v1/robot"
	"github.com/hphphp123321/mahjong-server/app/service/robot"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type GrpcRobotType int

const (
	RobotEvents GrpcRobotType = iota
	RobotBoardState
	RobotJSON
)

var _ robot.Robot = (*GrpcRobot)(nil)

type GrpcRobot struct {
	Name          string
	actionChooser ActionChooser

	client gprcRobot.GrpcRobotClient
}

func NewGrpcRobot(name string, robotType GrpcRobotType, addr string) (*GrpcRobot, error) {
	conn, err := grpc.Dial(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}
	client := gprcRobot.NewGrpcRobotClient(conn)
	var actionChooser ActionChooser
	switch robotType {
	case RobotEvents:
		actionChooser = EventsChooser{}
	case RobotBoardState:
		actionChooser = BoardStateChooser{}
	case RobotJSON:
		actionChooser = JSONChooser{}
	}
	return &GrpcRobot{
		Name:          name,
		actionChooser: actionChooser,
		client:        client,
	}, nil
}

func (g GrpcRobot) GetRobotType() string {
	return g.Name
}

func (g GrpcRobot) ChooseAction(events mahjong.Events, validActions mahjong.Calls) (actionIdx int, err error) {
	if g.client == nil {
		return 0, errors.New("client is nil")
	}

	actionIdx, err = g.actionChooser.ChooseAction(g.client, events, validActions)
	if err != nil {
		return 0, err
	}

	return actionIdx, nil
}
