package certificates

import (
	"github.com/deployment-io/deployment-runner-kit/types"
)

type UpdateCertificateDtoV1 struct {
	ID             string
	CertificateArn string
	CnameName      string
	CnameValue     string
	Verified       string
}

type UpdateCertificatesArgsV1 struct {
	types.AuthArgsV1
	Certificates []UpdateCertificateDtoV1
}

type UpdateCertificatesReplyV1 struct {
	Done bool
}
