package datastruct

type (
	User struct {
		Id                int32  `json:"id" db:"id"`
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

	UserAllRequest struct {
		Filter string
		Page   int32
		Limit  int32
		Offset int32
	}
)
