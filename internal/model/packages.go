package model

type Packages struct {
	ID          int      `json:"id"`
	Name        string   `json:"name"`
	Price       int      `json:"price"`
	Description string   `json:"description"`
	Benefits    []string `json:"benefits"`
}

type GetPackageByIdReq struct {
	ID int `json:"id"`
}
