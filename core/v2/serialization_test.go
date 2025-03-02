package v2_test

import (
	"encoding/hex"
	"math/big"
	"testing"

	"github.com/Layr-Labs/eigenda/core"
	v2 "github.com/Layr-Labs/eigenda/core/v2"
	"github.com/Layr-Labs/eigenda/encoding/utils/codec"
	"github.com/stretchr/testify/assert"
)

func TestBlobKey(t *testing.T) {
	blobKey := v2.BlobKey([32]byte{1, 2, 3})

	assert.Equal(t, "0102030000000000000000000000000000000000000000000000000000000000", blobKey.Hex())
	bk, err := v2.HexToBlobKey(blobKey.Hex())
	assert.NoError(t, err)
	assert.Equal(t, blobKey, bk)
}

func TestPaymentHash(t *testing.T) {
	pm := core.PaymentMetadata{
		AccountID:         "0x123",
		ReservationPeriod: 5,
		CumulativePayment: big.NewInt(100),
	}
	hash, err := pm.Hash()
	assert.NoError(t, err)
	// 0xf5894a8e9281b5687c0c7757d3d45fb76152bf659e6e61b1062f4c6bcb69c449 verified in solidity
	assert.Equal(t, "f5894a8e9281b5687c0c7757d3d45fb76152bf659e6e61b1062f4c6bcb69c449", hex.EncodeToString(hash[:]))
}

func TestBlobKeyFromHeader(t *testing.T) {
	data := codec.ConvertByPaddingEmptyByte(GETTYSBURG_ADDRESS_BYTES)
	commitments, err := p.GetCommitmentsForPaddedLength(data)
	if err != nil {
		t.Fatal(err)
	}

	bh := v2.BlobHeader{
		BlobVersion:     0,
		BlobCommitments: commitments,
		QuorumNumbers:   []core.QuorumID{0, 1},
		PaymentMetadata: core.PaymentMetadata{
			AccountID:         "0x123",
			ReservationPeriod: 5,
			CumulativePayment: big.NewInt(100),
		},
		Salt: 42,
	}
	blobKey, err := bh.BlobKey()
	assert.NoError(t, err)
	// 0x2bac85c7fc4c21ad02538a7eb44b120efbc64d25b1691470273f84c8cf82187a has verified in solidity  with chisel
	assert.Equal(t, "2bac85c7fc4c21ad02538a7eb44b120efbc64d25b1691470273f84c8cf82187a", blobKey.Hex())

	// same blob key should be generated for the blob header with shuffled quorum numbers
	bh2 := v2.BlobHeader{
		BlobVersion:     0,
		BlobCommitments: commitments,
		QuorumNumbers:   []core.QuorumID{1, 0},
		PaymentMetadata: core.PaymentMetadata{
			AccountID:         "0x123",
			ReservationPeriod: 5,
			CumulativePayment: big.NewInt(100),
		},
		Salt: 42,
	}

	blobKey2, err := bh2.BlobKey()
	assert.NoError(t, err)
	assert.Equal(t, blobKey2.Hex(), blobKey.Hex())
}

func TestBatchHeaderHash(t *testing.T) {
	batchRoot := [32]byte{}
	copy(batchRoot[:], []byte("1"))
	batchHeader := &v2.BatchHeader{
		ReferenceBlockNumber: 1,
		BatchRoot:            batchRoot,
	}

	hash, err := batchHeader.Hash()
	assert.NoError(t, err)
	// 0x891d0936da4627f445ef193aad63afb173409af9e775e292e4e35aff790a45e2 verified in solidity
	assert.Equal(t, "891d0936da4627f445ef193aad63afb173409af9e775e292e4e35aff790a45e2", hex.EncodeToString(hash[:]))
}

func TestBatchHeaderSerialization(t *testing.T) {
	batchRoot := [32]byte{}
	copy(batchRoot[:], []byte("batchRoot"))
	batchHeader := &v2.BatchHeader{
		ReferenceBlockNumber: 1000,
		BatchRoot:            batchRoot,
	}

	serialized, err := batchHeader.Serialize()
	assert.NoError(t, err)
	deserialized, err := v2.DeserializeBatchHeader(serialized)
	assert.NoError(t, err)
	assert.Equal(t, batchHeader, deserialized)
}

func TestBlobCertHash(t *testing.T) {
	data := codec.ConvertByPaddingEmptyByte(GETTYSBURG_ADDRESS_BYTES)
	commitments, err := p.GetCommitmentsForPaddedLength(data)
	if err != nil {
		t.Fatal(err)
	}

	blobCert := &v2.BlobCertificate{
		BlobHeader: &v2.BlobHeader{
			BlobVersion:     0,
			BlobCommitments: commitments,
			QuorumNumbers:   []core.QuorumID{0, 1},
			PaymentMetadata: core.PaymentMetadata{
				AccountID:         "0x123",
				ReservationPeriod: 5,
				CumulativePayment: big.NewInt(100),
			},
			Salt: 42,
		},
		Signature: []byte{1, 2, 3},
		RelayKeys: []v2.RelayKey{4, 5, 6},
	}

	hash, err := blobCert.Hash()
	assert.NoError(t, err)

	// afa39b4c45197f0254f7e8e2127c797c74578357e9f077eab7a8aa62e1402bec has verified in solidity with chisel
	assert.Equal(t, "afa39b4c45197f0254f7e8e2127c797c74578357e9f077eab7a8aa62e1402bec", hex.EncodeToString(hash[:]))
}

func TestBlobCertSerialization(t *testing.T) {
	data := codec.ConvertByPaddingEmptyByte(GETTYSBURG_ADDRESS_BYTES)
	commitments, err := p.GetCommitmentsForPaddedLength(data)
	if err != nil {
		t.Fatal(err)
	}

	blobCert := &v2.BlobCertificate{
		BlobHeader: &v2.BlobHeader{
			BlobVersion:     0,
			BlobCommitments: commitments,
			QuorumNumbers:   []core.QuorumID{0, 1},
			PaymentMetadata: core.PaymentMetadata{
				AccountID:         "0x123",
				ReservationPeriod: 5,
				CumulativePayment: big.NewInt(100),
			},
			Salt: 42,
		},
		Signature: []byte{1, 2, 3},
		RelayKeys: []v2.RelayKey{4, 5, 6},
	}

	serialized, err := blobCert.Serialize()
	assert.NoError(t, err)
	deserialized, err := v2.DeserializeBlobCertificate(serialized)
	assert.NoError(t, err)
	assert.Equal(t, blobCert, deserialized)
}
