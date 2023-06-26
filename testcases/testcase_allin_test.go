package pokerface

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/weedbox/pokerface"
)

func Test_Allin_Basic(t *testing.T) {

	pf := pokerface.NewPokerFace()

	opts := pokerface.NewStardardGameOptions()
	opts.Ante = 10

	// Preparing deck
	opts.Deck = pokerface.NewStandardDeckCards()

	// Preparing players
	players := []*pokerface.PlayerSetting{
		&pokerface.PlayerSetting{
			Bankroll:  10000,
			Positions: []string{"dealer"},
		},
		&pokerface.PlayerSetting{
			Bankroll:  10000,
			Positions: []string{"sb"},
		},
		&pokerface.PlayerSetting{
			Bankroll:  10000,
			Positions: []string{"bb"},
		},
	}
	opts.Players = append(opts.Players, players...)

	// Initializing game
	g := pf.NewGame(opts)
	assert.Nil(t, g.Start())

	// Check players
	for _, p := range g.GetPlayers() {
		assert.Equal(t, 0, len(p.State().HoleCards))
		assert.Equal(t, false, p.State().Fold)
		assert.Equal(t, int64(0), p.State().Wager)
		assert.Equal(t, int64(0), p.State().Pot)
		assert.Equal(t, int64(players[p.SeatIndex()].Bankroll), p.State().Bankroll)
		assert.Equal(t, int64(players[p.SeatIndex()].Bankroll), p.State().InitialStackSize)
		assert.Equal(t, int64(players[p.SeatIndex()].Bankroll), p.State().StackSize)

		// Position checks
		if p.SeatIndex() == 0 {
			assert.True(t, p.CheckPosition("dealer"))
		} else if p.SeatIndex() == 1 {
			assert.True(t, p.CheckPosition("sb"))
		} else if p.SeatIndex() == 2 {
			assert.True(t, p.CheckPosition("bb"))
		}
	}

	// Waiting for ready
	assert.Equal(t, "ReadyRequested", g.GetState().Status.CurrentEvent)
	assert.Nil(t, g.ReadyForAll())

	// ante
	assert.Equal(t, "AnteRequested", g.GetState().Status.CurrentEvent)

	for _, p := range g.GetPlayers() {
		assert.Equal(t, false, p.State().Acted)
		assert.Equal(t, 0, len(p.State().HoleCards))
		assert.Equal(t, false, p.State().Fold)
		assert.Equal(t, int64(0), p.State().Wager)
		assert.Equal(t, int64(0), p.State().Pot)
		assert.Equal(t, int64(players[p.SeatIndex()].Bankroll), p.State().Bankroll)
		assert.Equal(t, int64(players[p.SeatIndex()].Bankroll), p.State().InitialStackSize)
		assert.Equal(t, int64(players[p.SeatIndex()].Bankroll), p.State().StackSize)
	}

	assert.Nil(t, g.PayAnte())

	// Entering Preflop
	assert.Equal(t, "preflop", g.GetState().Status.Round)

	// Blinds
	assert.Equal(t, "BlindsRequested", g.GetState().Status.CurrentEvent)
	for _, p := range g.GetPlayers() {
		assert.Equal(t, false, p.State().Acted)
		assert.Equal(t, 2, len(p.State().HoleCards))
		assert.Equal(t, false, p.State().Fold)
		assert.Equal(t, int64(0), p.State().Wager)
		assert.Equal(t, int64(10), p.State().Pot)
	}

	assert.Nil(t, g.PayBlinds())

	// Current wager on the table should be equal to big blind
	assert.Equal(t, int64(10), g.GetState().Status.CurrentWager)
	assert.Equal(t, 2, g.GetState().Status.CurrentRaiser)

	// Waiting for ready
	assert.Nil(t, g.ReadyForAll())

	// Starting player loop
	assert.Equal(t, "RoundStarted", g.GetState().Status.CurrentEvent)

	// Dealer
	cp := g.GetCurrentPlayer()
	assert.Equal(t, 4, len(cp.State().AllowedActions))
	assert.Equal(t, "allin", cp.State().AllowedActions[0])
	assert.Equal(t, "fold", cp.State().AllowedActions[1])
	assert.Equal(t, "call", cp.State().AllowedActions[2])
	assert.Equal(t, "raise", cp.State().AllowedActions[3])

	// Dealer: Allin
	assert.Nil(t, cp.Allin())

	// SB
	cp = g.GetCurrentPlayer()
	assert.Equal(t, 2, len(cp.State().AllowedActions))
	assert.Equal(t, "allin", cp.State().AllowedActions[0])
	assert.Equal(t, "fold", cp.State().AllowedActions[1])

	// SB: Allin
	assert.Nil(t, cp.Allin())

	// BB
	cp = g.GetCurrentPlayer()
	assert.Equal(t, 2, len(cp.State().AllowedActions))
	assert.Equal(t, "allin", cp.State().AllowedActions[0])
	assert.Equal(t, "fold", cp.State().AllowedActions[1])

	// BB: fold
	assert.Nil(t, cp.Fold())

	// flop
	assert.Nil(t, g.Next())

	// turn
	assert.Nil(t, g.Next())

	// river
	assert.Nil(t, g.Next())

	// close game
	assert.Nil(t, g.Next())
	assert.Equal(t, "GameClosed", g.GetState().Status.CurrentEvent)

	//g.PrintState()
}

