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

func Test_Allin_PreviousRaiseSize(t *testing.T) {

	pf := pokerface.NewPokerFace()

	opts := pokerface.NewStardardGameOptions()
	opts.Blind.SB = 50
	opts.Blind.BB = 100
	opts.Ante = 0

	// Preparing deck
	opts.Deck = pokerface.NewStandardDeckCards()

	// Preparing players
	players := []*pokerface.PlayerSetting{
		&pokerface.PlayerSetting{
			Bankroll:  29800,
			Positions: []string{"dealer", "sb"},
		},
		&pokerface.PlayerSetting{
			Bankroll:  30200,
			Positions: []string{"bb"},
		},
	}
	opts.Players = append(opts.Players, players...)

	// Initializing game
	g := pf.NewGame(opts)
	g.GetState().Meta.Deck = []string{
		"H7", "HQ", "SQ", "H8", "C5", "H9", "H6", "S5", "S7", "D7", "D6", "C8", "D4", "H4",
		"CK", "D2", "SA", "HA", "DK", "CA", "HK", "DT", "C4", "SJ", "C3", "C2", "S3", "DJ",
		"S2", "S8", "S6", "H3", "HT", "S4", "CT", "SK", "ST", "DA", "S9", "C9", "H5", "C7",
		"CQ", "D5", "C6", "DQ", "H2", "D9", "HJ", "CJ", "D3", "D8",
	}
	g.GetState().Players[0].HoleCards = []string{"H7", "HQ"}
	g.GetState().Players[1].HoleCards = []string{"SQ", "H8"}

	assert.Nil(t, g.Start())

	// Waiting for ready
	assert.Equal(t, "ReadyRequested", g.GetState().Status.CurrentEvent)
	assert.Nil(t, g.ReadyForAll())

	// Entering Preflop
	assert.Equal(t, "preflop", g.GetState().Status.Round)

	// Blinds
	assert.Nil(t, g.PayBlinds())

	// Waiting for ready
	assert.Nil(t, g.ReadyForAll())

	// Dealer: Allin
	cp := g.GetCurrentPlayer()
	assert.Nil(t, cp.Allin())

	// BB: Allin
	cp = g.GetCurrentPlayer()
	assert.Equal(t, int64(29800), g.GetState().Status.CurrentWager)
	assert.Equal(t, int64(29800-g.GetState().Meta.Blind.BB), g.GetState().Status.PreviousRaiseSize)
	assert.ElementsMatch(t, cp.State().AllowedActions, []string{"allin", "call", "fold"})
	assert.Nil(t, cp.Allin())

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

func Test_Allin_PreflopFirstCurrentPlayer(t *testing.T) {

	pf := pokerface.NewPokerFace()

	opts := pokerface.NewStardardGameOptions()
	opts.Blind.SB = 150
	opts.Blind.BB = 300
	opts.Ante = 30

	// Preparing deck
	opts.Deck = pokerface.NewStandardDeckCards()

	// Preparing players
	players := []*pokerface.PlayerSetting{
		&pokerface.PlayerSetting{
			Bankroll:  23889,
			Positions: []string{"dealer"},
		},
		&pokerface.PlayerSetting{
			Bankroll:  8,
			Positions: []string{"sb"},
		},
		&pokerface.PlayerSetting{
			Bankroll:  176346,
			Positions: []string{"bb"},
		},
		&pokerface.PlayerSetting{
			Bankroll:  101206,
			Positions: []string{"ug"},
		},
		&pokerface.PlayerSetting{
			Bankroll:  92602,
			Positions: []string{"ug2"},
		},
		&pokerface.PlayerSetting{
			Bankroll:  84727,
			Positions: []string{"mp"},
		},
		&pokerface.PlayerSetting{
			Bankroll:  99009,
			Positions: []string{"hj"},
		},
		&pokerface.PlayerSetting{
			Bankroll:  61560,
			Positions: []string{"co"},
		},
	}
	opts.Players = append(opts.Players, players...)

	// Initializing game
	g := pf.NewGame(opts)

	assert.Nil(t, g.Start())

	// Waiting for ready
	assert.Equal(t, "ReadyRequested", g.GetState().Status.CurrentEvent)
	assert.Nil(t, g.ReadyForAll())

	// Ante
	assert.Nil(t, g.PayAnte())

	// Blinds
	assert.Nil(t, g.PayBlinds())

	// Waiting for ready
	assert.Nil(t, g.ReadyForAll())

	// Entering Preflop
	assert.Equal(t, "preflop", g.GetState().Status.Round)

	// Validate First Action Player (start from ug)
	assert.Equal(t, g.GetState().Status.CurrentPlayer, 3)

	// g.PrintState()
}

func Test_Allin_BigShortStackPlayers_AllowedActions(t *testing.T) {

	pf := pokerface.NewPokerFace()

	opts := pokerface.NewStardardGameOptions()
	opts.Blind.SB = 10
	opts.Blind.BB = 20
	opts.Ante = 0

	// Preparing deck
	opts.Deck = pokerface.NewStandardDeckCards()

	// Preparing players
	players := []*pokerface.PlayerSetting{
		&pokerface.PlayerSetting{
			Bankroll:  3005,
			Positions: []string{"dealer", "sb"},
		},
		&pokerface.PlayerSetting{
			Bankroll:  11995,
			Positions: []string{"bb"},
		},
	}
	opts.Players = append(opts.Players, players...)

	// Initializing game
	g := pf.NewGame(opts)

	assert.Nil(t, g.Start())

	// Waiting for ready
	assert.Equal(t, "ReadyRequested", g.GetState().Status.CurrentEvent)
	assert.Nil(t, g.ReadyForAll())

	// Entering Preflop
	assert.Equal(t, "preflop", g.GetState().Status.Round)

	// Blinds
	assert.Nil(t, g.PayBlinds())

	// Waiting for ready
	assert.Nil(t, g.ReadyForAll())

	// preflop movements
	assert.Nil(t, g.GetCurrentPlayer().Call())     // Dealer: Call
	assert.Nil(t, g.GetCurrentPlayer().Raise(100)) // BB: Raise 100
	assert.Nil(t, g.GetCurrentPlayer().Call())     // Dealer: Call

	// next
	assert.Nil(t, g.Next())

	// Entering Flop
	assert.Equal(t, "flop", g.GetState().Status.Round)

	// Waiting for ready
	assert.Nil(t, g.ReadyForAll())

	// flop movements
	assert.Nil(t, g.GetCurrentPlayer().Check())     // BB: Check
	assert.Nil(t, g.GetCurrentPlayer().Bet(20))     // Dealer: Bet 20
	assert.Nil(t, g.GetCurrentPlayer().Raise(4721)) // BB: Raise 4721
	assert.Nil(t, g.GetCurrentPlayer().Allin())     // Dealer: Allin

	/*
		FIXME:
		Conditions: dealer (small stack player), bb (big stack player)
		Context: In this flop case, dealer has already allin.
		Bug: Since bb is current raiser now, bb should not has any action and continue to "turn" round.
	*/
	assert.Equal(t, 1, g.GetState().Status.CurrentRaiser)
	assert.Equal(t, 1, g.GetState().Status.CurrentPlayer)
	assert.Equal(t, 0, len(g.GetCurrentPlayer().State().AllowedActions), "bb should not have any action")
	// g.PrintState()
}

func Test_Allin_BigShortStackPlayers_AllowedActions2(t *testing.T) {

	pf := pokerface.NewPokerFace()

	opts := pokerface.NewStardardGameOptions()
	opts.Blind.SB = 5
	opts.Blind.BB = 10
	opts.Ante = 0

	// Preparing deck
	opts.Deck = pokerface.NewStandardDeckCards()

	// Preparing players
	players := []*pokerface.PlayerSetting{
		&pokerface.PlayerSetting{
			Bankroll:  9005,
			Positions: []string{"dealer", "sb"},
		},
		&pokerface.PlayerSetting{
			Bankroll:  2995,
			Positions: []string{"bb"},
		},
	}
	opts.Players = append(opts.Players, players...)

	// Initializing game
	g := pf.NewGame(opts)

	assert.Nil(t, g.Start())

	// Waiting for ready
	assert.Equal(t, "ReadyRequested", g.GetState().Status.CurrentEvent)
	assert.Nil(t, g.ReadyForAll())

	// Entering Preflop
	assert.Equal(t, "preflop", g.GetState().Status.Round)

	// Blinds
	assert.Nil(t, g.PayBlinds())

	// Waiting for ready
	assert.Nil(t, g.ReadyForAll())

	// preflop movements
	assert.Nil(t, g.GetCurrentPlayer().Call())  // Dealer: Call
	assert.Nil(t, g.GetCurrentPlayer().Allin()) // BB: Allin

	/*
		FIXME:
		Conditions: dealer (big stack player), bb (small stack player)
		Context: In this preflop case, bb has already allin.
		Bug: Although dealer's bankroll is still greater than bb, dealer should only have call & fold options since there's no further players to move behind dealer's position.
	*/
	assert.Equal(t, 0, g.GetState().Status.CurrentPlayer)
	assert.Equal(t, 1, g.GetState().Status.CurrentRaiser)
	assert.ElementsMatch(t, []string{"fold", "call"}, g.GetCurrentPlayer().State().AllowedActions, "dealer should only have call & fold options")
	// g.PrintState()
}
