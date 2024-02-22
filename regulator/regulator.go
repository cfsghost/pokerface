package regulator

import (
	"errors"
	"fmt"
	"math"
	"sync"
)

var (
	ErrNotFoundTable    = errors.New("regulator: Table not found")
	ErrNoAvailableTable = errors.New("regulator: No available table")
	ErrAfterRegDealline = errors.New("regulator: Can't add players after the registration deadline")
)

type CompetitionStatus int

const (
	CompetitionStatus_Pending = iota
	CompetitionStatus_Normal
	CompetitionStatus_AfterRegDeadline
)

type Regulator interface {
	GetPlayerCount() int
	GetTableCount() int
	GetTable(tableID string) *Table
	SetStatus(status CompetitionStatus)
	AddPlayers(players []string) error
	SyncState(tableID string, playerCount int) (int, []string, error)
	ReleasePlayers(tableID string, players []string) error
}

type RequestTableFn func(players []string) (string, error)
type AssignPlayersFn func(tableID string, players []string) error

type regulator struct {
	maxPlayersPerTable int
	minInitialPlayers  int
	playerCount        int
	tableCount         int
	status             CompetitionStatus
	waitingQueue       []string
	tables             map[string]*Table
	mu                 sync.RWMutex
	requestTableFn     RequestTableFn
	assignPlayersFn    AssignPlayersFn
}

type Opt func(*regulator)

func WithRequestTableFn(fn RequestTableFn) Opt {
	return func(r *regulator) {
		r.requestTableFn = fn
	}
}

func WithAssignPlayersFn(fn AssignPlayersFn) Opt {
	return func(r *regulator) {
		r.assignPlayersFn = fn
	}
}

func NewRegulator(opts ...Opt) Regulator {
	r := &regulator{
		tableCount:         0,
		playerCount:        0,
		maxPlayersPerTable: 9,
		minInitialPlayers:  6,
		status:             CompetitionStatus_Pending,
		waitingQueue:       make([]string, 0),
		tables:             make(map[string]*Table),
		requestTableFn: func(players []string) (string, error) {
			panic("Not implemented")
		},
		assignPlayersFn: func(tableID string, players []string) error {
			panic("Not implemented")
		},
	}

	for _, opt := range opts {
		opt(r)
	}

	return r
}

func (r *regulator) getLowWaterLevelTableCount() int {

	requiredTables := int(math.Ceil(float64(r.playerCount) / float64(r.maxPlayersPerTable)))
	waterLevel := int(math.Floor(float64(r.playerCount) / float64(requiredTables)))

	tableCount := 0
	for _, t := range r.tables {

		if t.PlayerCount < waterLevel {
			tableCount++
			continue
		}
	}

	return tableCount
}

func (r *regulator) calculateLowerWaterLevel() float64 {

	requiredTables := int(math.Ceil(float64(r.playerCount) / float64(r.maxPlayersPerTable)))
	waterLevel := int(math.Floor(float64(r.playerCount) / float64(requiredTables)))

	tableCount := 0
	playerCount := r.playerCount
	for _, t := range r.tables {

		if t.PlayerCount <= waterLevel {
			tableCount++
			continue
		}

		playerCount -= t.PlayerCount
	}

	return float64(playerCount) / float64(tableCount)
}

func (r *regulator) requestPlayers(count int) []string {

	players := make([]string, 0)

	for i := 0; i < count; i++ {

		if len(r.waitingQueue) == 0 {
			break
		}

		player := r.waitingQueue[0]
		r.waitingQueue = r.waitingQueue[1:]

		players = append(players, player)
	}

	return players
}

func (r *regulator) getAvailableTable() (*Table, error) {

	for _, t := range r.tables {
		if t.Required > 0 {
			return t, nil
		}
	}

	return nil, nil
}

func (r *regulator) getPlayersFromWaitingQueue(count int) []string {

	if len(r.waitingQueue) == 0 {
		return []string{}
	}

	players := make([]string, 0)

	for i := 0; i < count; i++ {

		if len(r.waitingQueue) == 0 {
			break
		}

		player := r.waitingQueue[0]
		r.waitingQueue = r.waitingQueue[1:]

		players = append(players, player)
	}

	// TODO: test only: remove this later on
	fmt.Println("[MTT#DEBUG#regulator#getPlayersFromWaitingQueue] waitingQueue:", r.waitingQueue)

	return players
}

func (r *regulator) dispatchPlayer(players []string) ([]string, error) {

	// Find a table for the players
	t, err := r.getAvailableTable()
	if err != nil {
		fmt.Println("Failed to get available table:")
		fmt.Println(err)
		return players, err
	}

	if t == nil || t.Required == 0 {
		return players, ErrNoAvailableTable
	}

	candidates := players
	var picked []string

	if t.Required >= len(candidates) {

		// Pick up players this table needs
		picked = candidates
		players = []string{}

	} else if t.Required < len(candidates) {

		// Pick up players this table needs
		picked = candidates[:t.Required]
		candidates = candidates[t.Required:]
	}

	// Assign all players to the table
	err = r.assignPlayersFn(t.ID, picked)
	if err != nil {
		fmt.Println("Failed to assign players to table:")
		fmt.Println(err)
		return players, nil
	}

	// Table has no need to wait for more players
	t.Required -= len(picked)
	t.PlayerCount += len(picked)

	return candidates, nil
}

