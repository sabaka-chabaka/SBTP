package frame

import "crypto/sha256"

func computeChecksum(payload []byte) [ChecksumSize]byte {
	return sha256.Sum256(payload)
}

func verifyChecksum(payload []byte, expected [ChecksumSize]byte) bool {
	actual := computeChecksum(payload)
	return actual == expected
}

func (f *Frame) ApplyChecksum() {
	f.Checksum = computeChecksum(f.Payload)
	f.Flags = f.Flags.Set(FlagChecksum)
}

func (f *Frame) VerifyChecksum() error {
	if !f.Flags.Has(FlagChecksum) {
		return nil
	}
	if !verifyChecksum(f.Payload, f.Checksum) {
		return ErrChecksumMismatch
	}
	return nil
}
