package handler

import (
	error2 "error"
	"github.com/julienschmidt/httprouter"
	images2 "images"
	"net/http"
	s "session"
	user2 "user"
)

func HandlerUserNew(w http.ResponseWriter, r *http.Request, _ httprouter.Params){
	RenderTemplate(w, r, "users/new",nil)
}

func HandleUserCreate(w http.ResponseWriter, r * http.Request,_ httprouter.Params){
	user, err := user2.NewUser(
		r.FormValue("username"),
		r.FormValue("email"),
		r.FormValue("password"),
		)
	//fmt.Print(err)
	if err != nil {
		if error2.IsValidationError(err){
			RenderTemplate(w, r,"users/new", map[string]interface{}{
				"Error" : err.Error(),
				"User" : user,
			} )

			return
		}
		panic(err)
		return
	}
	err = user2.GlobalUserStore.Save(user)
	if err != nil {
		panic(err)
		return
	}

	//create a new sesssion
	session := s.NewSession(w)
	session.UserID = user.ID
	err = s.GlobalSessionStore.Save(session)
	if err != nil {
		panic(err)
		return
	}
	http.Redirect(w, r, "/?flash=User+created", http.StatusFound)
}

func HandleUserEdit(w http.ResponseWriter, r *http.Request, _ httprouter.Params)  {
	user := s.RequestUser(r)
	RenderTemplate(w, r, "users/edit",
		map[string]interface{}{
			"User": user,
		})
}

func HandleUserUpdate(w http.ResponseWriter, r *http.Request, _ httprouter.Params)  {
	currentUser := s.RequestUser(r)
	email := r.FormValue("email")
	currentPassword := r.FormValue("currentPassword")
	newPassword := r.FormValue("newPassword")
	user, err := user2.UpdateUser(currentUser, email, currentPassword, newPassword)

	if err != nil{
		if error2.IsValidationError(err){
			RenderTemplate(w, r, "users/edit", map[string]interface{}{
				"Error": err.Error(),
				"User" : user,
			})
			return
		}
		panic(err)
	}

	err = user2.GlobalUserStore.Save(*currentUser)
	if err != nil {
		panic(err)
	}

	http.Redirect(w, r, "/account?flash=User+updated", http.StatusFound)

}

func HandleUserShow(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	user, err := user2.GlobalUserStore.Find(params.ByName("userID"))
	if err != nil {
		panic(err)
	}

	// 404
	if user == nil {
		http.NotFound(w, r)
		return
	}

	images, err := images2.GlobalImageStore.FindAllByUser(user, 0)
	if err != nil {
		panic(err)
	}

	RenderTemplate(w, r, "users/show", map[string]interface{}{
		"Images": images,
		"User":   user,
	})
}
