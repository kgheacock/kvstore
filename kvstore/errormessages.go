package kvstore

//Data is the json repsonse
type Data struct {
	Value string `json:"value"`
}

//PutSuccess handles adds and replaces
type PutSuccess struct {
	Message  string `json:"message"`
	Replaced bool   `json:"replaced"`
}

//PutFailure handles keylength error and value missing
//Handles errors for new value and replace
type PutFailure struct {
	Error   string `json:"error"`
	Message string `json:"message"`
}

//GetSuccess handles a succesful get
type GetSuccess struct {
	Exists  bool   `json:"doesExist"`
	Message string `json:"message"`
	Value   string `json:"value"`
}

//GetFailure handles a failed get
type GetFailure struct {
	Exists  bool   `json:"doesExist"`
	Error   string `json:"error"`
	Message string `json:"message"`
}

//DeleteSuccess handles a succesful delete
type DeleteSuccess struct {
	Exists  bool   `json:"doesExist"`
	Message string `json:"message"`
}

//DeleteFailure handles a failed delete
type DeleteFailure struct {
	Exists  bool   `json:"doesExist"`
	Error   string `json:"error"`
	Message string `json:"message"`
}