func (r *regulator) allocateTables() error {

	requiredTables := int(math.Ceil(float64(r.playerCount) / float64(r.maxPlayersPerTable)))

	waterLevel := r.maxPlayersPerTable
	//	if r.tableCount > 0 {
	waterLevel = int(math.Floor(float64(r.playerCount) / float64(requiredTables)))
	//	}

	//	fmt.Println(waterLevel)

	//for len(r.waitingQueue) >= r.minInitialPlayers {
	for waterLevel >= r.minInitialPlayers && r.tableCount < requiredTables {

		// pull players from waiting queue
		players := r.getPlayersFromWaitingQueue(waterLevel)
		if len(players) == 0 {
			return nil
		}

		// Put players into a new table
		tableID, err := r.requestTableFn(players)
		if err != nil {
			return err
		}

		t := &Table{
			ID:          tableID,
			Required:    0,
			PlayerCount: len(players),
		}

		r.tableCount++

		if len(players) < waterLevel {
			// update table sheet
			t.Required = waterLevel - len(players)
		}

		r.tables[tableID] = t
	}

	return nil
}

func (r *regulator) breakTable(tableID string) error {

	_, ok := r.tables[tableID]
	if !ok {
		return ErrNotFoundTable
	}

	delete(r.tables, tableID)

	r.tableCount--

	return nil
}

func (r *regulator) GetPlayerCount() int {
	return r.playerCount
}

func (r *regulator) GetTableCount() int {
	return r.tableCount
}

func (r *regulator) GetTable(tableID string) *Table {
	r.mu.RLock()
	defer r.mu.RUnlock()

	if t, ok := r.tables[tableID]; ok {
		return t
	}

	return nil
}

func (r *regulator) SetStatus(status CompetitionStatus) {
	r.mu.Lock()
	defer r.mu.Unlock()

	if r.status == status {
		return
	}

	oldStatus := r.status

	r.status = status

	if oldStatus == CompetitionStatus_Pending && r.status == CompetitionStatus_Normal {
		r.drainWaitingQueue()
	}
}

func (r *regulator) AddPlayers(players []string) error {

	r.mu.Lock()
	defer r.mu.Unlock()

	if r.status == CompetitionStatus_AfterRegDeadline {
		return ErrAfterRegDealline
	}

	r.playerCount += len(players)

	return r.enterWaitingQueue(players)
}

func (r *regulator) SyncState(tableID string, playerCount int) (int, []string, error) {

	r.mu.Lock()
	defer r.mu.Unlock()

	t, ok := r.tables[tableID]
	if !ok {
		return 0, []string{}, ErrNotFoundTable
	}

	// update total player
	playerChanges := t.PlayerCount - playerCount
	r.playerCount -= playerChanges

	// Update table information
	t.PlayerCount = playerCount

	//fmt.Println(tableID, r.playerCount, -playerChanges, playerCount)

	// Figure out how many tables we need
	requiredTables := int(math.Ceil(float64(r.playerCount) / float64(r.maxPlayersPerTable)))

	//fmt.Println("Required tables:", requiredTables)

	if r.status == CompetitionStatus_AfterRegDeadline {
		// We can't add more tables after the registration deadline
		if r.playerCount <= r.maxPlayersPerTable && requiredTables < r.tableCount {

			// Break table
			err := r.breakTable(tableID)
			if err != nil {
				return 0, []string{}, err
			}

			return playerCount, []string{}, nil
		}
	}

	waterLevel := float64(r.playerCount) / float64(requiredTables)

	//fmt.Println("Water level:", waterLevel)

	if float64(playerCount) < waterLevel {

		// more than one table has low water level
		if r.getLowWaterLevelTableCount() >= 2 && requiredTables < r.tableCount {

			// Break table
			err := r.breakTable(tableID)
			if err != nil {
				return 0, []string{}, err
			}

			return playerCount, []string{}, nil
		}

		// We need more players
		count := int(math.Floor(waterLevel)) - playerCount

		// Request players
		players := r.requestPlayers(count)

		// Update table information
		stillRequired := count - len(players)
		if stillRequired > 0 {
			r.tables[tableID].Required = stillRequired
		}

		t.PlayerCount += len(players)

		return 0, players, nil
	}

	if float64(playerCount) > waterLevel {
		count := playerCount - int(math.Floor(waterLevel))

		picked := 0
		for i := 0; i < count; i++ {
			lwl := r.calculateLowerWaterLevel()

			// meet a condition
			if lwl >= math.Floor(waterLevel) {
				//fmt.Println("========= No need to move more players")
				break
			}

			picked++
			t.PlayerCount--
		}

		return picked, []string{}, nil
	}

	return 0, []string{}, nil
}

func (r *regulator) drainWaitingQueue() error {

	// First time to allocate tables
	if r.tableCount == 0 && len(r.waitingQueue) >= r.minInitialPlayers {
		return r.allocateTables()
	}

	if r.tableCount > 0 {

		var err error
		candidates := r.waitingQueue

		for len(candidates) > 0 {
			candidates, err = r.dispatchPlayer(candidates)
			if err == ErrNoAvailableTable {
				break
			}
		}

		r.waitingQueue = candidates

		// still have players
		if len(candidates) > 0 {
			return r.allocateTables()
		}
	}

	return nil
}

func (r *regulator) enterWaitingQueue(players []string) error {

	r.waitingQueue = append(r.waitingQueue, players...)

	// TODO: test only: remove this later on
	fmt.Println("[MTT#DEBUG#regulator#enterWaitingQueue] waitingQueue:", r.waitingQueue)

	if r.status == CompetitionStatus_Pending {
		return nil
	}

	return r.drainWaitingQueue()
}

func (r *regulator) ReleasePlayers(tableID string, players []string) error {

	r.mu.Lock()
	defer r.mu.Unlock()

	return r.enterWaitingQueue(players)
}
