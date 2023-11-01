package controller

import (
	"github.com/hphphp123321/go-common"
	"github.com/hphphp123321/mahjong-go/mahjong"
	pb "github.com/hphphp123321/mahjong-server/app/api/v1/mahjong"
)

func ToPbGameChatReply(in *pb.GameRequest, playerName string, seat int) *pb.GameReply {
	return &pb.GameReply{
		Message: ChatMsg(playerName, in.GetChat().Message),
		Reply: &pb.GameReply_Chat{Chat: &pb.ChatReply{
			Message: in.GetChat().Message,
			Seat:    int32(seat),
		}},
	}
}

func ToGameLimit(limit pb.Limit) mahjong.Limit {
	return mahjong.Limit(limit)
}

func ToMahjongGameRule(rule *pb.Rule) *mahjong.Rule {
	return &mahjong.Rule{
		GameLength:           int(rule.GameLength),
		IsOpenTanyao:         rule.IsOpenTanyao,
		HasAkaDora:           rule.HasAkaDora,
		RenhouLimit:          ToGameLimit(rule.RenhouLimit),
		IsHaiteiFromLiveOnly: rule.IsHaiteiFromLiveOnly,
		IsUra:                rule.IsUra,
		IsIpatsu:             rule.IsIpatsu,
		IsGreenRequired:      rule.IsGreenRequired,
		IsRinshanFu:          rule.IsRinshanFu,
		IsManganRound:        rule.IsManganRound,
		IsKazoeYakuman:       rule.IsKazoeYakuman,
		IsDoubleYakumans:     rule.IsDoubleYakumans,
		IsYakumanSum:         rule.IsYakumanSum,
		HonbaValue:           int(rule.HonbaValue),
		IsSanChaHou:          rule.IsSanChaHou,
		IsNagashiMangan:      rule.IsNagashiMangan,
	}
}

func ToPbGameRule(rule *mahjong.Rule) *pb.Rule {
	return &pb.Rule{
		GameLength:           int32(rule.GameLength),
		IsOpenTanyao:         rule.IsOpenTanyao,
		HasAkaDora:           rule.HasAkaDora,
		RenhouLimit:          pb.Limit(rule.RenhouLimit),
		IsHaiteiFromLiveOnly: rule.IsHaiteiFromLiveOnly,
		IsUra:                rule.IsUra,
		IsIpatsu:             rule.IsIpatsu,
		IsGreenRequired:      rule.IsGreenRequired,
		IsRinshanFu:          rule.IsRinshanFu,
		IsManganRound:        rule.IsManganRound,
		IsKazoeYakuman:       rule.IsKazoeYakuman,
		IsDoubleYakumans:     rule.IsDoubleYakumans,
		IsYakumanSum:         rule.IsYakumanSum,
		HonbaValue:           int32(rule.HonbaValue),
		IsSanChaHou:          rule.IsSanChaHou,
		IsNagashiMangan:      rule.IsNagashiMangan,
	}
}

func ToPbTile(t mahjong.Tile) pb.Tile {
	return pb.Tile(t + 1)
}

func ToMahjongTile(t pb.Tile) mahjong.Tile {
	return mahjong.Tile(t - 1)
}

func ToPbTiles(tiles mahjong.Tiles) []pb.Tile {
	return common.MapSlice(tiles, ToPbTile)
}

func ToMahjongTiles(tiles []pb.Tile) mahjong.Tiles {
	return common.MapSlice(tiles, ToMahjongTile)
}

func ToPbTileClass(tc mahjong.TileClass) pb.TileClass {
	return pb.TileClass(tc + 1)
}

func ToMahjongTileClass(tc pb.TileClass) mahjong.TileClass {
	return mahjong.TileClass(tc - 1)
}

func ToPbTileClasses(tcs mahjong.TileClasses) []pb.TileClass {
	return common.MapSlice(tcs, ToPbTileClass)
}

func ToMahjongTileClasses(tcs []pb.TileClass) mahjong.TileClasses {
	return common.MapSlice(tcs, ToMahjongTileClass)
}

