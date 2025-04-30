package data

type PlayerInfo struct {
	PlayerID  string   `json:"playerID"`
	Character string   `json:"character"`
	Position  Position `json:"position"`
}

// ExistingPlayersMessage представляет список уже подключённых игроков
type ExistingPlayersMessage struct {
	Type    string         `json:"type"`    // "existing_players"
	Players []SpawnMessage `json:"players"` // Список игроков
}

// InputMessage представляет действие игрока (например, движение, прыжок)
type InputMessage struct {
	Type      string `json:"type"`      // "input"
	PlayerID  string `json:"playerID"`  // ID игрока
	Action    string `json:"action"`    // Например, "moveLeft", "jump"
	Pressed   bool   `json:"pressed"`   // true (зажата) или false (отпущена)
	Timestamp int64  `json:"timestamp"` // Временная метка
}

// SyncMessage представляет текущее состояние игрока
type SyncMessage struct {
	Type     string   `json:"type"`     // "sync"
	PlayerID string   `json:"playerID"` // ID игрока
	Position Position `json:"position"` // Координаты игрока
	Velocity Velocity `json:"velocity"` // Скорость игрока
}

// PlayerState представляет состояние игрока (позиция, скорость)
type PlayerState struct {
	Position Position `json:"position"` // Координаты игрока
	Velocity Velocity `json:"velocity"` // Скорость игрока
}

// Position представляет координаты игрока
type Position struct {
	X float64 `json:"x"`
	Y float64 `json:"y"`
}

// Velocity представляет скорость игрока
type Velocity struct {
	DX float64 `json:"dx"`
	DY float64 `json:"dy"`
}

// CorrectionMessage представляет исправление состояния игрока
type CorrectionMessage struct {
	Type              string   `json:"type"`              // "correction"
	PlayerID          string   `json:"playerID"`          // ID игрока
	CorrectedPosition Position `json:"correctedPosition"` // Исправленные координаты
}

// SpawnMessage представляет сообщение о подключении нового игрока
type SpawnMessage struct {
	Type      string   `json:"type"`      // "spawn"
	PlayerID  string   `json:"playerID"`  // ID игрока
	Character string   `json:"character"` // Имя персонажа
	Position  Position `json:"position"`  // Начальная позиция
}

// DisconnectMessage представляет сообщение об отключении игрока
type DisconnectMessage struct {
	Type     string `json:"type"`     // "disconnect"
	PlayerID string `json:"playerID"` // ID игрока
}
