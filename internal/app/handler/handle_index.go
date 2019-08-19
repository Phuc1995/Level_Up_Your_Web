package handler

import (
	"github.com/julienschmidt/httprouter"
	 img "images"
	"net/http"
)

func HandleHome(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	images, err := img.GlobalImageStore.FindAll(0)
	if err != nil{
		panic(err)
	}

	// Display Home Page
	RenderTemplate(w, r, "index/home", map[string]interface{}{
		"Images": images,
	})
}
