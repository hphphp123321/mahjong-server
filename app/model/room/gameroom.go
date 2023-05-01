package room

import (
	"github.com/hphphp123321/mahjong-go/mahjong"
	"github.com/hphphp123321/mahjong-server/app/errs"
	"github.com/hphphp123321/mahjong-server/app/global"
	"time"
)

type GameRoom struct {
	game       *mahjong.Game
	players    map[int]*mahjong.Player // int: Seat
	seat2Order map[int]int             // seat -> order

	PlayersErrChan      map[int]chan error
	PlayersAction       map[int]chan *mahjong.Call
	PlayersEvents       map[int]chan mahjong.Events
	PlayersValidActions map[int]chan mahjong.Calls
}

func NewGameRoom(game *mahjong.Game, seatOrder []int) *GameRoom {
	var seat2Order = make(map[int]int)
	var players = make(map[int]*mahjong.Player)
	for i, seat := range seatOrder {
		seat2Order[seat] = i
		players[seat] = mahjong.NewMahjongPlayer()
	}
	return &GameRoom{
		game:                game,
		players:             players,
		seat2Order:          seat2Order,
		PlayersErrChan:      make(map[int]chan error, 4),
		PlayersAction:       make(map[int]chan *mahjong.Call, 4),
		PlayersEvents:       make(map[int]chan mahjong.Events, 4),
		PlayersValidActions: make(map[int]chan mahjong.Calls, 4),
	}
}

func (r *GameRoom) GetPlayer(seat int) *mahjong.Player {
	return r.players[seat]
}

func (r *GameRoom) RegisterSeat(seat int, action chan *mahjong.Call) (ech chan mahjong.Events, vch chan mahjong.Calls, errCh chan error) {
	r.PlayersAction[seat] = action

	ech = make(chan mahjong.Events, 20)
	r.PlayersEvents[seat] = ech

	vch = make(chan mahjong.Calls, 1)
	r.PlayersValidActions[seat] = vch

	errCh = make(chan error, 1)
	r.PlayersErrChan[seat] = errCh
	return ech, vch, errCh
}

func (r *GameRoom) CheckRegister() bool {
	if len(r.PlayersAction) != 4 || len(r.PlayersEvents) != 4 || len(r.PlayersValidActions) != 4 || len(r.PlayersErrChan) != 4 {
		return false
	}
	for _, v := range r.PlayersAction {
		if v == nil {
			return false
		}
	}
	return true
}

func (r *GameRoom) sendEvent(seat int, events mahjong.Events) {
	r.PlayersEvents[seat] <- events
}

func (r *GameRoom) sendValidActions(seat int, actions mahjong.Calls) {
	r.PlayersValidActions[seat] <- actions
}

func (r *GameRoom) getSeatByWind(wind mahjong.Wind) int {
	for seat, player := range r.players {
		if player.Wind == wind {
			return seat
		}
	}
	return -1
}

func (r *GameRoom) getWindBySeat(seat int) mahjong.Wind {
	return r.players[seat].Wind
}

func (r *GameRoom) getBoardStateBySeat(seat int) *mahjong.BoardState {
	return r.game.GetPosBoardState(r.getWindBySeat(seat), nil)
}

func (r *GameRoom) StartGame() {
	go func() {
		for !r.CheckRegister() {
			time.Sleep(time.Millisecond * 100)
		}
		global.Log.Debugln("start game")
		// start game
		var flag = mahjong.EndTypeNone
		var players []*mahjong.Player
		for seat, _ := range r.seat2Order {
			players = append(players, r.players[seat])
		}
		playersEventIdx := make(map[mahjong.Wind]int, 4)

		posCalls := r.game.Reset(players, nil)
		posCall := make(map[mahjong.Wind]*mahjong.Call, 4)

		for flag != mahjong.EndTypeGame {

			if flag == mahjong.EndTypeRound {
				// round end, clear EventIdx
				for wind := range r.game.PosPlayer {
					playersEventIdx[wind] = 0
				}
			} else {
				// send events
				for wind := range r.game.PosPlayer {
					events := r.game.GetPosEvents(wind, playersEventIdx[wind])
					r.sendEvent(r.getSeatByWind(wind), events)
					playersEventIdx[wind] += len(events)
				}
			}

			// send valid actions
			for wind, calls := range posCalls {
				r.sendValidActions(r.getSeatByWind(wind), calls)
			}
			// recv action
			for wind := range posCalls {
				posCall[wind] = <-r.PlayersAction[r.getSeatByWind(wind)]
			}

			// step
			posCalls, flag = r.game.Step(posCall)
			posCall = make(map[mahjong.Wind]*mahjong.Call, 4)

		}

		// send end error
		for seat := range r.seat2Order {
			r.PlayersErrChan[seat] <- errs.ErrGameEnd
		}
	}()
}
