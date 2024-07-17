package io

import (
	"crypto/rand"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type IOTestSuite struct {
	suite.Suite
}

func TestIOTestSuite(t *testing.T) {
	suite.Run(t, new(IOTestSuite))
}

func (suite *IOTestSuite) TestWriteAlignSmallFile() {
	// Create a temporary srcFile for testing
	srcFile, err := os.CreateTemp("", "512B_file")
	suite.Require().NoError(err)
	defer os.Remove(srcFile.Name())

	// Generate random data
	data := make([]byte, 512) // 512 bytes
	_, err = rand.Read(data)
	suite.Require().NoError(err)

	// Write the data to the file using IOWrite
	err = Write(srcFile, data, uint64(len(data)), 4096)
	suite.Require().NoError(err)

	// Read the written data from the srcFile
	readData := make([]byte, len(data))
	_, err = srcFile.ReadAt(readData, 0)
	suite.Require().NoError(err)

	// Assert that the written data matches the original data
	assert.Equal(suite.T(), data, readData)
}

func (suite *IOTestSuite) TestWriteUnalignSmallFile() {
	// Create a temporary srcFile for testing
	srcFile, err := os.CreateTemp("", "777B_file")
	suite.Require().NoError(err)
	defer os.Remove(srcFile.Name())

	// Generate random data
	data := make([]byte, 777) // 777 bytes
	_, err = rand.Read(data)
	suite.Require().NoError(err)

	// Write the data to the file using IOWrite
	err = Write(srcFile, data, uint64(len(data)), 4096)
	suite.Require().NoError(err)

	// Read the written data from the srcFile
	readData := make([]byte, len(data))
	_, err = srcFile.ReadAt(readData, 0)
	suite.Require().NoError(err)

	// Assert that the written data matches the original data
	assert.Equal(suite.T(), data, readData)
}

func (suite *IOTestSuite) TestWriteAlign() {
	// Create a temporary srcFile for testing
	srcFile, err := os.CreateTemp("", "4M_file")
	suite.Require().NoError(err)
	defer os.Remove(srcFile.Name())

	// Generate random data
	data := make([]byte, 4*1024*1024) // 4M
	_, err = rand.Read(data)
	suite.Require().NoError(err)

	// Write the data to the file using IOWrite
	err = Write(srcFile, data, uint64(len(data)), 4096)
	suite.Require().NoError(err)

	// Read the written data from the srcFile
	readData := make([]byte, len(data))
	_, err = srcFile.ReadAt(readData, 0)
	suite.Require().NoError(err)

	// Assert that the written data matches the original data
	assert.Equal(suite.T(), data, readData)
}

func (suite *IOTestSuite) TestWriteUnalign() {
	// Create a temporary srcFile for testing
	srcFile, err := os.CreateTemp("", "5M_file")
	suite.Require().NoError(err)
	defer os.Remove(srcFile.Name())

	// Generate random data
	data := make([]byte, 5*1024*1024) // 5M
	_, err = rand.Read(data)
	suite.Require().NoError(err)

	// Write the data to the file using IOWrite
	err = Write(srcFile, data, uint64(len(data)), 4096)
	suite.Require().NoError(err)

	// Read the written data from the srcFile
	readData := make([]byte, len(data))
	_, err = srcFile.ReadAt(readData, 0)
	suite.Require().NoError(err)

	// Assert that the written data matches the original data
	assert.Equal(suite.T(), data, readData)
}

func (suite *IOTestSuite) TestWriteUnalignError() {
	// Create a temporary srcFile for testing
	srcFile, err := os.CreateTemp("", "5M_file")
	suite.Require().NoError(err)
	defer os.Remove(srcFile.Name())

	// Generate random data
	data := make([]byte, 5*1024*1024) // 5M
	_, err = rand.Read(data)
	suite.Require().NoError(err)

	os.Setenv("HARV_FAULT", "1")
	// Write the data to the file using IOWrite
	err = Write(srcFile, data, uint64(len(data)), 4096)
	os.Setenv("HARV_FAULT", "")
	assert.Equal(suite.T(), err, ErrFaultInject)
}

func (suite *IOTestSuite) TestWriteAlignBigChunk() {
	// Create a temporary srcFile for testing
	srcFile, err := os.CreateTemp("", "128M_file")
	suite.Require().NoError(err)
	defer os.Remove(srcFile.Name())

	// Generate random data
	data := make([]byte, 128*1024*1024) // 128M
	_, err = rand.Read(data)
	suite.Require().NoError(err)

	// Write the data to the file using IOWrite
	err = Write(srcFile, data, uint64(len(data)), 4194304)
	suite.Require().NoError(err)

	// Read the written data from the srcFile
	readData := make([]byte, len(data))
	_, err = srcFile.ReadAt(readData, 0)
	suite.Require().NoError(err)

	// Assert that the written data matches the original data
	assert.Equal(suite.T(), data, readData)
}

func (suite *IOTestSuite) TestWriteUnalignBigChunk() {
	// Create a temporary srcFile for testing
	srcFile, err := os.CreateTemp("", "135M_file")
	suite.Require().NoError(err)
	defer os.Remove(srcFile.Name())

	// Generate random data
	data := make([]byte, 135*1024*1024) // 135M
	_, err = rand.Read(data)
	suite.Require().NoError(err)

	// Write the data to the file using IOWrite
	err = Write(srcFile, data, uint64(len(data)), 4194304)
	suite.Require().NoError(err)

	// Read the written data from the srcFile
	readData := make([]byte, len(data))
	_, err = srcFile.ReadAt(readData, 0)
	suite.Require().NoError(err)

	// Assert that the written data matches the original data
	assert.Equal(suite.T(), data, readData)
}

func (suite *IOTestSuite) TestCopyAlignSmallFile() {
	// Create a temporary srcFile for testing
	srcFile, err := os.CreateTemp("", "512B_file")
	suite.Require().NoError(err)
	defer os.Remove(srcFile.Name())

	// Generate random data
	data := make([]byte, 512) // 512 Bytes
	_, err = rand.Read(data)
	suite.Require().NoError(err)

	// Write the data to the file using IOWrite
	err = Write(srcFile, data, uint64(len(data)), 4096)
	suite.Require().NoError(err)

	// Read the written data from the srcFile
	readData := make([]byte, len(data))
	_, err = srcFile.ReadAt(readData, 0)
	suite.Require().NoError(err)

	// Assert that the written data matches the original data
	assert.Equal(suite.T(), data, readData)

	// Create a temporary dstFile for testing
	dstFile, err := os.CreateTemp("", "dstfile")
	suite.Require().NoError(err)
	defer os.Remove(dstFile.Name())

	// Copy the data from srcFile to dstFile using Copy
	err = Copy(srcFile, dstFile, 4096)
	suite.Require().NoError(err)

	// Read the written data from the dstFile
	dstData := make([]byte, len(data))
	_, err = dstFile.ReadAt(dstData, 0)
	suite.Require().NoError(err)
	// Assert that the written data in dstFile matches the original data
	assert.Equal(suite.T(), data, dstData)
}

func (suite *IOTestSuite) TestCopyUnalignSmallFile() {
	// Create a temporary srcFile for testing
	srcFile, err := os.CreateTemp("", "777B_file")
	suite.Require().NoError(err)
	defer os.Remove(srcFile.Name())

	// Generate random data
	data := make([]byte, 777) // 777 Bytes
	_, err = rand.Read(data)
	suite.Require().NoError(err)

	// Write the data to the file using IOWrite
	err = Write(srcFile, data, uint64(len(data)), 4096)
	suite.Require().NoError(err)

	// Read the written data from the srcFile
	readData := make([]byte, len(data))
	_, err = srcFile.ReadAt(readData, 0)
	suite.Require().NoError(err)

	// Assert that the written data matches the original data
	assert.Equal(suite.T(), data, readData)

	// Create a temporary dstFile for testing
	dstFile, err := os.CreateTemp("", "dstfile")
	suite.Require().NoError(err)
	defer os.Remove(dstFile.Name())

	// Copy the data from srcFile to dstFile using Copy
	err = Copy(srcFile, dstFile, 4096)
	suite.Require().NoError(err)

	// Read the written data from the dstFile
	dstData := make([]byte, len(data))
	_, err = dstFile.ReadAt(dstData, 0)
	suite.Require().NoError(err)
	// Assert that the written data in dstFile matches the original data
	assert.Equal(suite.T(), data, dstData)
}

func (suite *IOTestSuite) TestCopyAlign() {
	// Create a temporary srcFile for testing
	srcFile, err := os.CreateTemp("", "4M_file")
	suite.Require().NoError(err)
	defer os.Remove(srcFile.Name())

	// Generate random data
	data := make([]byte, 4*1024*1024) // 4M
	_, err = rand.Read(data)
	suite.Require().NoError(err)

	// Write the data to the file using IOWrite
	err = Write(srcFile, data, uint64(len(data)), 4096)
	suite.Require().NoError(err)

	// Read the written data from the srcFile
	readData := make([]byte, len(data))
	_, err = srcFile.ReadAt(readData, 0)
	suite.Require().NoError(err)

	// Assert that the written data matches the original data
	assert.Equal(suite.T(), data, readData)

	// Create a temporary dstFile for testing
	dstFile, err := os.CreateTemp("", "dstfile")
	suite.Require().NoError(err)
	defer os.Remove(dstFile.Name())

	// Copy the data from srcFile to dstFile using Copy
	err = Copy(srcFile, dstFile, 4096)
	suite.Require().NoError(err)

	// Read the written data from the dstFile
	dstData := make([]byte, len(data))
	_, err = dstFile.ReadAt(dstData, 0)
	suite.Require().NoError(err)
	// Assert that the written data in dstFile matches the original data
	assert.Equal(suite.T(), data, dstData)
}

func (suite *IOTestSuite) TestCopyUnalign() {
	// Create a temporary srcFile for testing
	srcFile, err := os.CreateTemp("", "5M_file")
	suite.Require().NoError(err)
	defer os.Remove(srcFile.Name())

	// Generate random data
	data := make([]byte, 5*1024*1024) // 5M
	_, err = rand.Read(data)
	suite.Require().NoError(err)

	// Write the data to the file using IOWrite
	err = Write(srcFile, data, uint64(len(data)), 4096)
	suite.Require().NoError(err)

	// Read the written data from the srcFile
	readData := make([]byte, len(data))
	_, err = srcFile.ReadAt(readData, 0)
	suite.Require().NoError(err)

	// Assert that the written data matches the original data
	assert.Equal(suite.T(), data, readData)

	// Create a temporary dstFile for testing
	dstFile, err := os.CreateTemp("", "dstfile")
	suite.Require().NoError(err)
	defer os.Remove(dstFile.Name())

	// Copy the data from srcFile to dstFile using Copy
	err = Copy(srcFile, dstFile, 4096)
	suite.Require().NoError(err)

	// Read the written data from the dstFile
	dstData := make([]byte, len(data))
	_, err = dstFile.ReadAt(dstData, 0)
	suite.Require().NoError(err)
	// Assert that the written data in dstFile matches the original data
	assert.Equal(suite.T(), data, dstData)
}

func (suite *IOTestSuite) TestCopyUnalignError() {
	// Create a temporary srcFile for testing
	srcFile, err := os.CreateTemp("", "5M_file")
	suite.Require().NoError(err)
	defer os.Remove(srcFile.Name())

	// Generate random data
	data := make([]byte, 5*1024*1024) // 5M
	_, err = rand.Read(data)
	suite.Require().NoError(err)

	// Write the data to the file using IOWrite
	err = Write(srcFile, data, uint64(len(data)), 4096)
	suite.Require().NoError(err)

	// Read the written data from the srcFile
	readData := make([]byte, len(data))
	_, err = srcFile.ReadAt(readData, 0)
	suite.Require().NoError(err)

	// Assert that the written data matches the original data
	assert.Equal(suite.T(), data, readData)

	// Create a temporary dstFile for testing
	dstFile, err := os.CreateTemp("", "dstfile")
	suite.Require().NoError(err)
	defer os.Remove(dstFile.Name())

	// Copy the data from srcFile to dstFile using Copy
	os.Setenv("HARV_FAULT", "1")
	err = Copy(srcFile, dstFile, 4096)
	os.Setenv("HARV_FAULT", "")
	assert.Equal(suite.T(), err, ErrFaultInject)
}

func (suite *IOTestSuite) TestCopyAlignBigChunk() {
	// Create a temporary srcFile for testing
	srcFile, err := os.CreateTemp("", "128M_file")
	suite.Require().NoError(err)
	defer os.Remove(srcFile.Name())

	// Generate random data
	data := make([]byte, 128*1024*1024) // 128M
	_, err = rand.Read(data)
	suite.Require().NoError(err)

	// Write the data to the file using IOWrite
	err = Write(srcFile, data, uint64(len(data)), 4194304)
	suite.Require().NoError(err)

	// Read the written data from the srcFile
	readData := make([]byte, len(data))
	_, err = srcFile.ReadAt(readData, 0)
	suite.Require().NoError(err)

	// Assert that the written data matches the original data
	assert.Equal(suite.T(), data, readData)

	// Create a temporary dstFile for testing
	dstFile, err := os.CreateTemp("", "dstfile")
	suite.Require().NoError(err)
	defer os.Remove(dstFile.Name())

	// Copy the data from srcFile to dstFile using Copy
	err = Copy(srcFile, dstFile, 4194304)
	suite.Require().NoError(err)

	// Read the written data from the dstFile
	dstData := make([]byte, len(data))
	_, err = dstFile.ReadAt(dstData, 0)
	suite.Require().NoError(err)
	// Assert that the written data in dstFile matches the original data
	assert.Equal(suite.T(), data, dstData)
}

func (suite *IOTestSuite) TestCopyUnalignBigChunk() {
	// Create a temporary srcFile for testing
	srcFile, err := os.CreateTemp("", "135M_file")
	suite.Require().NoError(err)
	defer os.Remove(srcFile.Name())

	// Generate random data
	data := make([]byte, 135*1024*1024) // 135M
	_, err = rand.Read(data)
	suite.Require().NoError(err)

	// Write the data to the file using IOWrite
	err = Write(srcFile, data, uint64(len(data)), 4194304)
	suite.Require().NoError(err)

	// Read the written data from the srcFile
	readData := make([]byte, len(data))
	_, err = srcFile.ReadAt(readData, 0)
	suite.Require().NoError(err)

	// Assert that the written data matches the original data
	assert.Equal(suite.T(), data, readData)

	// Create a temporary dstFile for testing
	dstFile, err := os.CreateTemp("", "dstfile")
	suite.Require().NoError(err)
	defer os.Remove(dstFile.Name())

	// Copy the data from srcFile to dstFile using Copy
	err = Copy(srcFile, dstFile, 4194304)
	suite.Require().NoError(err)

	// Read the written data from the dstFile
	dstData := make([]byte, len(data))
	_, err = dstFile.ReadAt(dstData, 0)
	suite.Require().NoError(err)
	// Assert that the written data in dstFile matches the original data
	assert.Equal(suite.T(), data, dstData)
}

func BenchmarkPWriteBigChunkZero(b *testing.B) {
	srcFile, _ := os.CreateTemp("", "135M_file")

	data := make([]byte, 135*1024*1024) // 135M

	b.ResetTimer()
	_ = Write(srcFile, data, uint64(len(data)), 4194304)
}

func BenchmarkPWriteBigChunk(b *testing.B) {
	srcFile, _ := os.CreateTemp("", "135M_file")

	data := make([]byte, 135*1024*1024) // 135M
	_, _ = rand.Read(data)

	b.ResetTimer()
	_ = Write(srcFile, data, uint64(len(data)), 4194304)
}

func BenchmarkPWriteSmallChunk(b *testing.B) {
	srcFile, _ := os.CreateTemp("", "135M_file")

	data := make([]byte, 135*1024*1024) // 135M
	_, _ = rand.Read(data)

	b.ResetTimer()
	_ = Write(srcFile, data, uint64(len(data)), 4096)
}
