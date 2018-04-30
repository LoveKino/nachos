package session

import (
	"net/http"
	"time"
)

/**
 * expire:  minutes
 */
func QuickSetSession(w http.ResponseWriter, sessionKey []byte, key string, value string, expire time.Duration) error {
	if value == "" {
		return nil
	}

	v, err := Encrypt(sessionKey, value)
	if err != nil {
		return err
	}

	expiration := time.Now().Add(expire * time.Minute)

	http.SetCookie(w, &http.Cookie{
		Name:    key,
		Value:   v,
		Path:    "/",
		Expires: expiration,
	})

	return nil
}

func GetSession(r *http.Request, sessionKey []byte, key string) string {
	cookie, err := r.Cookie(key)
	if err != nil {
		return ""
	}
	sourceText, err := Decrypt(sessionKey, cookie.Value)
	if err == nil && sourceText != "" {
		return sourceText
	} else {
		return ""
	}
}
