package datastruct

type User struct {
	Handle            string `json:"handle"`
	Email             string `json:"email"`
	Country           string `json:"country"`
	City              string `json:"city"`
	Rating            int32  `json:"rating"`
	MaxRating         int32  `json:"maxRating"`
	Avatar            string `json:"avatar"`
	Password          string `json:"password"`
	EncryptedPassword string `json:"encryptedPassword"`
}