package entity

type Infos struct {
	FieldID    string `json:"fieldID,omitempty"`
	FieldValue string `json:"fieldValue,omitempty"`
	Pepper     string `json:"pepper,omitempty"`
}

type CallBackData struct {
	ID     int64 `json:"omitempty"`
	Url    string
	Data   []byte
	User   string
	Fields []string
}