func ToPbWindRound(w mahjong.WindRound) pb.WindRound {
	return pb.WindRound(w)
}

func ToMahjongWindRound(w pb.WindRound) mahjong.WindRound {
	return mahjong.WindRound(w)
}

func ToPbWind(w mahjong.Wind) pb.Wind {
	return pb.Wind(w + 1)
}

func ToMahjongWind(w pb.Wind) mahjong.Wind {
	return mahjong.Wind(w - 1)
}

func ToPbWinds(ws []mahjong.Wind) []pb.Wind {
	return common.MapSlice(ws, ToPbWind)
}

func ToMahjongWinds(ws []pb.Wind) []mahjong.Wind {
	return common.MapSlice(ws, ToMahjongWind)
}

func ToPbCallType(ct mahjong.CallType) pb.CallType {
	return pb.CallType(ct)
}

func ToMahjongCallType(ct pb.CallType) mahjong.CallType {
	return mahjong.CallType(ct)
}

func ToPbCall(c *mahjong.Call) *pb.Call {
	return &pb.Call{
		Type:    ToPbCallType(c.CallType),
		Tiles:   ToPbTiles(c.CallTiles),
		FromWho: ToPbWinds(c.CallTilesFromWho),
	}
}

func ToMahjongCall(c *pb.Call) *mahjong.Call {
	return &mahjong.Call{
		CallType:         ToMahjongCallType(c.Type),
		CallTiles:        ToMahjongTiles(c.Tiles),
		CallTilesFromWho: ToMahjongWinds(c.FromWho),
	}
}

func ToPbCalls(cs []*mahjong.Call) []*pb.Call {
	return common.MapSlice(cs, ToPbCall)
}

func ToMahjongCalls(cs []*pb.Call) []*mahjong.Call {
	return common.MapSlice(cs, ToMahjongCall)
}

func ToPbPlayerState(ps *mahjong.PlayerState) *pb.PlayerState {
	return &pb.PlayerState{
		Points:         int32(ps.Points),
		Melds:          ToPbCalls(ps.Melds),
		DiscardTiles:   ToPbTiles(ps.DiscardTiles),
		TilesTsumoGiri: ps.TilesTsumoGiri,
		IsRiichi:       ps.IsRiichi,
	}
}

func ToMahjongPlayerState(ps *pb.PlayerState) *mahjong.PlayerState {
	return &mahjong.PlayerState{
		Points:         int(ps.Points),
		Melds:          ToMahjongCalls(ps.Melds),
		DiscardTiles:   ToMahjongTiles(ps.DiscardTiles),
		TilesTsumoGiri: ps.TilesTsumoGiri,
		IsRiichi:       ps.IsRiichi,
	}
}

func ToPbBoardState(b *mahjong.BoardState) *pb.BoardState {
	return &pb.BoardState{
		WindRound:      ToPbWindRound(b.WindRound),
		NumHonba:       int32(b.NumHonba),
		NumRiichi:      int32(b.NumRiichi),
		DoraIndicators: ToPbTiles(b.DoraIndicators),
		PlayerWind:     ToPbWind(b.PlayerWind),
		Position:       ToPbWind(b.Position),
		HandTiles:      ToPbTiles(b.HandTiles),
		ValidActions:   ToPbCalls(b.ValidActions),
		NumRemainTiles: int32(b.NumRemainTiles),
		PlayerStates: func() map[int32]*pb.PlayerState {
			var m = make(map[int32]*pb.PlayerState)
			for k, v := range b.PlayerStates {
				m[int32(ToPbWind(k))] = ToPbPlayerState(v)
			}
			return m
		}(),
	}
}

func ToMahjongBoardState(b *pb.BoardState) *mahjong.BoardState {
	return &mahjong.BoardState{
		WindRound:      ToMahjongWindRound(b.WindRound),
		NumHonba:       int(b.NumHonba),
		NumRiichi:      int(b.NumRiichi),
		DoraIndicators: ToMahjongTiles(b.DoraIndicators),
		PlayerWind:     ToMahjongWind(b.PlayerWind),
		Position:       ToMahjongWind(b.Position),
		HandTiles:      ToMahjongTiles(b.HandTiles),
		ValidActions:   ToMahjongCalls(b.ValidActions),
		NumRemainTiles: int(b.NumRemainTiles),
		PlayerStates: func() map[mahjong.Wind]*mahjong.PlayerState {
			var m = make(map[mahjong.Wind]*mahjong.PlayerState)
			for k, v := range b.PlayerStates {
				m[ToMahjongWind(pb.Wind(k))] = ToMahjongPlayerState(v)
			}
			return m
		}(),
	}
}

