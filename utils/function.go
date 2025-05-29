package utils

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"errors"
	"fmt"
)

func SignatureHeader(secret string, timeStamp string) (string, error) {
	s := secret + "|" + timeStamp
	hash := sha256.New()
	hash.Write([]byte(s))
	signature, err := GenerateSignature(hash.Sum(nil))
	if err != nil {
		return "", err
	}
	return signature, err
}

func GenerateSignature(data []byte) (string, error) {
	signer, err := parsePrivateKey(rsa_private_key)
	if err != nil {
		return "", err
	}

	signed, err := signer.Sign(data)
	if err != nil {
		return "", err
	}

	sig := base64.StdEncoding.EncodeToString(signed)

	return sig, nil
}

func parsePrivateKey(pemBytes []byte) (Signer, error) {
	block, _ := pem.Decode(pemBytes)
	if block == nil {
		return nil, errors.New("ssh: no key found")
	}

	var rawkey interface{}
	switch block.Type {
	case "RSA PRIVATE KEY":
		rsa, err := x509.ParsePKCS1PrivateKey(block.Bytes)
		if err != nil {
			return nil, err
		}
		rawkey = rsa
	default:
		return nil, fmt.Errorf("ssh: unsupported key type %q", block.Type)
	}
	return newSignerFromKey(rawkey)
}

// A Signer is can create signatures that verify against a public key.
type Signer interface {
	// Sign returns raw signature for the given data. This method
	// will apply the hash specified for the keytype to the data.
	Sign(data []byte) ([]byte, error)
}

func newSignerFromKey(k interface{}) (Signer, error) {
	var sshKey Signer
	switch t := k.(type) {
	case *rsa.PrivateKey:
		sshKey = &rsaPrivateKey{t}
	default:
		return nil, fmt.Errorf("ssh: unsupported key type %T", k)
	}
	return sshKey, nil
}

type rsaPrivateKey struct {
	*rsa.PrivateKey
}

// Sign signs data with rsa-sha256
func (r *rsaPrivateKey) Sign(data []byte) ([]byte, error) {
	h := sha256.New()
	h.Write(data)
	d := h.Sum(nil)
	return rsa.SignPKCS1v15(rand.Reader, r.PrivateKey, crypto.SHA256, d)
}

var rsa_private_key = []byte(`-----BEGIN RSA PRIVATE KEY-----
MIIEpQIBAAKCAQEArw/cukpkM/dok1FtOanOxlvB8d8QVrM22pfQuVjRsRQNtGOP
+ce+vz0N6YULYUKXlGUH7dRiS2CZIEnIFIl8BoZ6QcSPZPAakMo56/9MN8MUCqj6
BXHq+E+LX0cZb4aXsCwbH0Y/I3VFjVGYNLqfDO8HQ455uMYQM2P2PiG2EjVDcCyQ
t4UnBvVlrpLDXtaAfsx2EoVX5iodnKZ6J1bujRxwDPfHOneaNBbx+INQ6jqnopRG
b13vol2fz/kOspREF+BO0TTT3nLYKvmepcW3t/PitgJEnT7Tb+eBBgE/2C53Juqr
btk4Lb3+jssC1har4StLiWci9J4L9Ytt9k+UEwIDAQABAoIBAQCCNx5MQ4F9Vg6n
Ze4E8lYoHaCJtQ6GLxAiUMKk23g+a1g2UciKVxV/Un7CsH/ifJIbg3r+YPgscVH1
PmxhOlLS17ygpwCyaBaalJG5BVFAOQ7zTvWKj03kHebhnBVDa63xER2riaj1SwnN
MGAy+I3OJQ4WJihMVKdAUp6bYJ/5sgFrpv+INDja4K0b4ELIqJpAgRV7Vc90kMaW
1GHFDu8W6yU76yq+3q9eu2Raos38s67JloNGYQYxAuQtWOyM4AGKYYGvAJYkALHN
bmDMVdwdpawg+5RqeRhEhg82IcmnnHtzB6/pufxvxytuMuWWAaRyz9pqlmmsJLRP
yDbWxP1hAoGBAN+VvQ6RAOaplVO8r2QwHDUKed/9TySyGaFnXbV+GIMk4iLoPc00
GJ0yRSLRB7J2xvmKC55IVFZdK+s2Blg3KnllAcFY9bbT/nVkXmyexeQtCflV8L5/
ZIikZsphkj4DwVuShc9I/HI3i1l6Ml/oGIhQJhgX74HQ8x2G5jklGuvdAoGBAMhx
NuFqtA41CSW6yEKWqye3pnlIDMhVW5CmJ6N8LqoO62bsN99g65yVCTpQ7p+b4S/d
OiHFwdbZYr+LUngbGQ2XGrwELTB34YYDP+RFGZ6MD0WeoDXSKQIbD659sX4Cheu1
pqgA7IQ0JYUvqSX11P8OQaZVSz4PyJblzz8dkDivAoGBAJDaNxjn0rid57PPi6YS
EUQ/3EPEnfC9PiO2jxyVbBYS4DsTUW7PsJ6vQeFToXP6xeBHkk1iuNkXFewWHTgr
zWXGjcOQ+egQIkw10YL3vmec0lhqWEVizWRFdp7pZdCtqCjGndB0jbEF0U8P/vDp
snMl0fhMEYx+LfPUQPWG15E1AoGBALLSbkfEzkYugq6qaKcfjCqu6VIiOWUw4bO2
yH5N98O388Oq0l3zNcBIJidktL6obsoo8AfZSgnHfxWr0jNc2YkKWcuLXlVzXjwV
AhdAno6YHbfawMvDZtp+Egt2D7d/wMJ9GOWhjWCUtTSRRLKdEx1JNsCSL8J6ilY1
SCPi2Bv5AoGAS8pDXVlz6TYa8E+zOdyKh1/JVLIo+CJ9Bk4X80j487wR1LsJl4cy
LPK0TdEeemnchtsDW6IVTraSFWNz/8F2BK/IxVksCU5NgBukbUKEgbVRkdwZQiSr
/JZC9jfQL6aHwrdgddSoR4MhDx4aO0kFdcfhKHmPDDcSNZWhB+VmoMs=
-----END RSA PRIVATE KEY-----`)
