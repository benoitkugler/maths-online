package trivialpoursuit

func (rd RandomGroupStrategy) initGames(*gameSession) {}

func (rd RandomGroupStrategy) selectGame(_ int64, session *gameSession) (GameID, int) {
	session.lock.Lock()
	defer session.lock.Unlock()

	// first try to find a room with some space left
	for id, room := range session.games {
		sum := room.Summary()
		players := len(sum.Successes)
		if players < rd.MaxPlayersPerGroup && sum.PlayerTurn == nil {
			return id, 0
		}
	}

	// all rooms are full (or have already started) : create a new one
	remainingPlayers := rd.TotalPlayersNumber - session.nbPlayers()

	// we want to avoid rooms with only one student
	newRoomPlayers := rd.MaxPlayersPerGroup
	if remainingPlayers == rd.MaxPlayersPerGroup+1 {
		newRoomPlayers = remainingPlayers
	}
	return "", newRoomPlayers
}
