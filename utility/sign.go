package utility

import (
	"bytes"
	"encoding/hex"
	"fmt"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
)

func VerifyEIP191Signature(address string, message string, signature string) (bool, error) {
	// Step 1: Format the message with Ethereum's prefix
	prefixedMessage := fmt.Sprintf("\x19Ethereum Signed Message:\n%d%s", len(message), message)

	// Step 2: Hash the prefixed message
	messageHash := crypto.Keccak256([]byte(prefixedMessage))

	// Step 3: Decode the signature
	signatureBytes, err := hex.DecodeString(signature[2:]) // Remove "0x" prefix
	if err != nil {
		return false, fmt.Errorf("invalid signature format: %v", err)
	}
	if len(signatureBytes) != 65 {
		return false, fmt.Errorf("invalid signature length")
	}

	// Extract r, s, and v from the signature
	//r := signatureBytes[:32]
	//s := signatureBytes[32:64]
	v := signatureBytes[64]

	// Adjust v value for recovery (27 or 28 are valid)
	if v < 27 {
		v += 27
	}

	// Append v back into the signature
	fullSignature := append(signatureBytes[:64], v-27) // Remove 27 offset for go-ethereum

	// Step 4: Recover the public key
	publicKey, err := crypto.SigToPub(messageHash, fullSignature)
	if err != nil {
		return false, fmt.Errorf("failed to recover public key: %v", err)
	}

	// Step 5: Derive the address from the public key
	recoveredAddress := crypto.PubkeyToAddress(*publicKey)

	// Step 6: Compare the recovered address with the expected address
	if bytes.Equal(recoveredAddress.Bytes(), common.HexToAddress(address).Bytes()) {
		return true, nil
	}
	return false, fmt.Errorf("address mismatch: got %s, expected %s", recoveredAddress.Hex(), address)
}
