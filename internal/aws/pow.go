package aws

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"golang.org/x/crypto/scrypt"
	"math/bits"
	"strconv"
)

func checkDigest(digest []byte, difficulty int) bool {
	full := difficulty / 8
	rem := difficulty % 8
	
	for i := 0; i < full; i++ {
		if digest[i] != 0x00 {
			return false
		}
	}
	if rem != 0 && bits.LeadingZeros8(digest[full]) < rem {
		return false
	}
	return true
}

func HashPoW(challengeInput, checksum string, difficulty int) (string, error) {
	combined := []byte(challengeInput + checksum)
	
	for nonce := 0; ; nonce++ {
		data := append(combined, []byte(strconv.Itoa(nonce))...)
		digest := sha256.Sum256(data)
		if checkDigest(digest[:], difficulty) {
			return strconv.Itoa(nonce), nil
		}
	}
}

func ScryptFunc(inputStr, saltStr string, memoryCost int) (string, error) {
	hash, err := scrypt.Key([]byte(inputStr), []byte(saltStr), memoryCost, 8, 1, 16)
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(hash), nil
}

func ComputeScryptNonce(challengeInput, checksum string, difficulty int) (string, error) {
	combined := challengeInput + checksum
	salt := checksum
	memory := 128
	
	for nonce := 0; ; nonce++ {
		input := fmt.Sprintf("%s%d", combined, nonce)
		resultHex, err := ScryptFunc(input, salt, memory)
		if err != nil {
			return "", err
		}
		result, err := hex.DecodeString(resultHex)
		if err != nil {
			return "", err
		}
		if checkDigest(result, difficulty) {
			return strconv.Itoa(nonce), nil
		}
	}
}

/*
var ChallengeTypes = map[string]interface{}{
	"h72f957df656e80ba55f5d8ce2e8c7ccb59687dba3bfb273d54b08a261b2f3002": ComputeScryptNonce,
	"h7b0c470f0cfe3a80a9e26526ad185f484f6817d0832712a4a37a908786a6a67f": HashPoW,
	"ha9faaffd31b4d5ede2a2e19d2d7fd525f66fee61911511960dcbb52d3c48ce25": "mp_verify",
}*/

func SolveChallenge(challengeType, challengeInput, checksum string, difficulty int) (string, error) {
	if challengeType == "h72f957df656e80ba55f5d8ce2e8c7ccb59687dba3bfb273d54b08a261b2f3002" {
		return ComputeScryptNonce(challengeInput, checksum, difficulty)
	} else if challengeType == "h7b0c470f0cfe3a80a9e26526ad185f484f6817d0832712a4a37a908786a6a67f" {
		return HashPoW(challengeInput, checksum, difficulty)
	} else {
		return "", fmt.Errorf("unknown challengeType: %s", challengeType)
	}
}
