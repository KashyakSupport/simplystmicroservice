package hellodatastore

import (
	"net/http"

	"google.golang.org/appengine"
	"google.golang.org/appengine/log"
	"google.golang.org/appengine/memcache"
)

func getSession(req *http.Request) *memcache.Item {
	cookie, err := req.Cookie("session")

	if err != nil {
		return &memcache.Item{}
	}

	ctx := appengine.NewContext(req)
	item, err := memcache.Get(ctx, cookie.Value)
	if err != nil {
		return &memcache.Item{}
	}
	log.Infof(ctx, "item"+"%s", item)
	log.Infof(ctx, "item.Value"+"%s", string(item.Key))
	log.Infof(ctx, "item.Value"+"%s", string(item.Value))
	log.Infof(ctx, "cookie.Value"+"%s", cookie.Value)
	return item
}
