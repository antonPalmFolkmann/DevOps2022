package controllers

type UserReq struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type RegisterReq struct {
	UserReq
	Email string `json:"email"`
}

type LoginResp struct {
	Username string   `json:"username"`
	Email    string   `json:"email"`
	Avatar   string   `json:"avatar"`
	Follows  []string `json:"follows"`
}

type MsgResp struct {
	AuthorName string `json:"authorName"`
	Text       string `json:"text"`
	PubDate    string `json:"pubDate"`
	Flagged    bool   `json:"flagged"`
}

type MsgsPerUsernameResp struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Avatar   string `json:"avatar"`
	Msgs     []MsgResp
}

type AddMsgsReq struct {
	AuthorName string `json:"authorName"`
	Text       string `json:"text"`
	PubDate    string `json:"pubDate"`
}
