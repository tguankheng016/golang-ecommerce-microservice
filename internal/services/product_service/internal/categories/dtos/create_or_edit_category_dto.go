package dtos

type CreateOrEditCategoryDto struct {
	Id   *int   `json:"id"`
	Name string `json:"name"`
}

type CreateCategoryDto struct {
	CreateOrEditCategoryDto
}

type EditCategoryDto struct {
	CreateOrEditCategoryDto
}
