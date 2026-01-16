package catalog

type Product struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Price       int    `json:"price"`
}

type ProductDocument struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Price       int    `json:"price"`
}