func ToPbYakuSets(yakuSet mahjong.YakuSet) []*pb.YakuSet {
	yakuSets := make([]*pb.YakuSet, 0, len(yakuSet))
	for yaku, han := range yakuSet {
		yakuSets = append(yakuSets, &pb.YakuSet{
			Yaku: pb.Yaku(yaku),
			Han:  int32(han),
		})
	}
	return yakuSets
}

func ToMahjongYakuSets(yakuSets []*pb.YakuSet) mahjong.YakuSet {
	yakuSet := make(mahjong.YakuSet, len(yakuSets))
	for _, y := range yakuSets {
		yakuSet[mahjong.Yaku(y.Yaku)] = int(y.Han)
	}
	return yakuSet
}

func ToPbYakuman(yakuman mahjong.Yakuman) pb.Yakuman {
	return pb.Yakuman(yakuman)
}

func ToMahjongYakuman(yakuman pb.Yakuman) mahjong.Yakuman {
	return mahjong.Yakuman(yakuman)
}

func ToPbYakuMans(yakuMans mahjong.Yakumans) []pb.Yakuman {
	return common.MapSlice(yakuMans, ToPbYakuman)
}

func ToMahjongYakuMans(yakuMans []pb.Yakuman) mahjong.Yakumans {
	return common.MapSlice(yakuMans, ToMahjongYakuman)
}

func ToPbFuInfo(fuInfo *mahjong.FuInfo) *pb.FuInfo {
	return &pb.FuInfo{
		Fu:     pb.Fu(fuInfo.Fu),
		Points: int32(fuInfo.Points),
	}
}

func ToMahjongFuInfo(fuInfo *pb.FuInfo) *mahjong.FuInfo {
	return &mahjong.FuInfo{
		Fu:     mahjong.Fu(fuInfo.Fu),
		Points: int(fuInfo.Points),
	}
}

func ToPbFuInfos(fuInfos []*mahjong.FuInfo) []*pb.FuInfo {
	return common.MapSlice(fuInfos, ToPbFuInfo)
}

func ToMahjongFuInfos(fuInfos []*pb.FuInfo) []*mahjong.FuInfo {
	return common.MapSlice(fuInfos, ToMahjongFuInfo)
}

func ToPbYakuResult(yaku *mahjong.YakuResult) *pb.YakuResult {
	return &pb.YakuResult{
		YakuSets: ToPbYakuSets(yaku.Yaku),
		Yakumans: ToPbYakuMans(yaku.Yakumans),
		Bonuses:  ToPbYakuSets(yaku.Bonuses),
		Fus:      ToPbFuInfos(yaku.Fus),
		IsClosed: yaku.IsClosed,
	}
}

func ToMahjongYakuResult(yaku *pb.YakuResult) *mahjong.YakuResult {
	return &mahjong.YakuResult{
		Yaku:     ToMahjongYakuSets(yaku.YakuSets),
		Yakumans: ToMahjongYakuMans(yaku.Yakumans),
		Bonuses:  ToMahjongYakuSets(yaku.Bonuses),
		Fus:      ToMahjongFuInfos(yaku.Fus),
		IsClosed: yaku.IsClosed,
	}
}

func ToPbScoreResult(score *mahjong.ScoreResult) *pb.ScoreResult {
	return &pb.ScoreResult{
		PayRon:         int32(score.PayRon),
		PayRonDealer:   int32(score.PayRonDealer),
		PayTsumo:       int32(score.PayTsumo),
		PayTsumoDealer: int32(score.PayTsumoDealer),
		Special:        pb.Limit(score.Special),
		Han:            int32(score.Han),
		Fu:             int32(score.Fu),
	}
}

