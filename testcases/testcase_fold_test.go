package pokerface

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/weedbox/pokerface"
)

func Test_Fold(t *testing.T) {

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
	err := g.Start()
	assert.Nil(t, err)

	// Waiting for ready
	for _, p := range g.GetPlayers() {
		assert.Equal(t, "Initialized", g.GetState().Status.CurrentEvent.Name)
		assert.Equal(t, 0, len(p.State().HoleCards))
		assert.Equal(t, false, p.State().Fold)
		assert.Equal(t, int64(0), p.State().Wager)
		assert.Equal(t, int64(0), p.State().Pot)
		assert.Equal(t, int64(players[p.SeatIndex()].Bankroll), p.State().Bankroll)
		assert.Equal(t, int64(players[p.SeatIndex()].Bankroll), p.State().InitialStackSize)
		assert.Equal(t, int64(players[p.SeatIndex()].Bankroll), p.State().StackSize)
		assert.Equal(t, "ready", p.State().AllowedActions[0])

		// Position checks
		if p.SeatIndex() == 0 {
			assert.True(t, p.CheckPosition("dealer"))
		} else if p.SeatIndex() == 1 {
			assert.True(t, p.CheckPosition("sb"))
		} else if p.SeatIndex() == 2 {
			assert.True(t, p.CheckPosition("bb"))
		}

		err := p.Ready()
		assert.Nil(t, err)
	}

	// ante
	assert.Equal(t, "Prepared", g.GetState().Status.CurrentEvent.Name)

	for _, p := range g.GetPlayers() {
		assert.Equal(t, false, p.State().Acted)
		assert.Equal(t, 0, len(p.State().HoleCards))
		assert.Equal(t, false, p.State().Fold)
		assert.Equal(t, int64(0), p.State().Wager)
		assert.Equal(t, int64(0), p.State().Pot)
		assert.Equal(t, int64(players[p.SeatIndex()].Bankroll), p.State().Bankroll)
		assert.Equal(t, int64(players[p.SeatIndex()].Bankroll), p.State().InitialStackSize)
		assert.Equal(t, int64(players[p.SeatIndex()].Bankroll), p.State().StackSize)
		assert.Equal(t, "pay", p.State().AllowedActions[0])
		err := p.Pay(opts.Ante)
		assert.Nil(t, err)
	}

	// Entering Preflop
	assert.Equal(t, "RoundInitialized", g.GetState().Status.CurrentEvent.Name)
	assert.Equal(t, "preflop", g.GetState().Status.Round)

	// Blinds
	for _, p := range g.GetPlayers() {
		assert.Equal(t, false, p.State().Acted)
		assert.Equal(t, "RoundInitialized", g.GetState().Status.CurrentEvent.Name)
		assert.Equal(t, 2, len(p.State().HoleCards))
		assert.Equal(t, false, p.State().Fold)
		assert.Equal(t, int64(0), p.State().Wager)
		assert.Equal(t, int64(10), p.State().Pot)

		if p.SeatIndex() == 1 {
			assert.Equal(t, "pay", p.State().AllowedActions[0])

			// Small blind
			err := p.Pay(5)
			assert.Nil(t, err)
		} else if p.SeatIndex() == 2 {
			assert.Equal(t, "pay", p.State().AllowedActions[0])

			// Big blind
			err := p.Pay(10)
			assert.Nil(t, err)
		}
	}

	// Current wager on the table should be equal to big blind
	assert.Equal(t, int64(10), g.GetState().Status.CurrentWager)
	assert.Equal(t, 2, g.GetState().Status.CurrentRaiser)

	// get ready
	for _, p := range g.GetPlayers() {
		assert.Equal(t, "ready", p.State().AllowedActions[0])
		err := p.Ready()
		assert.Nil(t, err)
	}

	// Starting player loop
	assert.Equal(t, "RoundStarted", g.GetState().Status.CurrentEvent.Name)

	// Dealer
	cp := g.GetCurrentPlayer()
	assert.Equal(t, 4, len(cp.State().AllowedActions))
	assert.Equal(t, "allin", cp.State().AllowedActions[0])
	assert.Equal(t, "fold", cp.State().AllowedActions[1])
	assert.Equal(t, "call", cp.State().AllowedActions[2])
	assert.Equal(t, "raise", cp.State().AllowedActions[3])

	// Dealer: fold
	err = cp.Fold()
	assert.Nil(t, err)

	// SB
	cp = g.GetCurrentPlayer()
	assert.Equal(t, 4, len(cp.State().AllowedActions))
	assert.Equal(t, "allin", cp.State().AllowedActions[0])
	assert.Equal(t, "fold", cp.State().AllowedActions[1])
	assert.Equal(t, "call", cp.State().AllowedActions[2])
	assert.Equal(t, "raise", cp.State().AllowedActions[3])

	// SB: fold
	err = cp.Fold()
	assert.Nil(t, err)

	// This game should be closed immediately
	err = g.Next()
	assert.Nil(t, err)
	assert.Equal(t, "GameClosed", g.GetState().Status.CurrentEvent.Name)

	//g.PrintState()
}

