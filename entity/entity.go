package entity

type SendRequest struct {
	Channel string
	To      string
	From    string
	Payload string
	Title   string
}
