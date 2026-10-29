package main

import (
	"bytes"
	"compress/gzip"
	"crypto"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"flag"
	"fmt"
	"io"
	"math/big"
	"os"

	"github.com/sirupsen/logrus"
)

// openssl x509 -in certificate.cer -outform pem -out certificate.pem

// The UIDAI Secure QR Code data components are delimited by the byte 255 (or -1 as signed byte)
const Delimiter byte = 255

// AadhaarData holds the decoded structural data from the Secure QR code
type AadhaarData struct {
	ReferenceID string
	Name        string
	DOB         string
	Gender      string
	CareOf      string
	District    string
	Landmark    string
	House       string
	Location    string
	Pincode     string
	PostOffice  string
	State       string
	Street      string
	SubDistrict string
	VTC         string
	RawPhoto    []byte
	Signature   []byte
}

var (
	qrcode = flag.String("qrcode", "", "aadhar qr code")
)

func main() {
	flag.Parse()
	// 1. Example raw text scanned from an Aadhaar Secure QR Code (BigInt string)
	rawQrStr := *qrcode
	// 2. Decode the Secure QR Data
	aadhaarData, signedBytes, err := DecodeAadhaarQR(rawQrStr)
	if err != nil {
		fmt.Printf("❌ Failed to decode QR payload: %v\n", err)
		return
	}

	fmt.Printf("✅ Decoded Profile: %s (%s, DOB: %s)\n", aadhaarData.Name, aadhaarData.Gender, aadhaarData.DOB)

	// 3. Verify digital signature using UIDAI Public Certificate
	// You must download the official UIDAI public V3 certificate from uidai.gov.in
	certPath := "uidai_auth_prod_2028.cer.pem" //  "uidai_12_06_18_cer.pem"
	isValid, err := VerifyAadhaarSignature(signedBytes, aadhaarData.Signature, certPath)
	if err != nil {
		fmt.Printf("❌ Signature verification error: %v\n", err)
		return
	}

	if isValid {
		fmt.Println("🔒 SUCCESS: The QR Code is authentic and tamper-proof!")
	} else {
		fmt.Println("🚨 WARNING: QR Code signature validation failed! Data could be forged.")
	}
}

// DecodeAadhaarQR parses the numeric string format into structural text fields and text blocks
func DecodeAadhaarQR(qrStr string) (*AadhaarData, []byte, error) {
	// Convert base-10 string back into big.Int bytes
	bigInt := new(big.Int)
	bigInt, ok := bigInt.SetString(qrStr, 10)
	if !ok {
		return nil, nil, errors.New("invalid base-10 integer string format")
	}
	compressedBytes := bigInt.Bytes()

	// Decompress via GZIP
	gzipReader, err := gzip.NewReader(bytes.NewReader(compressedBytes))
	if err != nil {
		return nil, nil, fmt.Errorf("gzip decompression failed: %w", err)
	}
	defer gzipReader.Close()

	decompressedBytes, err := io.ReadAll(gzipReader)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to read decompressed stream: %w", err)
	}

	// Split data blocks using the delimiter byte (255)
	parts := bytes.Split(decompressedBytes, []byte{Delimiter})
	if len(parts) < 16 {
		return nil, nil, fmt.Errorf("malformed secure QR structure, expected at least 16 fragments, got %d", len(parts))
	}

	logrus.Infof("Num Parts %d", len(parts))
	// Structural extraction matching the UIDAI standard layout
	for i, part := range parts {
		logrus.Infof("Item %d : %s", i, string(part))
	}
	data := &AadhaarData{
		ReferenceID: string(parts[0]),
		Name:        string(parts[1]),
		DOB:         string(parts[2]),
		Gender:      string(parts[3]),
		CareOf:      string(parts[4]),
		District:    string(parts[5]),
		Landmark:    string(parts[6]),
		House:       string(parts[7]),
		Location:    string(parts[8]),
		Pincode:     string(parts[9]),
		PostOffice:  string(parts[10]),
		State:       string(parts[11]),
		Street:      string(parts[12]),
		SubDistrict: string(parts[13]),
		VTC:         string(parts[14]),
	}

	// The remaining components depend on whether a profile photo is baked into the template
	// Usually: Parts[15] is the raw ISO/JPEG image data, Parts[16] is the 256-byte signature block
	if len(parts) == 17 {
		data.RawPhoto = parts[15]
		data.Signature = parts[16]
	} else {
		data.Signature = parts[39]
	}
	logrus.Infof("Data  %v", data)
	// Reconstruct the exact sequence of signed bytes (excluding the signature itself)
	// UIDAI signs the exact raw content payload up to the point of signature boundaries
	signatureOffset := bytes.LastIndex(decompressedBytes, []byte{Delimiter})
	signedBytes := decompressedBytes[:signatureOffset]

	return data, signedBytes, nil
}

// VerifyAadhaarSignature matches the content array against the signature component
func VerifyAadhaarSignature(signedBytes, signature []byte, certPath string) (bool, error) {
	certData, err := os.ReadFile(certPath)
	if err != nil {
		return false, fmt.Errorf("unable to read public certificate file: %w", err)
	}

	block, _ := pem.Decode(certData)
	if block == nil {
		return false, errors.New("failed to parse certificate PEM block")
	}

	cert, err := x509.ParseCertificate(block.Bytes)
	if err != nil {
		return false, fmt.Errorf("failed to parse x509 certificate structure: %w", err)
	}

	pubKey, ok := cert.PublicKey.(*rsa.PublicKey)
	if !ok {
		return false, errors.New("certificate public key is not valid RSA type")
	}

	// Hash the content using SHA-256
	hashed := sha256.Sum256(signedBytes)

	// Validate RSA signature blocks
	err = rsa.VerifyPKCS1v15(pubKey, crypto.SHA256, hashed[:], signature)
	if err != nil {
		return false, nil // Signature is invalid
	}

	return true, nil // Signature verified perfectly
}
