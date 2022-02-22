package storage

type User struct {
	user_id   int64
	username  string
	email     string
	pw_hash   string
	follows   []User
	followers []User
}

func (u *User) FindAllMessages() {

}

func (u *User) FindMyMessages() {

}

func (u *User) GetAvatar() {

}

type Message struct {
	message_id int64
	authour_id int64
	text       string
	pub_date   int64
	flagged    bool
}
