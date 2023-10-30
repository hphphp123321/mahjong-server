package defaultloader

import (
	"github.com/hphphp123321/mahjong-server/app/bootloader"
	"github.com/hphphp123321/mahjong-server/app/bootloader/loader/baseloader"
	"github.com/hphphp123321/mahjong-server/app/bootloader/loader/configloader"
	"github.com/hphphp123321/mahjong-server/app/bootloader/loader/grpcloader"
	"github.com/hphphp123321/mahjong-server/app/bootloader/loader/robotloader"
	"github.com/hphphp123321/mahjong-server/app/bootloader/loader/zaploggerloader"
)

var (
	List = bootloader.LoaderList{
		&baseloader.BaseLoader{},
		&configloader.ConfigLoader{},
		&zaploggerloader.ZapLoggerLoader{},
		&robotloader.RobotLoader{},
		&grpcloader.GrpcLoader{},
	}
)
