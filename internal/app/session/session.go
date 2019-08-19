package session

import (
	"generateId"
	"net/http"
	"net/url"
	"time"
	"user"
)

type Session struct {
	ID string
	UserID string
	Expiry time.Time
}

const(
	//Keep users logged in for 3 days
	SessionLength      = 24 * 3 * time.Hour
	SessionCookieName = "GophrSession"
	SessionIDLength = 20
)

func NewSession(w http.ResponseWriter) *Session {
	expiry := time.Now().Add(SessionLength)

	session := &Session{
		ID:     generateId.GenerateID("sess", SessionIDLength),
		Expiry: expiry,
	}

	cookie := http.Cookie{
		Name: SessionCookieName,
		Value: session.ID,
		Expires: session.Expiry,
	}

	http.SetCookie(w, &cookie)
	return session
}

func RequestSession(r *http.Request) *Session {
	//fmt.Println("sesion_requestsession: ",Session{})
	cookie, err := r.Cookie(SessionCookieName)

	if err != nil{
		return nil
	}

	session, err := GlobalSessionStore.Find(cookie.Value)
	if err != nil{
		panic(err)
	}

	if session == nil{
		return nil
	}

	if session.Expired() {
		GlobalSessionStore.Delete(session)
		return nil
	}
	return session
}

func (session *Session) Expired() bool {
	return session.Expiry.Before(time.Now())
}

func RequestUser(r *http.Request) *user.User {
	session := RequestSession(r)
	if session == nil || session.UserID == ""{
		return nil
	}

	user, err := user.GlobalUserStore.Find(session.UserID)

	if  err != nil{
		panic(err)
	}
	return user
}

func RequireLogin(w http.ResponseWriter, r *http.Request)  {
	// Let the request pass if we've got a users
	if RequestUser(r) != nil{
		return
	}

	query := url.Values{}
	query.Add("next", url.QueryEscape(r.URL.String()))

	http.Redirect(w, r, "/login?"+ query.Encode(), http.StatusFound)
}

func FindOrCreateSession(w http.ResponseWriter, r *http.Request) *Session {
	session := RequestSession(r)
	if session == nil {
		session = NewSession(w)
	}
	return session
}