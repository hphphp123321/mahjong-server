package controller

import (
	"context"
	"github.com/hphphp123321/mahjong-go/mahjong"
	pb "github.com/hphphp123321/mahjong-server/app/api/v1"
	"github.com/hphphp123321/mahjong-server/app/errs"
	"github.com/hphphp123321/mahjong-server/app/global"
	"github.com/hphphp123321/mahjong-server/app/service/server"
	"io"
)

func AddGameStream(ctx context.Context, stream pb.Mahjong_GameServer, server *MahjongServer) error {
	cid, err := server.s.GetID(ctx)
	if err != nil {
		return err
	}
	server.clients[cid].gameStream = stream
	return nil
}

func RemoveGameStream(ctx context.Context, server *MahjongServer) error {
	cid, err := server.s.GetID(ctx)
	if err != nil {
		return err
	}
	server.clients[cid].gameStream = nil
	return nil
}

func BoardCastGameReply(ctx context.Context, server *MahjongServer, reply *pb.GameReply) error {
	r, err := server.s.ListPlayerIDs(ctx)
	if err != nil {
		return err
	}
	for _, pid := range r.PlayerIDs {
		c, ok := server.clients[pid]
		if !ok {
			return err
		}
		if c.gameStream == nil {
			continue
		}
		if err = c.gameStream.Send(reply); err != nil {
			return err
		}
	}
	return nil
}

func SendBackGame(ctx context.Context, server *MahjongServer, reply *pb.GameReply) error {
	cid, err := server.s.GetID(ctx)
	if err != nil {
		return err
	}
	c, ok := server.clients[cid]
	if !ok {
		return err
	}
	if c.gameStream == nil {
		return err
	}
	if err = c.gameStream.Send(reply); err != nil {
		return err
	}
	return nil
}

func StartGameSendStream(ctx context.Context, stream pb.Mahjong_GameServer, channels *server.StreamReply) (done chan error) {
	eventsChan := channels.Events
	validChan := channels.ValidCalls
	errChan := channels.Error
	done = make(chan error)
	go func() {
		for {
			select {
			case <-ctx.Done():
				global.Log.Infof("send game stream done")
			case err := <-errChan:
				if errS := SendGameEnd(stream); err != nil {
					global.Log.Warnf("send game end failed: %v", errS)
				}
				if err == errs.ErrGameEnd {
					global.Log.Infof("send game stream end")
					done <- nil
					return
				} else {
					global.Log.Warnf("send game stream error: %v", err)
					done <- nil
					return
				}
			case events := <-eventsChan:
				if err := SendEvents(stream, events); err != nil {
					global.Log.Warnf("send events failed: %v", err)
				}
			case validCalls := <-validChan:
				if err := SendValidActions(stream, validCalls); err != nil {
					global.Log.Warnf("send valid actions failed: %v", err)
				}
			}
		}
	}()
	return done
}

func SendGameEnd(stream pb.Mahjong_GameServer) error {
	end := true
	reply := &pb.GameReply{
		End: &end,
	}
	return stream.Send(reply)
}

func SendEvents(stream pb.Mahjong_GameServer, events mahjong.Events) error {
	pbEvents := ToPbEvents(events)
	reply := &pb.GameReply{
		Events: pbEvents,
	}
	return stream.Send(reply)
}

func SendValidActions(stream pb.Mahjong_GameServer, validCalls mahjong.Calls) error {
	pbValidCalls := ToPbCalls(validCalls)
	reply := &pb.GameReply{
		ValidActions: pbValidCalls,
	}
	return stream.Send(reply)
}

func StartGameRecvStream(ctx context.Context, stream pb.Mahjong_GameServer, server *MahjongServer) (done chan error, actionChan chan *mahjong.Call) {
	done = make(chan error)
	actionChan = make(chan *mahjong.Call, 1)
	go func() {
		for {
			in, err := stream.Recv()
			if err == io.EOF {
				global.Log.Infof("recv game stream EOF")
				done <- nil
				return
			}
			if err != nil {
				global.Log.Warnf("recv game stream failed: %v", err)
				done <- err
				return
			}
			switch in.GetRequest().(type) {
			case *pb.GameRequest_RefreshGame:
				if err := handleRefreshGame(ctx, server); err != nil {
					global.Log.Warnf("handle refresh game failed: %v", err)
				}
			case *pb.GameRequest_Chat:
				if err := handleGameChat(ctx, server, in); err != nil {
					global.Log.Warnf("handle chat failed: %v", err)
				}
			case *pb.GameRequest_Action:
				if err := handleGameAction(ctx, server, in, actionChan); err != nil {
					global.Log.Warnf("handle action failed: %v", err)
				}
			}
		}
	}()
	return done, actionChan
}

func handleGameAction(ctx context.Context, server *MahjongServer, in *pb.GameRequest, actionChan chan *mahjong.Call) error {
	cid, err := server.s.GetID(ctx)
	if err != nil {
		return err
	}
	c, ok := server.clients[cid]
	if !ok {
		return err
	}
	if c.gameStream == nil {
		return err
	}
	call := in.GetAction()
	mahjongCall := ToMahjongCall(call)
	actionChan <- mahjongCall
	return nil
}

func handleGameChat(ctx context.Context, server *MahjongServer, in *pb.GameRequest) error {
	cid, err := server.s.GetID(ctx)
	if err != nil {
		return err
	}
	c, ok := server.clients[cid]
	if !ok {
		return err
	}
	if c.gameStream == nil {
		return err
	}
	name, err := server.s.GetName(ctx)
	seat, err := server.s.GetSeat(ctx)
	reply := ToPbGameChatReply(in, name, seat)
	if err := BoardCastGameReply(ctx, server, reply); err != nil {
		return err
	}
	return nil
}

func handleRefreshGame(ctx context.Context, server *MahjongServer) (err error) {
	cid, err := server.s.GetID(ctx)
	if err != nil {
		return err
	}
	c, ok := server.clients[cid]
	if !ok {
		return err
	}
	if c.gameStream == nil {
		return err
	}
	b, err := server.s.GetBoardState(ctx)
	if err != nil {
		return err
	}
	reply := &pb.GameReply{
		Reply: &pb.GameReply_RefreshGameReply{
			RefreshGameReply: ToPbBoardState(b),
		},
	}
	if err := SendBackGame(ctx, server, reply); err != nil {
		return err
	}
	return nil
}
