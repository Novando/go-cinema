package uuid

import (
	"encoding/hex"
	"fmt"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"strings"
)

func ParseUUID(src string) (res pgtype.UUID, err error) {
	var dst [16]byte
	switch len(src) {
	case 36:
		src = src[0:8] + src[9:13] + src[14:18] + src[19:23] + src[24:]
	case 32:
		// dashes already stripped, assume valid
	default:
		// assume invalid.
		return pgtype.UUID{Valid: false}, fmt.Errorf("INVALID_UUID:%v", src)
	}

	buf, err := hex.DecodeString(src)
	if err != nil {
		return pgtype.UUID{Valid: false}, fmt.Errorf("INVALID_UUID")
	}

	copy(dst[:], buf)
	return pgtype.UUID{Bytes: dst, Valid: true}, err
}

func CheckNullUUID(uuidStr string) bool {
	return strings.Contains(uuidStr, "0000000000000000000000000000")
}

func GenerateUUID() pgtype.UUID {
	newUuid := uuid.New()
	generated, _ := ParseUUID(newUuid.String())
	return generated
}

func FormatDashedUUID(src [16]byte) string {
	res := fmt.Sprintf("%x", src)
	if len(res) == 32 {
		res = fmt.Sprintf("%v-%v-%v-%v-%v", res[0:8], res[9:13], res[14:18], res[19:23], res[24:])
	}
	return res
}
