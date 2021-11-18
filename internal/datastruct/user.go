package datastruct

type (
	User struct {
		Handle            string `json:"handle" db:"handle"`
		Email             string `json:"email" db:"email"`
		Country           string `json:"country" db:"country"`
		City              string `json:"city" db:"city"`
		Rating            int32  `json:"rating" db:"rating"`
		MaxRating         int32  `json:"max_rating" db:"max_rating"`
		Avatar            string `json:"avatar" db:"avatar"`
		Password          string `json:"password" db:"password"`
		EncryptedPassword string `json:"encrypted_password" db:"encrypted_password"`
	}

	UserQuery struct {
		Filter string `json:"filter"`
		Page   int32  `json:"page"`
		Limit  int32  `json:"limit"`
		Offset int32  `json:"offset"`
	}
)
