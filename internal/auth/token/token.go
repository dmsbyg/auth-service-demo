package token

type MakerVerifier interface {
	Maker
	Verifier
}

type Maker interface {
	Make(userID, userEmail string) (tokenString string, err error)
}

type Verifier interface {
	Verify(tokenString string) (payload JwtClaims, err error)
}
