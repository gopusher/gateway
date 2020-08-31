package helper

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
)

type Signature struct {
	secret []byte
}

func NewSignature(secret []byte) *Signature {
	return &Signature{
		secret: secret,
	}
}

func (s *Signature) Sign(plain []byte) (string, error) {
	h := hmac.New(sha256.New, s.secret)
	_, err := h.Write(plain)
	if err != nil {
		return "", err
	}

	//return fmt.Sprintf("%x", h.Sum(nil)), nil
	//return hex.EncodeToString(h.Sum(nil)), nil
	//上边两行相同，但是和下边一行意义不同

	return base64.StdEncoding.EncodeToString(h.Sum(nil)), nil
}
