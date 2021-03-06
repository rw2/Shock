// Package auth implements http request authentication
package auth

import (
	"errors"
	//"github.com/MG-RAST/Shock/shock-server/auth/basic"
	"github.com/MG-RAST/Shock/shock-server/auth/globus"
	"github.com/MG-RAST/Shock/shock-server/auth/mgrast"
	"github.com/MG-RAST/Shock/shock-server/conf"
	e "github.com/MG-RAST/Shock/shock-server/errors"
	"github.com/MG-RAST/Shock/shock-server/user"
)

// authCache is a
var authCache cache
var authMethods []func(string) (*user.User, error)

func Initialize() {
	authCache = cache{m: make(map[string]cacheValue)}
	authMethods = []func(string) (*user.User, error){}
	if conf.AUTH_GLOBUS_TOKEN_URL != "" && conf.AUTH_GLOBUS_PROFILE_URL != "" {
		authMethods = append(authMethods, globus.Auth)
	}
	if conf.AUTH_MGRAST_OAUTH_URL != "" {
		authMethods = append(authMethods, mgrast.Auth)
	}
}

func Authenticate(header string) (u *user.User, err error) {
	if u = authCache.lookup(header); u != nil {
		return u, nil
	} else {
		for _, auth := range authMethods {
			if u, err := auth(header); u != nil && err == nil {
				authCache.add(header, u)
				return u, nil
			}
		}
	}
	return nil, errors.New(e.InvalidAuth)
}
