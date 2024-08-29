package aes

import "testing"

func TestAES(t *testing.T) {
	var text = "hello"
	var key = []byte{1, 2, 3, 4, 5, 6, 7, 8, 9, 0, 1, 2, 3, 4, 5, 6}
	b, err := Encrypt([]byte(text), key, key)
	if err != nil {
		t.Error(err)
		return
	}
	res, err := Decrypt(b, key, key)
	if err != nil {
		t.Error(err)
		return
	}
	t.Log(string(res))
	if string(res) != text {
		t.Fail()
	}
}
