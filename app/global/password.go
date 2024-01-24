package global

import "github.com/hphphp123321/mahjong-server/app/service/password"

var PasswordService password.Password = password.Bcrypt{}
