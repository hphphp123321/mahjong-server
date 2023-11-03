package robot

import (
	"github.com/hphphp123321/go-common"
	"github.com/hphphp123321/mahjong-go/mahjong"
	pb "github.com/hphphp123321/mahjong-server/app/api/v1/robot"
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
		Who:         ToPbWind(event.Who),
		Tile:        ToPbTile(event.Tile),
		TenpaiInfos: ToPbTenpaiInfos(event.TenpaiInfos),
	}
}

func ToPbEventDiscard(event *mahjong.EventDiscard) *pb.EventDiscard {
	return &pb.EventDiscard{
		Who:        ToPbWind(event.Who),
		Tile:       ToPbTile(event.Tile),
		TsumoGiri:  false,
		TenpaiInfo: ToPbTenpaiInfo(event.TenpaiInfo),
	}
}

func ToPbEventTsumoGiri(event *mahjong.EventTsumoGiri) *pb.EventDiscard {
	return &pb.EventDiscard{
		Who:        ToPbWind(event.Who),
		Tile:       ToPbTile(event.Tile),
		TsumoGiri:  true,
		TenpaiInfo: ToPbTenpaiInfo(event.TenpaiInfo),
	}
}

func ToPbEventChi(event *mahjong.EventChi) *pb.EventCall {
	return &pb.EventCall{
		Who:         ToPbWind(event.Who),
		Call:        ToPbCall(event.Call),
		TenpaiInfos: ToPbTenpaiInfos(event.TenpaiInfos),
	}
}

