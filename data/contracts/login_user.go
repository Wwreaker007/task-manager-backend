package contracts

type LoginRequest struct{
	Name 		string			`json:"name"`
	Password 	string			`json:"password"`
}

type LoginResponse struct{
	Status		string			`json:"status"`
	Token 		string			`json:"token"`		// To be cached at the user end as a cookie
}