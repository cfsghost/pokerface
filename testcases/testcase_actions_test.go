package pokerface

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/weedbox/pokerface"
)

func Test_Actions_Basic(t *testing.T) {

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
		&pokerface.PlayerSetting{
			Bankroll: 10000,
		},
		&pokerface.PlayerSetting{
			Bankroll: 10000,
		},
		&pokerface.PlayerSetting{
			Bankroll: 10000,
		},
		&pokerface.PlayerSetting{
			Bankroll: 10000,
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
	assert.Equal(t, "ReadyRequested", g.GetState().Status.CurrentEvent)
	assert.Nil(t, g.ReadyForAll())

	// Ante
	assert.Equal(t, "AnteRequested", g.GetState().Status.CurrentEvent)
	assert.Nil(t, g.PayAnte())

	// Blinds
	assert.Equal(t, "BlindsRequested", g.GetState().Status.CurrentEvent)
	assert.Nil(t, g.PayBlinds())

	// Round: Preflop
	assert.Equal(t, "ReadyRequested", g.GetState().Status.CurrentEvent)
	assert.Nil(t, g.ReadyForAll()) // ready for the round

	assert.Equal(t, "RoundStarted", g.GetState().Status.CurrentEvent)
	assert.Nil(t, g.Call())
	assert.Nil(t, g.Call())
	assert.Nil(t, g.Call())
	assert.Nil(t, g.Call())
	assert.Nil(t, g.Call())
	assert.Nil(t, g.Call())
	assert.Nil(t, g.Call())  // Dealer
	assert.Nil(t, g.Call())  // SB
	assert.Nil(t, g.Check()) // BB
	assert.Equal(t, "RoundClosed", g.GetState().Status.CurrentEvent)

	// Round: Flop
	assert.Nil(t, g.Next())
	assert.Equal(t, "ReadyRequested", g.GetState().Status.CurrentEvent)
	assert.Nil(t, g.ReadyForAll()) // ready for the round
	assert.Equal(t, "RoundStarted", g.GetState().Status.CurrentEvent)
	assert.Equal(t, true, g.GetCurrentPlayer().CheckPosition("sb"))
	assert.Nil(t, g.Check()) // SB
	assert.Nil(t, g.Check()) // BB
	assert.Nil(t, g.Bet(100))
	assert.Nil(t, g.Call())
	assert.Nil(t, g.Call())
	assert.Nil(t, g.Call())
	assert.Nil(t, g.Call())
	assert.Nil(t, g.Call())
	assert.Nil(t, g.Call()) // Dealer
	assert.Nil(t, g.Call()) // SB
	assert.Nil(t, g.Call()) // BB
	assert.Equal(t, "RoundClosed", g.GetState().Status.CurrentEvent)

	// Round: Turn
	assert.Nil(t, g.Next())
	assert.Equal(t, "ReadyRequested", g.GetState().Status.CurrentEvent)
	assert.Nil(t, g.ReadyForAll()) // ready for the round
	assert.Equal(t, "RoundStarted", g.GetState().Status.CurrentEvent)
	assert.Equal(t, true, g.GetCurrentPlayer().CheckPosition("sb"))
	assert.Nil(t, g.Check())  // SB
	assert.Nil(t, g.Bet(100)) // BB
	assert.Nil(t, g.Raise(200))
	assert.Nil(t, g.Raise(300))
	assert.Nil(t, g.Call())
	assert.Nil(t, g.Call())
	assert.Nil(t, g.Call())
	assert.Nil(t, g.Call())
	assert.Nil(t, g.Call()) // Dealer
	assert.Nil(t, g.Call()) // SB
	assert.Nil(t, g.Call()) // BB
	assert.Nil(t, g.Call())
	assert.Equal(t, "RoundClosed", g.GetState().Status.CurrentEvent)

	// Round: River
	assert.Nil(t, g.Next())
	assert.Equal(t, "ReadyRequested", g.GetState().Status.CurrentEvent)
	assert.Nil(t, g.ReadyForAll()) // ready for the round
	assert.Equal(t, "RoundStarted", g.GetState().Status.CurrentEvent)
	assert.Equal(t, true, g.GetCurrentPlayer().CheckPosition("sb"))
	assert.Nil(t, g.Check()) // SB
	assert.Nil(t, g.Check()) // BB
	assert.Nil(t, g.Check())
	assert.Nil(t, g.Check())
	assert.Nil(t, g.Check())
	assert.Nil(t, g.Check())
	assert.Nil(t, g.Check())
	assert.Nil(t, g.Check())
	assert.Nil(t, g.Check()) // Dealer
	assert.Equal(t, "RoundClosed", g.GetState().Status.CurrentEvent)

	// Game closed
	assert.Nil(t, g.Next())
	assert.Equal(t, "GameClosed", g.GetState().Status.CurrentEvent)
}

func Test_Actions_TwoPlayers(t *testing.T) {

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
			Positions: []string{"dealer", "sb"},
		},
		&pokerface.PlayerSetting{
			Bankroll:  10000,
			Positions: []string{"bb"},
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
	assert.Nil(t, g.PayBlinds())

	// Round: Preflop
	assert.Nil(t, g.ReadyForAll()) // ready for the round
	assert.Equal(t, true, g.GetCurrentPlayer().CheckPosition("sb"))
	assert.Nil(t, g.Call())  // SB, Dealer
	assert.Nil(t, g.Check()) // BB

	// Round: Flop
	assert.Nil(t, g.Next())
	assert.Nil(t, g.ReadyForAll()) // ready for the round
	assert.Equal(t, true, g.GetCurrentPlayer().CheckPosition("bb"))
	assert.Nil(t, g.Check()) // BB
	assert.Nil(t, g.Check()) // SB, Dealer

	// Round: Turn
	assert.Nil(t, g.Next())
	assert.Nil(t, g.ReadyForAll()) // ready for the round
	assert.Equal(t, true, g.GetCurrentPlayer().CheckPosition("bb"))
	assert.Nil(t, g.Check())    // BB
	assert.Nil(t, g.Bet(100))   // SB, Dealer
	assert.Nil(t, g.Raise(200)) // BB
	assert.Nil(t, g.Raise(300)) // SB, Dealer
	assert.Nil(t, g.Call())     // BB

	// Round: River
	assert.Nil(t, g.Next())
	assert.Nil(t, g.ReadyForAll()) // ready for the round
	assert.Equal(t, true, g.GetCurrentPlayer().CheckPosition("bb"))
	assert.Nil(t, g.Check()) // BB
	assert.Nil(t, g.Check()) // SB, Dealer

	// Game closed
	assert.Nil(t, g.Next())
}

func Test_Fold_RaiseInPreflop(t *testing.T) {

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
	assert.Nil(t, g.PayBlinds())

	// Round: Preflop
	assert.Nil(t, g.ReadyForAll()) // ready for the round
	assert.Equal(t, false, g.GetCurrentPlayer().CheckPosition("dealer"))
	assert.Equal(t, false, g.GetCurrentPlayer().CheckPosition("sb"))
	assert.Equal(t, false, g.GetCurrentPlayer().CheckPosition("bb"))
	assert.Nil(t, g.Raise(20))
	assert.Nil(t, g.Raise(30)) // Dealer
	assert.Nil(t, g.Call())    // SB
	assert.Nil(t, g.Call())    // BB
	assert.Nil(t, g.Call())

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

func Test_Actions_CallTo1BBInPreflop(t *testing.T) {

	pf := pokerface.NewPokerFace()

	// Options
	opts := pokerface.NewStardardGameOptions()
	opts.Blind.SB = 10
	opts.Blind.BB = 20
	opts.Ante = 0

	// Preparing deck
	opts.Deck = pokerface.NewStandardDeckCards()

	// Preparing players
	players := []*pokerface.PlayerSetting{
		&pokerface.PlayerSetting{
			Bankroll:  100,
			Positions: []string{"dealer"},
		},
		&pokerface.PlayerSetting{
			Bankroll:  200,
			Positions: []string{"sb"},
		},
		&pokerface.PlayerSetting{
			Bankroll:  15,
			Positions: []string{"bb"},
		},
		&pokerface.PlayerSetting{
			Bankroll:  150,
			Positions: []string{"ug"},
		},
	}
	opts.Players = append(opts.Players, players...)

	// Initializing game
	g := pf.NewGame(opts)

	// Start the game
	assert.Nil(t, g.Start())

	// Waiting for initial ready
	assert.Equal(t, "ReadyRequested", g.GetState().Status.CurrentEvent)
	assert.Nil(t, g.ReadyForAll())

	// Blinds
	assert.Equal(t, "BlindsRequested", g.GetState().Status.CurrentEvent)
	assert.Nil(t, g.PayBlinds())

	// Round: Preflop
	assert.Equal(t, "ReadyRequested", g.GetState().Status.CurrentEvent)
	assert.Nil(t, g.ReadyForAll()) // ready for the round

	assert.Equal(t, "RoundStarted", g.GetState().Status.CurrentEvent)
	assert.Nil(t, g.Call()) // UG action

	// ug wager should be equal to BB
	assert.Equal(t, g.Player(g.GetState().Status.LastAction.Source).State().Wager, g.GetState().Meta.Blind.BB)

	// g.PrintState()
}
