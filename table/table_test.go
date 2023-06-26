package table

import (
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_Table_Basic(t *testing.T) {

	var wg sync.WaitGroup
	wg.Add(1)

	backend := NewNativeBackend()
	opts := NewOptions()
	opts.MaxGames = 1 // Only one game

	table := NewTable(opts, WithBackend(backend))
	table.SetAnte(10)

	table.Join(0, &PlayerInfo{
		ID:       "player_1",
		Bankroll: 10000,
	})
	table.Join(1, &PlayerInfo{
		ID:       "player_2",
		Bankroll: 10000,
	})
	table.Join(2, &PlayerInfo{
		ID:       "player_3",
		Bankroll: 10000,
	})
	table.Join(3, &PlayerInfo{
		ID:       "player_4",
		Bankroll: 10000,
	})
	table.Join(4, &PlayerInfo{
		ID:       "player_5",
		Bankroll: 10000,
	})

	roundRunner := func(ts *State) {

		//t.Log(ts.GameState.Status.Round)

		switch ts.GameState.Status.Round {
		case "preflop":
			assert.Nil(t, table.Call("player_4"))
			assert.Nil(t, table.Call("player_5"))
			assert.Nil(t, table.Call("player_1"))  // Dealer
			assert.Nil(t, table.Call("player_2"))  // SB
			assert.Nil(t, table.Check("player_3")) // BB
		case "flop":
			assert.Nil(t, table.Check("player_2")) // SB
			assert.Nil(t, table.Check("player_3")) // BB
			assert.Nil(t, table.Bet("player_4", 100))
			assert.Nil(t, table.Call("player_5"))
			assert.Nil(t, table.Call("player_1")) // Dealer
			assert.Nil(t, table.Call("player_2")) // SB
			assert.Nil(t, table.Call("player_3")) // BB
		case "turn":
			assert.Nil(t, table.Check("player_2"))    // SB
			assert.Nil(t, table.Bet("player_3", 100)) // BB
			assert.Nil(t, table.Raise("player_4", 200))
			assert.Nil(t, table.Raise("player_5", 300))
			assert.Nil(t, table.Call("player_1")) // Dealer
			assert.Nil(t, table.Call("player_2")) // SB
			assert.Nil(t, table.Call("player_3")) // BB
			assert.Nil(t, table.Call("player_4"))
		case "river":
			assert.Nil(t, table.Check("player_2")) // SB
			assert.Nil(t, table.Check("player_3")) // BB
			assert.Nil(t, table.Check("player_4"))
			assert.Nil(t, table.Check("player_5"))
			assert.Nil(t, table.Check("player_1")) // Dealer
		}
	}

	roundStates := map[string]bool{
		"preflop": false,
		"flop":    false,
		"turn":    false,
		"river":   false,
	}

	table.OnStateUpdated(func(ts *State) {

		if ts.GameState == nil {
			return
		}

		assert.True(t, ts.GameState.HasPosition(0, "dealer"))
		assert.True(t, ts.GameState.HasPosition(1, "sb"))
		assert.True(t, ts.GameState.HasPosition(2, "bb"))

		//t.Log("OnStateUpdated >", ts.GameState.Status.CurrentEvent)

		switch ts.GameState.Status.CurrentEvent {
		case "ReadyRequested":
			assert.Nil(t, table.Ready("player_1"))
			assert.Nil(t, table.Ready("player_2"))
			assert.Nil(t, table.Ready("player_3"))
			assert.Nil(t, table.Ready("player_4"))
			assert.Nil(t, table.Ready("player_5"))
		case "AnteRequested":
			assert.Nil(t, table.Pay("player_1", 10))
			assert.Nil(t, table.Pay("player_2", 10))
			assert.Nil(t, table.Pay("player_3", 10))
			assert.Nil(t, table.Pay("player_4", 10))
			assert.Nil(t, table.Pay("player_5", 10))
		case "BlindsRequested":
			assert.Nil(t, table.Pay("player_2", 5))
			assert.Nil(t, table.Pay("player_3", 10))
		case "RoundStarted":

			if !roundStates[ts.GameState.Status.Round] {
				roundStates[ts.GameState.Status.Round] = true
				roundRunner(ts)

			}

		case "GameClosed":
			assert.NotNil(t, ts.GameState.Result)
			wg.Done()
		}
	})

	// Starting table
	assert.Equal(t, "idle", table.GetState().Status)
	assert.Nil(t, table.Start())

	wg.Wait()

	assert.Equal(t, "closed", table.GetState().Status)
	assert.Equal(t, opts.MaxGames, table.GetGameCount())
}
