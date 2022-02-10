package transaction

type GetStoreTransaction struct {
	Name     string `json:"name"`
	Quantity int    `json:"quantity"`
}
