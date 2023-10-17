package server

import (
	"context"
	"net/http"

	"github.com/google/uuid"
	"github.com/osintami/fingerprintz/common"
)

const (
	COOKIE_ID        = "osintami"
	COOKIE_DOMAIN    = "osintami.com"
	ONE_YEAR_SECONDS = 60 * 60 * 24 * 365
)

func (x *GatewayServer) PixelFireHandler(w http.ResponseWriter, r *http.Request) {
	var cookieId string
	cookie, err := r.Cookie(COOKIE_ID)
	if err != nil {
		cookieId = uuid.NewString()
	} else {
		cookieId = cookie.Value
	}
	cookie = &http.Cookie{}
	cookie.Domain = COOKIE_DOMAIN
	cookie.HttpOnly = true
	cookie.MaxAge = ONE_YEAR_SECONDS
	cookie.Name = COOKIE_ID
	cookie.Path = "/"
	cookie.SameSite = http.SameSiteNoneMode
	cookie.Secure = true
	cookie.Value = cookieId

	http.SetCookie(w, cookie)

	go x.pixels.PixelFire(context.Background(),
		&Pixel{
			CookieID:  cookieId,
			IpAddr:    common.IpAddr(r),
			UserAgent: r.UserAgent(),
			Referrer:  r.Referer(),
			Count:     1,
		})
}
