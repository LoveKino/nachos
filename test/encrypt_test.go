package nachos

import (
	"../session"
	"testing"
)

func TestEncrypt(t *testing.T) {
	text := "123"
	key := "AES256Key-ddch09acfers1234567890"
	stext, err := session.Encrypt([]byte(key), text)
	if err != nil || text == stext {
		t.Error("Got Exception: " + err.Error())
	} else {
		t.Log(stext)

		dtext, derr := session.Decrypt([]byte(key), stext)
		if derr != nil || dtext != text {
			t.Error("can not decrypt text")
		}
	}
}
