package signer

//nolint:tagliatelle
type SignOperationResponse struct {
	PublicKey     string `json:"publicKey"`
	Signature     string `json:"signature"`
	CorrelationID string `json:"correlationId,omitempty"`
}

type Signer interface {
	Sign(string, []byte) (*SignOperationResponse, error)
}
