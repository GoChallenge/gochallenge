package auth

import (
	"crypto/md5"
	"fmt"
	"math/rand"
	"net/http"
	"net/url"

	"github.com/gochallenge/gochallenge/model"
	"github.com/julienschmidt/httprouter"
)

const authResultURL = "/#api_key=%s"
const authErrorURL = "/#error=%s"

// GithubInit initiates github authentication workflow
func GithubInit(gh model.GithubAPI) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request,
		_ httprouter.Params) {
		// TODO: write state into a signed cookie, and validate it in GithubVerify.
		// see http://godoc.org/golang.org/x/oauth2#Config.AuthCodeURL for details
		url := gh.AuthURL(state())
		http.RedirectHandler(url, http.StatusFound).ServeHTTP(w, r)
	}
}

// GithubVerify verifies github callback information, and inits
// user record
func GithubVerify(gh model.GithubAPI, us model.Users) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request,
		_ httprouter.Params) {
		var (
			u   *model.User
			loc string
		)

		gu, err := getGithubUser(gh, r)
		u, err = setupUser(err, us, gu)

		if err != nil {
			loc = fmt.Sprintf(authErrorURL, url.QueryEscape(err.Error()))
		} else {
			loc = fmt.Sprintf(authResultURL, u.APIKey)
		}

		http.RedirectHandler(loc, http.StatusFound).ServeHTTP(w, r)
	}
}

func state() string {
	b := make([]byte, 32)
	for i := range b {
		b[i] = byte(rand.Intn(256))
	}
	return fmt.Sprintf("%x", md5.Sum(b))
}

func getGithubUser(gh model.GithubAPI, r *http.Request) (*model.GithubUser, error) {
	r.ParseForm()
	gc, err := gh.NewClientWithToken(r.FormValue("code"))
	if err != nil {
		return nil, err
	}

	return gh.User(gc)
}

func setupUser(err error, us model.Users,
	gu *model.GithubUser) (*model.User, error) {
	var u *model.User

	if err != nil {
		return u, err
	}

	if u, err = us.FindByGithubID(gu.ID); err == nil {
		// found existing user record, just return it
		return u, nil
	}

	if err != model.ErrNotFound {
		// got an error that does not indicate a missing user - fail
		return nil, err
	}

	// Couldn't find a user - add a new one
	if u, err = model.NewUser(); err != nil {
		return u, err
	}

	gu.Populate(u)
	err = us.Add(u)

	return u, err
}