func ToPbEventPon(event *mahjong.EventPon) *pb.EventCall {
	return &pb.EventCall{
		Who:         ToPbWind(event.Who),
		Call:        ToPbCall(event.Call),
		TenpaiInfos: ToPbTenpaiInfos(event.TenpaiInfos),
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
		Reason:    ToPbRyuuKyoKuReason(event.Reason),
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

func ToPbTenpaiEnd(event *mahjong.EventTenpaiEnd) *pb.EventTenpaiEnd {
	return &pb.EventTenpaiEnd{
		Who:         ToPbWind(event.Who),
		HandTiles:   ToPbTiles(event.HandTiles),
		TenpaiSlice: ToPbTileClasses(event.TenpaiSlice),
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
	case mahjong.EventTypeTenpaiEnd:
		return &pb.Event{
			Event: &pb.Event_EventTenpaiEnd{
				EventTenpaiEnd: ToPbTenpaiEnd(event.(*mahjong.EventTenpaiEnd)),
			}}
	default:
		panic("unknown event type")
	}
	return nil
}

// -------------------------------------------------------------------------------
// To Mahjong

func ToMahjongEventStart(event *pb.EventStart) *mahjong.EventStart {
	return &mahjong.EventStart{
		WindRound:         ToMahjongWindRound(event.WindRound),
		InitWind:          ToMahjongWind(event.InitWind),
		Seed:              event.Seed,
		NumGame:           int(event.NumGame),
		NumHonba:          int(event.NumHonba),
		NumRiichi:         int(event.NumRiichi),
		InitDoraIndicator: ToMahjongTile(event.InitDoraIndicator),
		InitTiles:         ToMahjongTiles(event.InitTiles),
		Rule:              ToMahjongGameRule(event.GameRule),
		PlayersPoints: map[mahjong.Wind]int{
			mahjong.East:  int(event.PlayersPoints[0]),
			mahjong.South: int(event.PlayersPoints[1]),
			mahjong.West:  int(event.PlayersPoints[2]),
			mahjong.North: int(event.PlayersPoints[3]),
		},
	}
}

func ToMahjongEventGet(event *pb.EventGet) *mahjong.EventGet {
	return &mahjong.EventGet{
		Who:         ToMahjongWind(event.Who),
		Tile:        ToMahjongTile(event.Tile),
		TenpaiInfos: ToMahjongTenpaiInfos(event.TenpaiInfos),
	}
}

func ToMahjongEventDiscard(event *pb.EventDiscard) mahjong.Event {
	if event.TsumoGiri {
		return &mahjong.EventTsumoGiri{
			Who:        ToMahjongWind(event.Who),
			Tile:       ToMahjongTile(event.Tile),
			TenpaiInfo: ToMahjongTenpaiInfo(event.TenpaiInfo),
		}
	}
	return &mahjong.EventDiscard{
		Who:        ToMahjongWind(event.Who),
		Tile:       ToMahjongTile(event.Tile),
		TenpaiInfo: ToMahjongTenpaiInfo(event.TenpaiInfo),
	}
}

func ToMahjongEventCall(event *pb.EventCall) mahjong.Event {
	who := ToMahjongWind(event.Who)
	call := ToMahjongCall(event.Call)
	switch call.CallType {
	case mahjong.Chi:
		return &mahjong.EventChi{Who: who, Call: call}
	case mahjong.Pon:
		return &mahjong.EventPon{Who: who, Call: call}
	case mahjong.DaiMinKan:
		return &mahjong.EventDaiMinKan{Who: who, Call: call}
	case mahjong.ShouMinKan:
		return &mahjong.EventShouMinKan{Who: who, Call: call}
	case mahjong.AnKan:
		return &mahjong.EventAnKan{Who: who, Call: call}
	default:
		panic("unknown call type")
	}
}

func ToMahjongEventRiichi(event *pb.EventRiichi) *mahjong.EventRiichi {
	return &mahjong.EventRiichi{
		Who:  ToMahjongWind(event.Who),
		Step: int(event.Step),
	}
}

func ToMahjongEventRon(event *pb.EventRon) *mahjong.EventRon {
	return &mahjong.EventRon{
		Who:       ToMahjongWind(event.Who),
		FromWho:   ToMahjongWind(event.FromWho),
		HandTiles: ToMahjongTiles(event.HandTiles),
		WinTile:   ToMahjongTile(event.WinTile),
		Result:    ToMahjongResult(event.Result),
	}
}

func ToMahjongEventTsumo(event *pb.EventTsumo) *mahjong.EventTsumo {
	return &mahjong.EventTsumo{
		Who:       ToMahjongWind(event.Who),
		HandTiles: ToMahjongTiles(event.HandTiles),
		WinTile:   ToMahjongTile(event.WinTile),
		Result:    ToMahjongResult(event.Result),
	}
}

func ToMahjongEventNewIndicator(event *pb.EventNewIndicator) *mahjong.EventNewIndicator {
	return &mahjong.EventNewIndicator{
		Tile: ToMahjongTile(event.Tile),
	}
}

func ToMahjongEventRyuuKyoku(event *pb.EventRyuuKyoKu) *mahjong.EventRyuuKyoku {
	return &mahjong.EventRyuuKyoku{
		Who:       ToMahjongWind(event.Who),
		HandTiles: ToMahjongTiles(event.HandTiles),
		Reason:    ToMahjongRyuuKyoKuReason(event.Reason),
	}
}

func ToMahjongEventEnd(event *pb.EventEnd) *mahjong.EventEnd {
	return &mahjong.EventEnd{
		PointsChange: map[mahjong.Wind]int{
			mahjong.East:  int(event.EastPointsChange),
			mahjong.South: int(event.SouthPointsChange),
			mahjong.West:  int(event.WestPointsChange),
			mahjong.North: int(event.NorthPointsChange),
		},
	}
}

func ToMahjongEventFuriten(event *pb.EventFuriten) *mahjong.EventFuriten {
	return &mahjong.EventFuriten{
		Who:           ToMahjongWind(event.Who),
		FuritenReason: mahjong.FuritenReason(event.Reason),
	}
}

func ToMahjongEventNagashiMangan(event *pb.EventNagashiMangan) *mahjong.EventNagashiMangan {
	return &mahjong.EventNagashiMangan{
		Who: ToMahjongWind(event.Who),
	}
}

func ToMahjongEventTenhaiEnd(event *pb.EventTenpaiEnd) *mahjong.EventTenpaiEnd {
	return &mahjong.EventTenpaiEnd{
		Who:         ToMahjongWind(event.Who),
		HandTiles:   ToMahjongTiles(event.HandTiles),
		TenpaiSlice: ToMahjongTileClasses(event.TenpaiSlice),
	}
}

func ToMahjongEvents(events []*pb.Event) []mahjong.Event {
	return common.MapSlice(events, ToMahjongEvent)
}

func ToMahjongEvent(event *pb.Event) mahjong.Event {
	switch e := event.Event.(type) {
	case *pb.Event_EventStart:
		return ToMahjongEventStart(e.EventStart)
	case *pb.Event_EventGet:
		return ToMahjongEventGet(e.EventGet)
	case *pb.Event_EventDiscard:
		return ToMahjongEventDiscard(e.EventDiscard)
	case *pb.Event_EventCall:
		return ToMahjongEventCall(e.EventCall)
	case *pb.Event_EventRiichi:
		return ToMahjongEventRiichi(e.EventRiichi)
	case *pb.Event_EventRon:
		return ToMahjongEventRon(e.EventRon)
	case *pb.Event_EventTsumo:
		return ToMahjongEventTsumo(e.EventTsumo)
	case *pb.Event_EventNewIndicator:
		return ToMahjongEventNewIndicator(e.EventNewIndicator)
	case *pb.Event_EventRyuuKyoKu:
		return ToMahjongEventRyuuKyoku(e.EventRyuuKyoKu)
	case *pb.Event_EventEnd:
		return ToMahjongEventEnd(e.EventEnd)
	case *pb.Event_EventFuriten:
		return ToMahjongEventFuriten(e.EventFuriten)
	case *pb.Event_EventNagashiMangan:
		return ToMahjongEventNagashiMangan(e.EventNagashiMangan)
	case *pb.Event_EventTenpaiEnd:
		return ToMahjongEventTenhaiEnd(e.EventTenpaiEnd)
	default:
		panic("unknown event type")
	}
	return nil
}
