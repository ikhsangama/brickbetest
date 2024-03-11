package signature

func Verify(signature []string) bool {
	// should be verifying with public key, but do simplicity here
	return signature[0] == "brick"
}
