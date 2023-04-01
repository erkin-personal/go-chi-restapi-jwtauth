package models

type Certificate struct {
	PrivateKey string `gorm:"size:100;not null;" json:"private_key"`
	Certificate string `gorm:"size:100;not null;" json:"certificate"`
	CertificateRequest string `gorm:"size:100;not null;" json:"certificate_request"`
	E string `gorm:"size:100;not null;" json:"e"`
	Kid string `gorm:"size:100;not null;" json:"kid"`
	Kty string `gorm:"size:100;not null;" json:"kty"`
	N string `gorm:"size:100;not null;" json:"n"`
	Ttl string `gorm:"size:100;not null;" json:"ttl"`
	//**********
	EncPublicKey string `gorm:"size:100;not null;" json:"enc_public_key"`

}

type Signature struct {
	RequestBody string `gorm:"size:100;not null;" json:"req_body"`
	E string `gorm:"size:100;not null;" json:"e"`
	N string `gorm:"size:100;not null;" json:"n"`
	Signature string `gorm:"size:100;not null;" json:"signature"`
	NowUnix string `gorm:"size:100;not null;" json:"now_unix"`
	//********
	//true/false
}


type CreateSignature struct {
	ResponseBody string `gorm:"size:100;not null;" json:"res_body"`
	PrivateKey string `gorm:"size:100;not null;" json:"private_key"`
	//*******
	Base64Signature string `gorm:"size:100;not null;" json:"base64_signature"`

}








