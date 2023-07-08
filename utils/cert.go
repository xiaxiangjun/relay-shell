package utils

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"math/big"
	"time"
)

// 创建需要证书文件
func CreateCertificate() (string, string, error) {
	// 生成密钥
	key, err := rsa.GenerateKey(rand.Reader, 2048)
	if nil != err {
		return "", "", err
	}

	// 证书模板
	certTpl := x509.Certificate{
		SerialNumber: big.NewInt(0),
		Subject: pkix.Name{
			CommonName:         "os",
			Country:            []string{"CN"},        // 国家
			OrganizationalUnit: []string{"xxj.org"},   // 组织单位
			Province:           []string{"chongqing"}, // 省
			Locality:           []string{"nanchuan"},  // 地点
		},
		NotBefore:             time.Now(),                                                   // 开始时间
		NotAfter:              time.Now().AddDate(50, 0, 0),                                 // 过期时间
		BasicConstraintsValid: true,                                                         // 基本的有效性约束
		IsCA:                  true,                                                         // 是否根证书
		KeyUsage:              x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature, // 数字签名, 密钥加密
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth, x509.ExtKeyUsageClientAuth},
	}

	// 生成证书
	cert, err := x509.CreateCertificate(rand.Reader, &certTpl, &certTpl, &key.PublicKey, key)
	if nil != err {
		return "", "", err
	}

	keyPem := pem.EncodeToMemory(&pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: x509.MarshalPKCS1PrivateKey(key),
	})

	certPem := pem.EncodeToMemory(&pem.Block{
		Type:  "CERTIFICATE",
		Bytes: cert,
	})

	return string(certPem), string(keyPem), nil
}
