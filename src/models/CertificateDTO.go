package models

type Certificate struct {
	cert_id         int    `json:"cert_info_id" gorm:"primary_key;auto_increment;not null"`
	Target          string `json:"url_id" gorm:"not null"`
	Subject         string `json:"subject"`
	Organization    string `json:"issuer"`
	PreDate         string `json:"valid_since"`
	PostDate        string `json:"valid_until"`
	CertificateLife string `json:"certificate_life"`
	Protocol        string `json:"protocol"`
}