func Test_Fold_PassRequired(t *testing.T) {

	pf := pokerface.NewPokerFace()

	// Options
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
		&pokerface.PlayerSetting{
			Bankroll: 10000,
		},
	}
	opts.Players = append(opts.Players, players...)

	// Initializing game
	g := pf.NewGame(opts)

	// Start the game
	assert.Nil(t, g.Start())

	// Waiting for initial ready
	assert.Nil(t, g.ReadyForAll())

	// Ante
	assert.Nil(t, g.PayAnte())

	// Blinds
	assert.Nil(t, g.Pay(5))
	assert.Nil(t, g.Pay(10))

	// Round: Preflop
	assert.Nil(t, g.ReadyForAll()) // ready for the round
	assert.Equal(t, false, g.GetCurrentPlayer().CheckPosition("dealer"))
	assert.Equal(t, false, g.GetCurrentPlayer().CheckPosition("sb"))
	assert.Equal(t, false, g.GetCurrentPlayer().CheckPosition("bb"))
	assert.Nil(t, g.Call())
	assert.Nil(t, g.Call())  // Dealer
	assert.Nil(t, g.Call())  // SB
	assert.Nil(t, g.Check()) // BB

	// Round: Flop
	assert.Nil(t, g.Next())
	assert.Nil(t, g.ReadyForAll()) // ready for the round
	assert.Equal(t, true, g.GetCurrentPlayer().CheckPosition("sb"))
	assert.Nil(t, g.Check()) // SB
	assert.Nil(t, g.Check()) // BB
	assert.Nil(t, g.Bet(100))
	assert.Nil(t, g.Call()) // Dealer
	assert.Nil(t, g.Fold()) // SB
	assert.Nil(t, g.Call()) // BB

	// Round: Turn
	assert.Nil(t, g.Next())
	assert.Nil(t, g.ReadyForAll()) // ready for the round
	assert.Equal(t, true, g.GetCurrentPlayer().CheckPosition("sb"))
	assert.Nil(t, g.Pass())   // SB
	assert.Nil(t, g.Bet(100)) // BB
	assert.Nil(t, g.Raise(200))
	assert.Nil(t, g.Call()) // Dealer
	assert.Nil(t, g.Pass()) // SB
	assert.Nil(t, g.Call()) // BB

	// Round: River
	assert.Nil(t, g.Next())
	assert.Nil(t, g.ReadyForAll()) // ready for the round
	assert.Equal(t, true, g.GetCurrentPlayer().CheckPosition("sb"))
	assert.Nil(t, g.Pass())  // SB
	assert.Nil(t, g.Check()) // BB
	assert.Nil(t, g.Check())
	assert.Nil(t, g.Check()) // Dealer

	// Game closed
	assert.Nil(t, g.Next())
}