func Test_Allin_NoOneCanMove(t *testing.T) {

	pf := pokerface.NewPokerFace()

	opts := pokerface.NewStardardGameOptions()
	opts.Ante = 10

	// Preparing deck
	opts.Deck = pokerface.NewStandardDeckCards()

	// Preparing players
	players := []*pokerface.PlayerSetting{
		&pokerface.PlayerSetting{
			Bankroll:  10,
			Positions: []string{"dealer"},
		},
		&pokerface.PlayerSetting{
			Bankroll:  15,
			Positions: []string{"sb"},
		},
		&pokerface.PlayerSetting{
			Bankroll:  20,
			Positions: []string{"bb"},
		},
	}
	opts.Players = append(opts.Players, players...)

	// Initializing game
	g := pf.NewGame(opts)
	assert.Nil(t, g.Start())

	// Waiting for ready
	for _, p := range g.GetPlayers() {
		assert.Equal(t, 0, len(p.State().HoleCards))
		assert.Equal(t, false, p.State().Fold)
		assert.Equal(t, int64(0), p.State().Wager)
		assert.Equal(t, int64(0), p.State().Pot)
		assert.Equal(t, int64(players[p.SeatIndex()].Bankroll), p.State().Bankroll)
		assert.Equal(t, int64(players[p.SeatIndex()].Bankroll), p.State().InitialStackSize)
		assert.Equal(t, int64(players[p.SeatIndex()].Bankroll), p.State().StackSize)

		// Position checks
		if p.SeatIndex() == 0 {
			assert.True(t, p.CheckPosition("dealer"))
		} else if p.SeatIndex() == 1 {
			assert.True(t, p.CheckPosition("sb"))
		} else if p.SeatIndex() == 2 {
			assert.True(t, p.CheckPosition("bb"))
		}
	}

	// Waiting for ready
	assert.Equal(t, "ReadyRequested", g.GetState().Status.CurrentEvent)
	assert.Nil(t, g.ReadyForAll())

	// ante
	assert.Equal(t, "AnteRequested", g.GetState().Status.CurrentEvent)

	for _, p := range g.GetPlayers() {
		assert.Equal(t, false, p.State().Acted)
		assert.Equal(t, 0, len(p.State().HoleCards))
		assert.Equal(t, false, p.State().Fold)
		assert.Equal(t, int64(0), p.State().Wager)
		assert.Equal(t, int64(0), p.State().Pot)
		assert.Equal(t, int64(players[p.SeatIndex()].Bankroll), p.State().Bankroll)
		assert.Equal(t, int64(players[p.SeatIndex()].Bankroll), p.State().InitialStackSize)
		assert.Equal(t, int64(players[p.SeatIndex()].Bankroll), p.State().StackSize)
	}

	assert.Nil(t, g.PayAnte())

	// Entering Preflop
	assert.Equal(t, "preflop", g.GetState().Status.Round)

	// Blinds
	assert.Equal(t, "BlindsRequested", g.GetState().Status.CurrentEvent)
	for _, p := range g.GetPlayers() {
		assert.Equal(t, false, p.State().Acted)
		assert.Equal(t, 2, len(p.State().HoleCards))
		assert.Equal(t, false, p.State().Fold)
		assert.Equal(t, int64(0), p.State().Wager)
		assert.Equal(t, int64(10), p.State().Pot)
	}

	assert.Nil(t, g.PayBlinds())

	// Current wager on the table should be equal to big blind
	assert.Equal(t, int64(10), g.GetState().Status.CurrentWager)
	assert.Equal(t, 2, g.GetState().Status.CurrentRaiser)

	// Waiting for ready
	assert.Nil(t, g.ReadyForAll())

	// flop
	assert.Nil(t, g.Next())

	// turn
	assert.Nil(t, g.Next())

	// river
	assert.Nil(t, g.Next())

	// close game
	assert.Nil(t, g.Next())
	assert.Equal(t, "GameClosed", g.GetState().Status.CurrentEvent)

	//g.PrintState()
}
