package global

import "testing"

func TestGenerateJoinCode(t *testing.T) {
	t.Log(GenerateJoinCode())
}

func TestVerifyLicense(t *testing.T) {
	t.Log(VerifyLicense("ELFLR-OOOO-EAORQ-NAMCM"))
}
