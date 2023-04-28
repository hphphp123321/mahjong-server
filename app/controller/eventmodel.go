package controller

import (
	"github.com/hphphp123321/go-common"
	"github.com/hphphp123321/mahjong-go/mahjong"
	pb "github.com/hphphp123321/mahjong-server/app/api/v1"
)

func ToPbEventStart(event *mahjong.EventStart) *pb.EventStart {
	return &pb.EventStart{
		WindRound:         ToPbWindRound(event.WindRound),
		InitWind:          ToPbWind(event.InitWind),
		Seed:              event.Seed,
		NumGame:           int32(event.NumGame),
		NumHonba:          int32(event.NumHonba),
		NumRiichi:         int32(event.NumRiichi),
		InitDoraIndicator: ToPbTile(event.InitDoraIndicator),
		InitTiles:         ToPbTiles(event.InitTiles),
		GameRule:          ToPbGameRule(event.Rule),
		PlayersPoints: map[int32]int32{
			0: int32(event.PlayersPoints[mahjong.East]),
			1: int32(event.PlayersPoints[mahjong.South]),
			2: int32(event.PlayersPoints[mahjong.West]),
			3: int32(event.PlayersPoints[mahjong.North]),
		},
	}
}

func ToPbEventGet(event *mahjong.EventGet) *pb.EventGet {
	return &pb.EventGet{
		Who:  ToPbWind(event.Who),
		Tile: ToPbTile(event.Tile),
	}
}

func ToPbEventDiscard(event *mahjong.EventDiscard) *pb.EventDiscard {
	return &pb.EventDiscard{
		Who:         ToPbWind(event.Who),
		Tile:        ToPbTile(event.Tile),
		TsumoGiri:   false,
		TenhaiSlice: ToPbTileClasses(event.TenhaiSlice),
	}
}

func ToPbEventTsumoGiri(event *mahjong.EventTsumoGiri) *pb.EventDiscard {
	return &pb.EventDiscard{
		Who:         ToPbWind(event.Who),
		Tile:        ToPbTile(event.Tile),
		TsumoGiri:   true,
		TenhaiSlice: ToPbTileClasses(event.TenhaiSlice),
	}
}

func ToPbEventChi(event *mahjong.EventChi) *pb.EventCall {
	return &pb.EventCall{
		Who:  ToPbWind(event.Who),
		Call: ToPbCall(event.Call),
	}
}

func ToPbEventPon(event *mahjong.EventPon) *pb.EventCall {
	return &pb.EventCall{
		Who:  ToPbWind(event.Who),
		Call: ToPbCall(event.Call),
	}
}

func ToPbEventDaiMinKan(event *mahjong.EventDaiMinKan) *pb.EventCall {
	return &pb.EventCall{
		Who:  ToPbWind(event.Who),
		Call: ToPbCall(event.Call),
	}
}

func ToPbEventShouMinKan(event *mahjong.EventShouMinKan) *pb.EventCall {
	return &pb.EventCall{
		Who:  ToPbWind(event.Who),
		Call: ToPbCall(event.Call),
	}
}

func ToPbEventAnKan(event *mahjong.EventAnKan) *pb.EventCall {
	return &pb.EventCall{
		Who:  ToPbWind(event.Who),
		Call: ToPbCall(event.Call),
	}
}

func ToPbEventRiichi(event *mahjong.EventRiichi) *pb.EventRiichi {
	return &pb.EventRiichi{
		Who:  ToPbWind(event.Who),
		Step: int32(event.Step),
	}
}

func ToPbEventRon(event *mahjong.EventRon) *pb.EventRon {
	return &pb.EventRon{
		Who:       ToPbWind(event.Who),
		FromWho:   ToPbWind(event.FromWho),
		HandTiles: ToPbTiles(event.HandTiles),
		WinTile:   ToPbTile(event.WinTile),
		Result:    ToPbResult(event.Result),
	}
}

func ToPbEventTsumo(event *mahjong.EventTsumo) *pb.EventTsumo {
	return &pb.EventTsumo{
		Who:       ToPbWind(event.Who),
		HandTiles: ToPbTiles(event.HandTiles),
		WinTile:   ToPbTile(event.WinTile),
		Result:    ToPbResult(event.Result),
	}
}

func ToPbEventNewIndicator(event *mahjong.EventNewIndicator) *pb.EventNewIndicator {
	return &pb.EventNewIndicator{
		Tile: ToPbTile(event.Tile),
	}
}

func ToPbEventChanKan(event *mahjong.EventChanKan) *pb.EventRon {
	return &pb.EventRon{
		Who:       ToPbWind(event.Who),
		FromWho:   ToPbWind(event.FromWho),
		HandTiles: ToPbTiles(event.HandTiles),
		WinTile:   ToPbTile(event.WinTile),
		Result:    ToPbResult(event.Result),
	}
}

func ToPbEventRyuuKyoku(event *mahjong.EventRyuuKyoku) *pb.EventRyuuKyoKu {
	return &pb.EventRyuuKyoKu{
		Who:       ToPbWind(event.Who),
		HandTiles: ToPbTiles(event.HandTiles),
		Reason:    pb.RyuuKyoKuReason(event.Reason),
	}
}

func ToPbEventEnd(event *mahjong.EventEnd) *pb.EventEnd {
	return &pb.EventEnd{
		EastPointsChange:  int32(event.PointsChange[mahjong.East]),
		SouthPointsChange: int32(event.PointsChange[mahjong.South]),
		WestPointsChange:  int32(event.PointsChange[mahjong.West]),
		NorthPointsChange: int32(event.PointsChange[mahjong.North]),
	}
}

