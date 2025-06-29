package api

import (
	"fmt"
	"net/http"
)

type ctxKey string

const (
	CtxUserID ctxKey = "userId"
	CtxGroup  ctxKey = "group"
)

func GetUserIDAndGroup(r *http.Request) (string, string) {
	uid, _ := r.Context().Value(CtxUserID).(string)
	g, _ := r.Context().Value(CtxGroup).(string)
	fmt.Println("uid:", uid, "group:", g)
	return uid, g
}
