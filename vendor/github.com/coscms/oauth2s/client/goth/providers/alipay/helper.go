package alipay

import (
	"crypto"
	"crypto/rsa"
	"encoding/base64"
	"net/url"
	"sort"
	"strings"
)

func SignRSA2(param url.Values, privateKey *rsa.PrivateKey) (string, error) {
	return SignRSAx(param, privateKey, crypto.SHA256)
}

func SignRSA(param url.Values, privateKey *rsa.PrivateKey) (string, error) {
	return SignRSAx(param, privateKey, crypto.SHA1)
}

func SignRSAx(param url.Values, privateKey *rsa.PrivateKey, hash crypto.Hash) (string, error) {
	if param == nil {
		param = make(url.Values, 0)
	}
	pList := make([]string, 0, len(param))
	for key := range param {
		var value = strings.TrimSpace(param.Get(key))
		if len(value) > 0 {
			pList = append(pList, key+"="+value)
		}
	}
	sort.Strings(pList)
	src := strings.Join(pList, "&")
	sig, err := RSASignWithKey([]byte(src), privateKey, hash)
	if err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(sig), nil
}

func VerifySign(val url.Values, key []byte) error {
	sign, err := base64.StdEncoding.DecodeString(val.Get("sign"))
	if err != nil {
		return err
	}
	signType := val.Get("sign_type")
	keys := make([]string, 0, len(val))
	for key := range val {
		if key == `sign` || key == `sign_type` || key == `alipay_cert_sn` {
			continue
		}
		keys = append(keys, key)
	}
	sort.Strings(keys)
	pList := make([]string, 0, len(keys))
	for _, key := range keys {
		pList = append(pList, key+"="+val.Get(key))
	}
	s := strings.Join(pList, "&")
	if signType == `RSA` {
		err = RSAVerify([]byte(s), sign, key, crypto.SHA1)
	} else {
		err = RSAVerify([]byte(s), sign, key, crypto.SHA256)
	}
	return err
}

func VerifyResponseData(data []byte, signType, sign string, key []byte) error {
	signBytes, err := base64.StdEncoding.DecodeString(sign)
	if err != nil {
		return err
	}

	if signType == `RSA` {
		err = RSAVerify(data, signBytes, key, crypto.SHA1)
	} else {
		err = RSAVerify(data, signBytes, key, crypto.SHA256)
	}
	return err
}
