package models

type PersonalPage struct {
	Id        string `json:"id"`
	Password  string `json:"password"`
	Email     string `json:"email"`
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	Age       int    `json:"age"`
	Sex       string `json:"sex"`
	Address   string `json:"address"`
}
