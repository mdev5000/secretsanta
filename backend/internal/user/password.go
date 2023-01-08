package user

import "golang.org/x/crypto/bcrypt"

func hashPassword(password []byte) ([]byte, error) {
	return bcrypt.GenerateFromPassword(password, bcrypt.DefaultCost)
}

func comparePassword(hashPassword []byte, password []byte) error {
	err := bcrypt.CompareHashAndPassword(hashPassword, password)
	return err
}
