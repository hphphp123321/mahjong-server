package controller

import (
	"context"
	pb "github.com/hphphp123321/mahjong-server/app/api/v1"
	"github.com/hphphp123321/mahjong-server/app/errs"
	"io"
)

func AddReadyStream(ctx context.Context, stream pb.Mahjong_ReadyServer, server *MahjongServer) error {
	cid, err := server.s.GetID(ctx)
	if err != nil {
		return err
	}
	server.clients[cid].readyStream = stream
	return nil
}

func BoardCastReadyReply(ctx context.Context, server *MahjongServer, reply *pb.ReadyReply) error {
	r, err := server.s.ListPlayerIDs(ctx)
	if err != nil {
		return err
	}
	for _, pid := range r.PlayerIDs {
		c, ok := server.clients[pid]
		if !ok {
			return errs.ErrClientNotFound
		}
		if err = c.readyStream.Send(reply); err != nil {
			// TODO LOG
			continue
		}
	}
	return err
}

func StartReadyStream(ctx context.Context, stream pb.Mahjong_ReadyServer, server *MahjongServer) (done chan error, replyChan chan *pb.ReadyReply) {
	// TODO
	done = make(chan error)
	replyChan = make(chan *pb.ReadyReply)
	go func() {
		for {
			in, err := stream.Recv()
			if err == io.EOF {
				done <- nil
				return
			} else if err != nil {
				done <- err
				return
			}
			switch in.GetRequest().(type) {
			case *pb.ReadyRequest_Ping:
				reply, err := handlePing(ctx, server, in)
				sendChan(replyChan, done, reply, err)
			case *pb.ReadyRequest_GetReady:
				reply, err := handleGetReady(ctx, server, in)
				sendChan(replyChan, done, reply, err)
			case *pb.ReadyRequest_CancelReady:
				reply, err := handleCancelReady(ctx, server, in)
				sendChan(replyChan, done, reply, err)
			case *pb.ReadyRequest_LeaveRoom:
				reply, err := handleLeaveRoom(ctx, server, in)
				sendChan(replyChan, done, reply, err)
			case *pb.ReadyRequest_RemovePlayer:
				reply, err := handleRemovePlayer(ctx, server, in)
				sendChan(replyChan, done, reply, err)
			case *pb.ReadyRequest_AddRobot:
				reply, err := handleAddRobot(ctx, server, in)
				sendChan(replyChan, done, reply, err)
			case *pb.ReadyRequest_ListRobots:
				reply, err := handleListRobots(ctx, server, in)
				sendChan(replyChan, done, reply, err)
			case *pb.ReadyRequest_Chat:
				reply, err := handleChat(ctx, server, in)
				sendChan(replyChan, done, reply, err)
			case *pb.ReadyRequest_StartGame:
				reply, err := handleStartGame(ctx, server, in)
				sendChan(replyChan, done, reply, err)
			}
		}
	}()
	return done, replyChan
}

func handlePing(ctx context.Context, server *MahjongServer, in *pb.ReadyRequest) (reply *pb.ReadyReply, err error) {
	return &pb.ReadyReply{
		Message: "",
		Reply:   &pb.ReadyReply_Pong{Pong: &pb.Empty{}},
	}, nil
}

func handleGetReady(ctx context.Context, server *MahjongServer, in *pb.ReadyRequest) (reply *pb.ReadyReply, err error) {
	r, err := server.s.GetReady(ctx)
	if err != nil {
		return nil, err
	}
	return ToPbGetReadyReply(r), nil
}

func handleCancelReady(ctx context.Context, server *MahjongServer, in *pb.ReadyRequest) (reply *pb.ReadyReply, err error) {
	r, err := server.s.CancelReady(ctx)
	if err != nil {
		return nil, err
	}
	return ToPbCancelReadyReply(r), nil
}

func handleLeaveRoom(ctx context.Context, server *MahjongServer, in *pb.ReadyRequest) (reply *pb.ReadyReply, err error) {
	r, err := server.s.LeaveRoom(ctx)
	if err != nil {
		return nil, err
	}
	return ToPbLeaveRoomReply(r), nil
}

func handleRemovePlayer(ctx context.Context, server *MahjongServer, in *pb.ReadyRequest) (reply *pb.ReadyReply, err error) {
	r, err := server.s.RemovePlayer(ctx, ToServerRemovePlayerRequest(in))
	if err != nil {
		return nil, err
	}
	return ToPbRemovePlayerReply(r), nil
}

func handleAddRobot(ctx context.Context, server *MahjongServer, in *pb.ReadyRequest) (reply *pb.ReadyReply, err error) {
	r, err := server.s.AddRobot(ctx, ToServerAddRobotRequest(in))
	if err != nil {
		return nil, err
	}
	return ToPbAddRobotReply(r), nil
}

func handleListRobots(ctx context.Context, server *MahjongServer, in *pb.ReadyRequest) (reply *pb.ReadyReply, err error) {
	r, err := server.s.ListRobots(ctx)
	if err != nil {
		return nil, err
	}
	return ToPbListRobotsReply(r), nil
}

func handleChat(ctx context.Context, server *MahjongServer, in *pb.ReadyRequest) (reply *pb.ReadyReply, err error) {
	name, err := server.s.GetName(ctx)
	if err != nil {
		return nil, err
	}
	return ToPbChatReply(in, name), nil
}

func handleStartGame(ctx context.Context, server *MahjongServer, in *pb.ReadyRequest) (reply *pb.ReadyReply, err error) {
	return nil, nil
	// TODO
}

func sendChan(replyChan chan<- *pb.ReadyReply, done chan<- error, reply *pb.ReadyReply, err error) {
	if err != nil {
		done <- err
	} else {
		replyChan <- reply
	}
}
