package room

import (
	"github.com/hphphp123321/mahjong-go/mahjong"
	"github.com/hphphp123321/mahjong-server/app/errs"
	"github.com/hphphp123321/mahjong-server/app/global"
	"github.com/hphphp123321/mahjong-server/app/model/player"
	"time"
)

type GameRoom struct {
	game       *mahjong.Game
	players    map[int]*mahjong.Player // int: Seat
	seat2Order map[int]int             // seat -> order

	//PlayersErrChan      map[int]chan error
	PlayersAction map[int]chan *mahjong.Call
	//PlayersEvents       map[int]chan mahjong.Events
	//PlayersValidActions map[int]chan mahjong.Calls
	PlayersEventsChan map[int]chan *player.GameEventChannel
}

func NewGameRoom(game *mahjong.Game, seatOrder []int) *GameRoom {
	var seat2Order = make(map[int]int)
	var players = make(map[int]*mahjong.Player)
	for i, seat := range seatOrder {
		seat2Order[seat] = i
		players[seat] = mahjong.NewMahjongPlayer()
	}
	return &GameRoom{
		game:              game,
		players:           players,
		seat2Order:        seat2Order,
		PlayersAction:     make(map[int]chan *mahjong.Call, 4),
		PlayersEventsChan: make(map[int]chan *player.GameEventChannel, 4),
	}
}

func (r *GameRoom) GetPlayer(seat int) *mahjong.Player {
	return r.players[seat]
}

func (r *GameRoom) RegisterSeat(seat int, action chan *mahjong.Call) chan *player.GameEventChannel {
	r.PlayersAction[seat] = action

	eventChan := make(chan *player.GameEventChannel)
	r.PlayersEventsChan[seat] = eventChan
	return eventChan
}

func (r *GameRoom) CheckRegister() bool {
	if len(r.PlayersAction) != 4 || len(r.PlayersEventsChan) != 4 {
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
	r.PlayersEventsChan[seat] <- &player.GameEventChannel{
		Events: events,
	}
}

func (r *GameRoom) sendValidActions(seat int, actions mahjong.Calls) {
	r.PlayersEventsChan[seat] <- &player.GameEventChannel{
		ValidActions: actions,
	}
}

func (r *GameRoom) getSeatByWind(wind mahjong.Wind) int {
	for seat, p := range r.players {
		if p.Wind == wind {
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

		defer func() {
			if err := recover(); err != nil {
				// send end error
				for seat := range r.seat2Order {
					select {
					case r.PlayersEventsChan[seat] <- &player.GameEventChannel{Err: errs.ErrGameEndUnexpect}:
					default:
						global.Log.Warnln("send end error failed")
					}
				}
			}
		}()

		// start game
		var flag = mahjong.EndTypeNone
		var players []*mahjong.Player
		for seat, _ := range r.seat2Order {
			players = append(players, r.players[seat])
		}
		playersEventIdx := make(map[mahjong.Wind]int, 4)

		posCalls := r.game.Reset(players, nil)
		posCall := make(map[mahjong.Wind]*mahjong.Call, 4)

		for {

			// send events
			for wind := range r.game.PosPlayer {
				events := r.game.GetPosEvents(wind, playersEventIdx[wind])
				r.sendEvent(r.getSeatByWind(wind), events)
				playersEventIdx[wind] += len(events)
			}

			if flag == mahjong.EndTypeRound {
				// round end, clear EventIdx
				for wind := range r.game.PosPlayer {
					playersEventIdx[wind] = 0
				}
			}

			if flag == mahjong.EndTypeGame {
				// game end
				break
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
			r.PlayersEventsChan[seat] <- &player.GameEventChannel{Err: errs.ErrGameEnd}
		}
	}()
}
