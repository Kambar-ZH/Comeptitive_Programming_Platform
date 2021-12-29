package datastruct

type Tag struct {
	Id   int32  `json:"id,omitempty" db:"id,omitempty"`
	Name string `json:"name,omitempty" db:"name,omitempty"`
}
