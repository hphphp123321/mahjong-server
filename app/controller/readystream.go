package controller

import (
	"context"
	pb "github.com/hphphp123321/mahjong-server/app/api/v1"
	"io"
)

func StartReadyStream(ctx context.Context, stream pb.Mahjong_ReadyServer, done chan error, replyChan chan *pb.ReadyReply) {
	// TODO
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
				reply, err := handlePing(ctx, in)
				sendChan(replyChan, done, reply, err)
			case *pb.ReadyRequest_GetReady:
				reply, err := handleGetReady(ctx, in)
				sendChan(replyChan, done, reply, err)
			case *pb.ReadyRequest_CancelReady:
				reply, err := handleCancelReady(ctx, in)
				sendChan(replyChan, done, reply, err)
			case *pb.ReadyRequest_LeaveRoom:
				reply, err := handleLeaveRoom(ctx, in)
				sendChan(replyChan, done, reply, err)
			case *pb.ReadyRequest_RemovePlayer:
				reply, err := handleRemovePlayer(ctx, in)
				sendChan(replyChan, done, reply, err)
			case *pb.ReadyRequest_AddRobot:
				reply, err := handleAddRobot(ctx, in)
				sendChan(replyChan, done, reply, err)
			case *pb.ReadyRequest_ListRobots:
				reply, err := handleListRobots(ctx, in)
				sendChan(replyChan, done, reply, err)
			case *pb.ReadyRequest_Chat:
				reply, err := handleChat(ctx, in)
				sendChan(replyChan, done, reply, err)
			case *pb.ReadyRequest_StartGame:
				reply, err := handleStartGame(ctx, in)
				sendChan(replyChan, done, reply, err)
			}
		}
	}()
}

func handlePing(ctx context.Context, in *pb.ReadyRequest) (reply *pb.ReadyReply, err error) {
	return nil, nil
	// TODO
}

func handleGetReady(ctx context.Context, in *pb.ReadyRequest) (reply *pb.ReadyReply, err error) {
	return nil, nil
	// TODO
}

func handleCancelReady(ctx context.Context, in *pb.ReadyRequest) (reply *pb.ReadyReply, err error) {
	return nil, nil
	// TODO
}

func handleLeaveRoom(ctx context.Context, in *pb.ReadyRequest) (reply *pb.ReadyReply, err error) {
	return nil, nil
	// TODO
}

func handleRemovePlayer(ctx context.Context, in *pb.ReadyRequest) (reply *pb.ReadyReply, err error) {
	return nil, nil
	// TODO
}

func handleAddRobot(ctx context.Context, in *pb.ReadyRequest) (reply *pb.ReadyReply, err error) {
	return nil, nil
	// TODO
}

func handleListRobots(ctx context.Context, in *pb.ReadyRequest) (reply *pb.ReadyReply, err error) {
	return nil, nil
	// TODO
}

func handleChat(ctx context.Context, in *pb.ReadyRequest) (reply *pb.ReadyReply, err error) {
	return nil, nil
	// TODO
}

func handleStartGame(ctx context.Context, in *pb.ReadyRequest) (reply *pb.ReadyReply, err error) {
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
