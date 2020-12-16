package user

type userFilter struct {
	Name string `db:"name" form:"name"`
	Email string `db:"email" form:"email"`
	Phone string `db:"phone" form:"phone"`
	PageSize int `db:"LIMIT" form:"page_size"`
	Page int	`db:"OFFSET" form:"page"`
}

func (u *userFilter) GetPageSize() int {
	return u.PageSize
}

func (u *userFilter) GetPage() int {
	return u.Page
}