func ToMahjongScoreResult(score *pb.ScoreResult) *mahjong.ScoreResult {
	return &mahjong.ScoreResult{
		PayRon:         int(score.PayRon),
		PayRonDealer:   int(score.PayRonDealer),
		PayTsumo:       int(score.PayTsumo),
		PayTsumoDealer: int(score.PayTsumoDealer),
		Special:        mahjong.Limit(score.Special),
		Han:            int(score.Han),
		Fu:             int(score.Fu),
	}
}

func ToPbResult(result *mahjong.Result) *pb.Result {
	if result == nil {
		return nil
	}
	return &pb.Result{
		YakuResult:  ToPbYakuResult(result.YakuResult),
		ScoreResult: ToPbScoreResult(result.ScoreResult),
	}
}

func ToMahjongResult(result *pb.Result) *mahjong.Result {
	if result == nil {
		return nil
	}
	return &mahjong.Result{
		YakuResult:  ToMahjongYakuResult(result.YakuResult),
		ScoreResult: ToMahjongScoreResult(result.ScoreResult),
	}
}

func ToPbTenpaiResult(result *mahjong.TenpaiResult) *pb.TenpaiResult {
	if result == nil {
		return nil
	}
	return &pb.TenpaiResult{
		RemainNum: int32(result.RemainNum),
		Result:    ToPbResult(result.Result),
	}
}

func ToMahjongTenpaiResult(result *pb.TenpaiResult) *mahjong.TenpaiResult {
	if result == nil {
		return nil
	}
	return &mahjong.TenpaiResult{
		RemainNum: int(result.RemainNum),
		Result:    ToMahjongResult(result.Result),
	}
}

func ToPbTenpaiInfo(info *mahjong.TenpaiInfo) *pb.TenpaiInfo {
	if info == nil {
		return nil
	}
	var tileClassesTenpaiResult = make(map[int32]*pb.TenpaiResult)
	for tileClass, result := range info.TileClassesTenpaiResult {
		tileClassesTenpaiResult[int32(ToPbTileClass(tileClass))] = ToPbTenpaiResult(result)
	}
	return &pb.TenpaiInfo{
		TileClassesTenpaiResult: tileClassesTenpaiResult,
		Furiten:                 info.Furiten,
	}
}

func ToMahjongTenpaiInfo(info *pb.TenpaiInfo) *mahjong.TenpaiInfo {
	if info == nil {
		return nil
	}
	var tileClassesTenpaiResult = make(map[mahjong.TileClass]*mahjong.TenpaiResult)
	for tileClass, result := range info.TileClassesTenpaiResult {
		tileClassesTenpaiResult[ToMahjongTileClass(pb.TileClass(tileClass))] = ToMahjongTenpaiResult(result)
	}
	return &mahjong.TenpaiInfo{
		TileClassesTenpaiResult: tileClassesTenpaiResult,
		Furiten:                 info.Furiten,
	}
}

func ToPbTenpaiInfos(infos mahjong.TenpaiInfos) map[int32]*pb.TenpaiInfo {
	if infos == nil {
		return nil
	}
	var m = make(map[int32]*pb.TenpaiInfo)
	for tile, info := range infos {
		m[int32(ToPbTile(tile))] = ToPbTenpaiInfo(info)
	}
	return m
}

func ToMahjongTenpaiInfos(infos map[int32]*pb.TenpaiInfo) mahjong.TenpaiInfos {
	if infos == nil {
		return nil
	}
	var m = make(mahjong.TenpaiInfos)
	for tile, info := range infos {
		m[ToMahjongTile(pb.Tile(tile))] = ToMahjongTenpaiInfo(info)
	}
	return m
}

func ToPbRyuuKyoKuReason(reason mahjong.RyuuKyokuReason) pb.RyuuKyoKuReason {
	return pb.RyuuKyoKuReason(reason - 1)
}

func ToMahjongRyuuKyoKuReason(reason pb.RyuuKyoKuReason) mahjong.RyuuKyokuReason {
	return mahjong.RyuuKyokuReason(reason + 1)
}
