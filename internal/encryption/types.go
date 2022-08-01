package encryption

type Encrypter interface {
	//Encrypts file and returnes encrypted byte sequence
	Encrypt() ([]byte, error)
	//Encrypts file and saves it to the destination; retunrs encrypted byte sequence
	EncryptAndSave(destination string) ([]byte, error)
	//Decrypts and returns decrypted byte sequence
	Decrypt() ([]byte, error)
	//Decrypts and and saves it to a file; returns decrypted byte sequence
	DecryptAndSave(destination string) ([]byte, error)
}

func Encrypt(e Encrypter) ([]byte, error) {
	return e.Encrypt()
}

func EncryptAndSave(e Encrypter, destination string) ([]byte, error) {
	return e.EncryptAndSave(destination)
}
func Decrypt(e Encrypter) ([]byte, error) {
	return e.Decrypt()
}

func DecryptAndSave(e Encrypter, destination string) ([]byte, error) {
	return e.DecryptAndSave(destination)
}
