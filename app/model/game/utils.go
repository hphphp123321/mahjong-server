package game

import (
	"github.com/hphphp123321/go-common"
	"github.com/hphphp123321/mahjong-go/mahjong"
)

func GetTestTiles(mode int) mahjong.Tiles {
	if mode >= 0 {
		return nil
	}

	switch mode {
	case -1:
		return GetAllTiles(KuokushiTiles)
	case -2:
		return GetAllTiles(TsuyisouTiles)
	}

	return nil
}

var KuokushiTiles = mahjong.Tiles{
	mahjong.Man1T1, mahjong.Man9T1, mahjong.Pin1T1, mahjong.Pin9T1, mahjong.Sou1T1, mahjong.Sou9T1,
	mahjong.Ton1, mahjong.Nan1, mahjong.Shaa1, mahjong.Pei1, mahjong.Haku1, mahjong.Hatsu1, mahjong.Chun1, mahjong.Man1T2, // 这里多给了一张Man1T2
}

var TsuyisouTiles = mahjong.Tiles{
	108, 109, 110, 112, 113, 114, 116, 117, 118, 120, 121, 122, 124, 125, // 这里多给了一张Haku1 2作为雀头
}

func GetAllTiles(tonTiles mahjong.Tiles) mahjong.Tiles {
	if len(tonTiles) != 14 {
		return nil
	}
	var tiles = mahjong.Tiles{}
	for i := 0; i < 136; i++ {
		tiles = append(tiles, mahjong.Tile(i))
	}
	for _, t := range tonTiles {
		tiles, _ = common.Remove(tiles, t)
	}

	// 将tonTiles的前13张添加到tiles的前13张位置
	tiles = append(tonTiles[:13], tiles...)

	// 将tonTiles的最后一张插入到tiles[52]的位置
	tiles = append(tiles[:52], append([]mahjong.Tile{tonTiles[13]}, tiles[52:]...)...)

	return tiles
}
