package domain

type Author struct {
	ID     uint
	UserID uint
}

func NewAuthor(userID uint) *Author {
	return &Author{
		UserID: userID,
	}
}
