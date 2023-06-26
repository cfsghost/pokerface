package table

import (
	"errors"
	"time"

	"github.com/weedbox/syncsaga"
)

var (
	ErrPlayerNotInGame      = errors.New("table: player not in the game")
	ErrTimesUp              = errors.New("table: time's up")
	ErrGameConditionsNotMet = errors.New("table: game conditions not met")
	ErrMaxGamesExceeded     = errors.New("table: reach the maximum number of games")
)

type TableOpt func(*table)

type Table interface {
	Start() error
	Close() error
	Resume() error
	Pause() error
	GetState() *State
	GetGame() Game
	GetGameCount() int
	GetPlayerIdx(playerID string) int

	SetAnte(chips int64)
	SetBlinds(dealer int64, sb int64, bb int64)

	// Event
	OnStateUpdated(func(*State))

	// Actions
	Ready(playerID string) error
	Pass(playerID string) error
	Pay(playerID string, chips int64) error
	Fold(playerID string) error
	Check(playerID string) error
	Call(playerID string) error
	Allin(playerID string) error
	Bet(playerID string, chips int64) error
	Raise(playerID string, chipLevel int64) error
}

type table struct {
	g              Game
	b              Backend
	isPaused       bool
	inPosition     bool
	options        *Options
	gameCount      int
	ts             *State
	rg             *syncsaga.ReadyGroup
	sm             *SeatManager
	onStateUpdated func(*State)
}

func WithBackend(b Backend) TableOpt {
	return func(t *table) {
		t.b = b
	}
}

func NewTable(options *Options, opts ...TableOpt) *table {

	t := &table{
		options:        options,
		rg:             syncsaga.NewReadyGroup(),
		sm:             NewSeatManager(options.MaxSeats),
		ts:             NewState(),
		onStateUpdated: func(*State) {},
	}

	for _, opt := range opts {
		opt(t)
	}

	t.ts.Status = "idle"

	return t
}

func (t *table) OnStateUpdated(fn func(*State)) {
	t.onStateUpdated = fn
}

func (t *table) GetState() *State {
	return t.ts
}

func (t *table) GetGame() Game {
	return t.g
}

func (t *table) GetGameCount() int {
	return t.gameCount
}

func (t *table) SetAnte(chips int64) {
	t.options.Ante = chips
}

func (t *table) SetBlinds(dealer int64, sb int64, bb int64) {
	t.options.Blind.Dealer = dealer
	t.options.Blind.SB = sb
	t.options.Blind.BB = bb
}

func (t *table) Start() error {

	t.ts.StartTime = time.Now().Unix()
	t.ts.EndTime = t.ts.StartTime + int64(t.options.Duration)

	go t.nextGame(0)

	return nil
}

func (t *table) Close() error {

	t.ts.Status = "closed"

	return nil
}

func (t *table) Resume() error {

	if !t.isPaused {
		return nil
	}

	return t.run(0)
}

func (t *table) Pause() error {

	t.isPaused = true

	return nil
}

func (t *table) Join(seatID int, p *PlayerInfo) error {
	return t.sm.Join(seatID, p)
}

func (t *table) Leave(seatID int) error {
	return t.sm.Leave(seatID)
}

func (t *table) GetPlayerIdx(playerID string) int {

	for _, s := range t.sm.seats {
		if s.Player == nil {
			continue
		}

		if s.Player.ID == playerID {
			return s.Player.GameIdx
		}
	}

	return -1
}

// Actions
func (t *table) Ready(playerID string) error {

	idx := t.GetPlayerIdx(playerID)
	if idx == -1 {
		return ErrPlayerNotInGame
	}

	err := t.g.Ready(idx)
	if err != nil {
		return err
	}

	return nil
}

func (t *table) Pass(playerID string) error {

	idx := t.GetPlayerIdx(playerID)
	if idx == -1 {
		return ErrPlayerNotInGame
	}

	err := t.g.Pass(idx)
	if err != nil {
		return err
	}

	return nil
}

func (t *table) Pay(playerID string, chips int64) error {

	idx := t.GetPlayerIdx(playerID)
	if idx == -1 {
		return ErrPlayerNotInGame
	}

	err := t.g.Pay(idx, chips)
	if err != nil {
		return err
	}

	return nil
}

func (t *table) Fold(playerID string) error {

	idx := t.GetPlayerIdx(playerID)
	if idx == -1 {
		return ErrPlayerNotInGame
	}

	err := t.g.Fold(idx)
	if err != nil {
		return err
	}

	return nil
}

func (t *table) Check(playerID string) error {

	idx := t.GetPlayerIdx(playerID)
	if idx == -1 {
		return ErrPlayerNotInGame
	}

	err := t.g.Check(idx)
	if err != nil {
		return err
	}

	return nil
}

func (t *table) Call(playerID string) error {

	idx := t.GetPlayerIdx(playerID)
	if idx == -1 {
		return ErrPlayerNotInGame
	}

	err := t.g.Call(idx)
	if err != nil {
		return err
	}

	return nil
}

func (t *table) Allin(playerID string) error {

	idx := t.GetPlayerIdx(playerID)
	if idx == -1 {
		return ErrPlayerNotInGame
	}

	err := t.g.Allin(idx)
	if err != nil {
		return err
	}

	return nil
}

func (t *table) Bet(playerID string, chips int64) error {

	idx := t.GetPlayerIdx(playerID)
	if idx == -1 {
		return ErrPlayerNotInGame
	}

	err := t.g.Bet(idx, chips)
	if err != nil {
		return err
	}

	return nil
}

func (t *table) Raise(playerID string, chipLevel int64) error {

	idx := t.GetPlayerIdx(playerID)
	if idx == -1 {
		return ErrPlayerNotInGame
	}

	err := t.g.Raise(idx, chipLevel)
	if err != nil {
		return err
	}

	return nil
}
