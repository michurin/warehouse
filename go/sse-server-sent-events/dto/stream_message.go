package dto

type Message struct {
	Color      string `json:"color"`
	Message    string `json:"message"`
	Name       string `json:"name"`
	TimeStamep int64  `json:"ts"`
}

type User struct {
	Name  string `json:"name"`
	Color string `json:"color"`
}

type RoomStatus struct {
	Locked bool   `json:"locked"`
	Users  []User `json:"users"`
}

type StreamMessage struct {
	Message    *Message    `json:"message,omitempty"`
	RoomStatus *RoomStatus `json:"status,omitempty"`
}
