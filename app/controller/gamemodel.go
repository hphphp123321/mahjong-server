package controller

import (
	"github.com/hphphp123321/mahjong-go/mahjong"
	pb "github.com/hphphp123321/mahjong-server/app/api/v1"
)

func ToGameLimit(limit pb.Limit) mahjong.Limit {
	return mahjong.Limit(limit)
}

func ToGameRule(rule *pb.Rule) *mahjong.Rule {
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
