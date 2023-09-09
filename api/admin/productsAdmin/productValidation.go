package productsadmin

import "errors"

func (a *apiProduct) ValidateProductRequest() error {
	if a.Name == "" {
		return errors.New("name not provided")
	}
	if a.Price == 0 {
		return errors.New("price not provided")
	}
	if a.Description == "" {
		return errors.New("description not provided")
	}
	if a.Image == "" {
		return errors.New("image not provided")
	}
	if a.Category == "" {
		return errors.New("category not provided")
	}
	if a.Stock == 0 {
		return errors.New("stock not provided")
	}

	return nil
}

func (a *apiProduct) ValidateProductUpdateRequest() error {
	if a.Price < 0 {
		return errors.New("price cannot be negative")
	}
	if a.Stock < 0 {
		return errors.New("stock cannot be negative")
	}
	return nil
}