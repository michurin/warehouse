package dto

type Message struct {
	Color      string `json:"color"`
	Message    string `json:"message"`
	Name       string `json:"name"`
	TimeStamep int64  `json:"ts"`
}

type RoomStatus struct {
	Locked bool     `json:"locked"`
	Users  []string `json:"users"`
}

type StreamMessage struct {
	Message    *Message    `json:"message,omitempty"`
	RoomStatus *RoomStatus `json:"status,omitempty"`
}
