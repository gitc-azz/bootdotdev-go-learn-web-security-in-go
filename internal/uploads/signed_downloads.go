package uploads

import (
	"crypto/hmac"
	"crypto/sha256"
	"crypto/subtle"
	"encoding/hex"
	"fmt"
	"strconv"
	"time"
)

const signedDownloadTTL = 5 * time.Minute

func CreateSignedDownloadPath(signingKey [32]byte, fileID int64, now time.Time) string {
	expires := now.Unix() + int64(signedDownloadTTL.Seconds())
	signature := signDownload(signingKey, fileID, expires)
	return fmt.Sprintf("/files/%d/signed-download?expires=%d&signature=%s", fileID, expires, signature)
}

func VerifySignedDownload(signingKey [32]byte, fileID int64, expiresValue, signature string, now time.Time) bool {
	expires, err := strconv.ParseInt(expiresValue, 10, 64)
	if err != nil || expires <= now.Unix() {
		return false
	}
	decodedSignature, err := hex.DecodeString(signature)
	if err != nil {
		return false
	}
	expectedSignature := computeMAC(signingKey, fileID, expires)

	return subtle.ConstantTimeCompare(decodedSignature, expectedSignature) == 1
}

func computeMAC(signingKey [32]byte, fileID, expires int64) []byte {
	mac := hmac.New(sha256.New, signingKey[:])
	fmt.Fprintf(mac, "GET\n/files/%d/signed-download\n%d", fileID, expires)
	return mac.Sum(nil)
}

func signDownload(signingKey [32]byte, fileID, expires int64) string {
	ret := hex.EncodeToString(computeMAC(signingKey, fileID, expires))
	return ret
}
