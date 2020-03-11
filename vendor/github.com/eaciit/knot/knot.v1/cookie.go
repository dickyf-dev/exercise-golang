package knot

import (
	"net/http"
	"net/url"
	"sync"
	"time"
)

type CookieStore struct {
	sync.RWMutex
	data map[string]*http.Cookie
}

var DefaultCookieExpire time.Duration

func (cs *CookieStore) initCookies() {
	cs.Lock()
	if cs.data == nil {
		cs.data = make(map[string]*http.Cookie)
	}
	cs.Unlock()
}

func (cs *CookieStore) getCookie(r *WebContext, name string, def string) (*http.Cookie, bool) {
	cs.initCookies()

	// first search on new cookies
	cs.RLock()
	c, exist := cs.data[name]
	cs.RUnlock()

	// when not found, try to search on request cookies
	if exist == false {
		var err error
		c, err = r.Request.Cookie(name)
		if err == nil {
			exist = true
		}
	}

	// when not exist and default is set
	// put cookie with default expire time
	if exist == false && len(def) > 0 {
		if int(DefaultCookieExpire) == 0 {
			DefaultCookieExpire = 30 * 24 * time.Hour
		}

		cs.setCookie(r, name, def, DefaultCookieExpire)
	}

	return c, exist
}

func (cs *CookieStore) setCookie(r *WebContext, name string, value string, expiresAfter time.Duration) *http.Cookie {
	cs.initCookies()

	c := &http.Cookie{}
	c.Name = name
	c.Value = value
	c.Path = "/"
	u, e := url.Parse(r.Request.URL.String())
	if e == nil {
		c.Expires = time.Now().Add(expiresAfter)
		c.Domain = u.Host
	}

	cs.Lock()
	cs.data[name] = c
	cs.Unlock()

	return c
}

func (cs *CookieStore) getAllCookies() map[string]*http.Cookie {
	cs.initCookies()

	cs.RLock()
	data := cs.data
	cs.RUnlock()

	return data
}
