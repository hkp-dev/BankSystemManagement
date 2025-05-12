package user

type RegisterRequest struct {
	Identifier  string `json:"identifier"`
	Email       string `json:"email"`
	Password    string `json:"password"`
	FirstName   string `json:"firstName"`
	LastName    string `json:"lastName"`
	PhoneNumber string `json:"phoneNumber"`
	Address     string `json:"address"`
	City        string `json:"city"`
	ZipCode     string `json:"zipCode"`
	IDCardFront string `json:"idCardFront"`
	IDCardBack  string `json:"idCardBack"`
}

type LoginRequest struct {
	Identifier string `json:"identifier"`
	Password   string `json:"password"`
}

type ReqGetUsrInfo struct {
	UserId string `json:"userId"`
}

