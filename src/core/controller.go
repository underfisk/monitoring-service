package core

import (
	"fmt"
	"net/http"
)


type CoreController struct { }


func (c *CoreController) onLogPost (w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "We got a new log post to insert")
}