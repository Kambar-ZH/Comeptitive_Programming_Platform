package datastruct

type User struct {
	Id                int32  `json:"id" db:"id"`
	Handle            string `json:"handle" db:"handle"`
	Email             string `json:"email" db:"email"`
	Country           string `json:"country" db:"country"`
	City              string `json:"city" db:"city"`
	Rating            int32  `json:"rating" db:"rating"`
	MaxRating         int32  `json:"max_rating" db:"max_rating"`
	Avatar            string `json:"avatar" db:"avatar"`
	Password          string `json:"password,omitempty" db:"password,omitempty"`
	EncryptedPassword string `json:"encrypted_password,omitempty" db:"encrypted_password,omitempty"`
}