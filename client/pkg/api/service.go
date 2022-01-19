package api

type SturdyAPI interface {
	GetView(id string) (View, error)
	GetCodebase(id string) (Codebase, error)
	AddPublicKey(publicKey string) error
	RenewAuth() (RenewAuthResponse, error)
	GetUser() (GetUserResponse, error)
	GetIgnores(viewID string) (GetIgnoresResponse, error)
}

type View struct {
	ID                 string `json:"id"`
	UserID             string `json:"user_id"`
	CodebaseID         string `json:"codebase_id"`
	CodebaseName       string `json:"codebase_name"`
	CodebaseIsArchived bool   `json:"codebase_is_archived"`
}

type Codebase struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}
