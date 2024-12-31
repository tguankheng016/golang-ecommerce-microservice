package dtos

type CreateOrEditCategoryDto struct {
	Id   int64  `json:"id"`
	Name string `json:"name" validate:"required"`
} // @name CreateOrEditCategoryDto

type CreateCategoryDto struct {
	*CreateOrEditCategoryDto
} // @name CreateCategoryDto

type EditCategoryDto struct {
	*CreateOrEditCategoryDto
} // @name EditCategoryDto
