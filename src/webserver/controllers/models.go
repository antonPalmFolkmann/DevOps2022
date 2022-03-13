package controllers

type UserReq struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type RegisterReq struct {
	UserReq
	Email string `json:"email"`
}

type RegisterResp struct {
	Error string `json:"error"`
}

type LoginResp struct {
	UserReq
	Email   string   `json:"email"`
	Follows []string `json:"follows"`
}

type MsgResp struct {
	AuthorName string `json:"authorName"`
	Text       string `json:"text"`
	PubDate    string `json:"pubDate"`
	Flagged    bool   `json:"flagged"`
}

type MsgsPerUsernameResp struct {
	UserReq
	Email  string `json:"email"`
	Avatar string `json:"avatar"`
	Msgs   []MsgResp
}

type AddMsgsReq struct {
	AuthorName string `json:"authorName"`
	Text       string `json:"text"`
	PubDate    string `json:"pubDate"`
}
