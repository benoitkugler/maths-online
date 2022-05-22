package trivialpoursuit

import "fmt"

func (RandomGroupStrategy) kind() int { return RandomGroupStrategyGrKind }

func (rd RandomGroupStrategy) initGames(*gameSession) {}

func (rd RandomGroupStrategy) selectOrCreateGame(_ string, _ int64, session *gameSession) (GameID, error) {
	// first try to find a room with some space left
	for id, room := range session.games {
		sum := room.Summary()
		players := len(sum.Successes)
		if players < rd.MaxPlayersPerGroup && sum.PlayerTurn == nil {
			return id, nil
		}
	}

	// all rooms are full (or have already started) : create a new one
	remainingPlayers := rd.TotalPlayersNumber - session.nbPlayers()

	// we want to avoid rooms with only one student
	newRoomPlayers := rd.MaxPlayersPerGroup
	if remainingPlayers == rd.MaxPlayersPerGroup+1 {
		newRoomPlayers = remainingPlayers
	}

	ProgressLogger.Printf("Creating game room (for %d players)", newRoomPlayers)
	gameID := session.createGame(newRoomPlayers)

	return gameID, nil
}

func (FixedSizeGroupStrategy) kind() int { return FixedSizeGroupStrategyGrKind }

func (fs FixedSizeGroupStrategy) initGames(session *gameSession) {
	for _, groupSize := range fs.Groups {
		session.createGame(groupSize)
	}
}

func (fs FixedSizeGroupStrategy) selectOrCreateGame(clientGameID string, _ int64, gs *gameSession) (GameID, error) {
	gs.lock.Lock()
	defer gs.lock.Unlock()
	if _, has := gs.games[clientGameID]; !has {
		return "", fmt.Errorf("unknown game ID %s", clientGameID)
	}
	return clientGameID, nil
}