func ToPbEventFuriten(event *mahjong.EventFuriten) *pb.EventFuriten {
	return &pb.EventFuriten{
		Who:    ToPbWind(event.Who),
		Reason: pb.FuritenReason(event.FuritenReason),
	}
}

func ToPbNagashiMangan(event *mahjong.EventNagashiMangan) *pb.EventNagashiMangan {
	return &pb.EventNagashiMangan{
		Who: ToPbWind(event.Who),
	}
}

func ToPbTenhaiEnd(event *mahjong.EventTenhaiEnd) *pb.EventTenhaiEnd {
	return &pb.EventTenhaiEnd{
		Who:         ToPbWind(event.Who),
		HandTiles:   ToPbTiles(event.HandTiles),
		TenhaiSlice: ToPbTileClasses(event.TenhaiSlice),
	}
}

func ToPbEvents(events mahjong.Events) []*pb.Event {
	return common.MapSlice(events, ToPbEvent)
}

func ToPbEvent(event mahjong.Event) *pb.Event {
	switch event.GetType() {
	case mahjong.EventTypeStart:
		return &pb.Event{
			Event: &pb.Event_EventStart{
				EventStart: ToPbEventStart(event.(*mahjong.EventStart)),
			}}
	case mahjong.EventTypeGet:
		return &pb.Event{
			Event: &pb.Event_EventGet{
				EventGet: ToPbEventGet(event.(*mahjong.EventGet)),
			}}
	case mahjong.EventTypeDiscard:
		return &pb.Event{
			Event: &pb.Event_EventDiscard{
				EventDiscard: ToPbEventDiscard(event.(*mahjong.EventDiscard)),
			}}
	case mahjong.EventTypeTsumoGiri:
		return &pb.Event{
			Event: &pb.Event_EventDiscard{
				EventDiscard: ToPbEventTsumoGiri(event.(*mahjong.EventTsumoGiri)),
			}}
	case mahjong.EventTypeChi:
		return &pb.Event{
			Event: &pb.Event_EventCall{
				EventCall: ToPbEventChi(event.(*mahjong.EventChi)),
			}}
	case mahjong.EventTypePon:
		return &pb.Event{
			Event: &pb.Event_EventCall{
				EventCall: ToPbEventPon(event.(*mahjong.EventPon)),
			}}
	case mahjong.EventTypeDaiMinKan:
		return &pb.Event{
			Event: &pb.Event_EventCall{
				EventCall: ToPbEventDaiMinKan(event.(*mahjong.EventDaiMinKan)),
			}}
	case mahjong.EventTypeShouMinKan:
		return &pb.Event{
			Event: &pb.Event_EventCall{
				EventCall: ToPbEventShouMinKan(event.(*mahjong.EventShouMinKan)),
			}}
	case mahjong.EventTypeAnKan:
		return &pb.Event{
			Event: &pb.Event_EventCall{
				EventCall: ToPbEventAnKan(event.(*mahjong.EventAnKan)),
			}}
	case mahjong.EventTypeRiichi:
		return &pb.Event{
			Event: &pb.Event_EventRiichi{
				EventRiichi: ToPbEventRiichi(event.(*mahjong.EventRiichi)),
			}}
	case mahjong.EventTypeRon:
		return &pb.Event{
			Event: &pb.Event_EventRon{
				EventRon: ToPbEventRon(event.(*mahjong.EventRon)),
			}}
	case mahjong.EventTypeTsumo:
		return &pb.Event{
			Event: &pb.Event_EventTsumo{
				EventTsumo: ToPbEventTsumo(event.(*mahjong.EventTsumo)),
			}}
	case mahjong.EventTypeNewIndicator:
		return &pb.Event{
			Event: &pb.Event_EventNewIndicator{
				EventNewIndicator: ToPbEventNewIndicator(event.(*mahjong.EventNewIndicator)),
			}}
	case mahjong.EventTypeChanKan:
		return &pb.Event{
			Event: &pb.Event_EventRon{
				EventRon: ToPbEventChanKan(event.(*mahjong.EventChanKan)),
			}}
	case mahjong.EventTypeRyuuKyoku:
		return &pb.Event{
			Event: &pb.Event_EventRyuuKyoKu{
				EventRyuuKyoKu: ToPbEventRyuuKyoku(event.(*mahjong.EventRyuuKyoku)),
			}}
	case mahjong.EventTypeEnd:
		return &pb.Event{
			Event: &pb.Event_EventEnd{
				EventEnd: ToPbEventEnd(event.(*mahjong.EventEnd)),
			}}
	case mahjong.EventTypeFuriten:
		return &pb.Event{
			Event: &pb.Event_EventFuriten{
				EventFuriten: ToPbEventFuriten(event.(*mahjong.EventFuriten)),
			}}
	case mahjong.EventTypeNagashiMangan:
		return &pb.Event{
			Event: &pb.Event_EventNagashiMangan{
				EventNagashiMangan: ToPbNagashiMangan(event.(*mahjong.EventNagashiMangan)),
			}}
	case mahjong.EventTypeTenhaiEnd:
		return &pb.Event{
			Event: &pb.Event_EventTenhaiEnd{
				EventTenhaiEnd: ToPbTenhaiEnd(event.(*mahjong.EventTenhaiEnd)),
			}}
	default:
		panic("unknown event type")
	}
	return nil
}
