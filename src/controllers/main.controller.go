package controllers

type MainController struct {}

func NewMainController() *MainController {
	return &MainController{}
}
/*
func ( mc *MainController) Root (req  *http.Request, w http.ResponseWriter) {
	fmt.Fprintf(w, "Hello, %q", html.EscapeString(req.URL.Path))
}*/