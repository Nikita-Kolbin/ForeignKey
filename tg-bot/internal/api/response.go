package api

type Response struct {
	Status string `json:"status"`
	Error  string `json:"error"`
}

func responseOk() Response {
	return Response{Status: "Ok", Error: ""}
}

func responseError(err string) Response {
	return Response{Status: "Error", Error: err}
}
