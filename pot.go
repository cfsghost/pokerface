package main

import "github.com/cfsghost/pokerface/pot"

func (g *game) updatePots() error {

	pots := pot.NewPotList()

	for _, p := range g.gs.Players {
		pots.AddContributer(p.Wager, p.Idx)
	}

	// Merge pots
	for i, pot := range pots.GetPots() {

		// More side pots
		if i > 0 {
			g.gs.Status.Pots = append(g.gs.Status.Pots, pot)
			continue
		}

		// Getting the last pot
		lastPot := g.gs.Status.Pots[len(g.gs.Status.Pots)-1]

		// Merge pot
		lastPot.Wager += pot.Wager
		lastPot.Total += pot.Total
	}

	return nil
}
