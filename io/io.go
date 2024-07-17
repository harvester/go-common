package io

/*
#include <fcntl.h>
#include <stdlib.h>
#include <unistd.h>
#include <string.h>
#include <errno.h>

// Safe_pwrite is a wrapper around pwrite(2) that retries on EINTR.
ssize_t safe_pwrite(int fd, const void *buf, size_t count, off_t offset)
{
        while (count > 0) {
                ssize_t r = pwrite(fd, buf, count, offset);
                if (r < 0) {
                        if (errno == EINTR)
                                continue;
                        return -errno;
                }
                count -= r;
                buf = (char *)buf + r;
                offset += r;
        }
        return 0;
}
ssize_t safe_pread(int fd, void *buf, size_t count, off_t offset)
{
        size_t cnt = 0;
        char *b = (char*)buf;

        while (cnt < count) {
                ssize_t r = pread(fd, b + cnt, count - cnt, offset + cnt);
                if (r <= 0) {
                        if (r == 0) {
                                // EOF
                                return cnt;
                        }
                        if (errno == EINTR)
                                continue;
                        return -errno;
                }
                cnt += r;
        }
        return cnt;
}
ssize_t safe_pread_exact(int fd, void *buf, size_t count, off_t offset)
{
        ssize_t ret = safe_pread(fd, buf, count, offset);
        if (ret < 0)
                return ret;
        if ((size_t)ret != count)
                return -EDOM;
        return 0;
}
// Write data to a file descriptor with O_DIRECT
int directWrite(int fd, void *buf, size_t count, off_t offset) {
    return safe_pwrite(fd, buf, count, offset);
}
// Read data from a file descriptor with O_DIRECT
int directRead(int fd, void *buf, size_t count, off_t offset) {
	return safe_pread(fd, buf, count, offset);
}
*/
import "C"

import (
	"errors"
	"fmt"
	"os"
	"reflect"
	"sync"
	"syscall"
	"unsafe"
)

var ErrFaultInject = errors.New("fault injection")

const (
	baseAlignSize  = 4096
	maxChunkSize   = 4194304
	maxProducerNum = 8
	BLKGETSIZE64   = 0x80081272
)

type Content struct {
	offset uint64
	buf    []byte
}

func Copy(src *os.File, dst *os.File, chunkSize int) error {
	var ioError error
	var producerNum = maxProducerNum
	var faultInject = error(nil)
	if os.Getenv("HARV_FAULT") != "" {
		faultInject = ErrFaultInject
	}
	ioProducer := func(start, boundary uint64, src *os.File, ioQ chan<- Content, ioWG, producerWG *sync.WaitGroup, errChan chan error, completedChan chan struct{}) {
		var startOffset uint64
		startOffset = start
		defer producerWG.Done()
		for {
			select {
			case _, open := <-completedChan:
				if !open {
					return
				}
			default:
				if startOffset >= boundary {
					return
				}
				end := startOffset + uint64(chunkSize)
				if end > boundary {
					end = boundary
				}
				count := end - startOffset
				buf := make([]byte, count)
				_, err := PReadExact(src, buf, int(count), startOffset)
				if err != nil || faultInject != nil {
					if err == nil && faultInject != nil {
						err = faultInject
					}
					errChan <- err
					return
				}
				// check empty zero buffer
				emptyBuf := make([]byte, count)
				if !reflect.DeepEqual(buf, emptyBuf) {
					ioWG.Add(1)
					ioQ <- Content{offset: uint64(startOffset), buf: buf}
				}
				startOffset = end
			}
		}
	}

	srcSize, err := getSourceVolSize(src)
	if err != nil {
		return fmt.Errorf("error getting file size")
	}

	if chunkSize > maxChunkSize {
		return fmt.Errorf("chunk size is too large, max chunk size is %d", maxChunkSize)
	}
	if chunkSize%baseAlignSize != 0 {
		return fmt.Errorf("chunk size must be a multiple of %d", baseAlignSize)
	}

	// Calculate the number of chunks based on the chunk size
	numChunks := srcSize / uint64(chunkSize)
	if int(srcSize)%chunkSize != 0 {
		numChunks++
	}
	if int(numChunks) < producerNum {
		producerNum = int(numChunks)
	}

	// Create a channel to receive the results
	ioQ := make(chan Content, producerNum*2)
	completedChan := make(chan struct{})
	errChan := make(chan error, producerNum*2)
	defer close(errChan)

	var ioWG, producerWG, workerWG sync.WaitGroup

	for id := 0; id < producerNum; id++ {
		respSize := srcSize / uint64(producerNum)
		mod := srcSize % uint64(producerNum)
		if id == producerNum-1 {
			respSize += mod
		}
		startOffset := uint64(id) * respSize
		boundary := startOffset + respSize
		producerWG.Add(1)
		go ioProducer(startOffset, boundary, src, ioQ, &ioWG, &producerWG, errChan, completedChan)
	}
	for i := 0; i < producerNum; i++ {
		workerWG.Add(1)
		go ioWorker(dst, ioQ, &ioWG, &workerWG, errChan, completedChan, faultInject)
	}

	go func() {
		producerWG.Wait()
		ioWG.Wait()
		if ioError == nil {
			close(ioQ)
			close(completedChan)
		}
	}()

	select {
	case <-completedChan:
		break
	case errVal := <-errChan:
		ioError = errVal
		close(completedChan)
		break
	}

	workerWG.Wait()
	if ioError != nil {
		go ioQflusher(ioQ)
		producerWG.Wait()
		close(ioQ)
		return ioError
	}
	return nil

}

func Write(dst *os.File, data []byte, size uint64, chunkSize int) error {
	var ioError error
	var producerNum = maxProducerNum
	var faultInject = error(nil)
	if os.Getenv("HARV_FAULT") != "" {
		faultInject = ErrFaultInject
	}
	ioProducer := func(start, boundary uint64, ioQ chan<- Content, ioWG, producerWG *sync.WaitGroup, completedChan chan struct{}) {
		//var finalIOQueued = false
		var startOffset uint64
		startOffset = start
		defer producerWG.Done()
		for {
			select {
			case _, open := <-completedChan:
				if !open {
					return
				}
			default:
				if startOffset >= boundary {
					return
				}
				end := startOffset + uint64(chunkSize)
				if end > boundary {
					end = boundary
				}
				chunk := data[startOffset:end]
				count := end - startOffset
				emptyBuf := make([]byte, count)
				if !reflect.DeepEqual(chunk, emptyBuf) {
					ioWG.Add(1)
					ioQ <- Content{offset: uint64(startOffset), buf: chunk}
				}
				startOffset = end
			}
		}
	}

	if chunkSize > maxChunkSize {
		return fmt.Errorf("chunk size is too large, max chunk size is %d", maxChunkSize)
	}
	if chunkSize%baseAlignSize != 0 {
		return fmt.Errorf("chunk size must be a multiple of %d", baseAlignSize)
	}

	// Calculate the number of chunks based on the alignment size
	numChunks := int(size) / chunkSize
	if len(data)%chunkSize != 0 {
		numChunks++
	}
	if int(numChunks) < producerNum {
		producerNum = int(numChunks)
	}

	// Create a channel to receive the results
	ioQ := make(chan Content, producerNum*2)
	completedChan := make(chan struct{})
	errChan := make(chan error, producerNum)

	var ioWG, producerWG, workerWG sync.WaitGroup

	for id := 0; id < producerNum; id++ {
		respSize := size / uint64(producerNum)
		mod := size % uint64(producerNum)
		if id == producerNum-1 {
			respSize += mod
		}
		startOffset := uint64(id) * respSize
		boundary := startOffset + respSize
		producerWG.Add(1)
		go ioProducer(startOffset, boundary, ioQ, &ioWG, &producerWG, completedChan)
	}
	for i := 0; i < producerNum; i++ {
		workerWG.Add(1)
		go ioWorker(dst, ioQ, &ioWG, &workerWG, errChan, completedChan, faultInject)
	}

	go func() {
		producerWG.Wait()
		ioWG.Wait()
		close(ioQ)
		if ioError == nil {
			close(completedChan)
		}
	}()

	select {
	case <-completedChan:
		break
	case errVal := <-errChan:
		ioError = errVal
		close(completedChan)
		break
	}

	// all workers should be done
	workerWG.Wait()
	close(errChan)

	if ioError != nil {
		go ioQflusher(ioQ)
		producerWG.Wait()
		close(ioQ)
		return ioError
	}
	return nil
}

func PWrite(dst *os.File, data []byte, size int, offset uint64) (int, error) {
	var writeBuffer unsafe.Pointer
	if C.posix_memalign((*unsafe.Pointer)(unsafe.Pointer(&writeBuffer)), C.size_t(baseAlignSize), C.size_t(size)) != 0 {
		fmt.Printf("Error allocating aligned memory\n")
		return 0, fmt.Errorf("error allocating aligned memory")
	}
	defer C.free(unsafe.Pointer(writeBuffer))

	// Copy the Go data into the C buffer
	C.memcpy(writeBuffer, unsafe.Pointer(&data[0]), C.size_t(size))

	// Call the C function to write with O_DIRECT
	ret := C.directWrite(C.int(dst.Fd()), writeBuffer, C.size_t(size), C.off_t(offset))
	if ret < 0 {
		fmt.Printf("Error writing data: %v\n", ret)
		return 0, fmt.Errorf("error writing data")
	}

	return int(ret), nil
}

func PReadExact(src *os.File, buf []byte, count int, offset uint64) (int, error) {
	var readBuffer unsafe.Pointer
	if C.posix_memalign((*unsafe.Pointer)(unsafe.Pointer(&readBuffer)), C.size_t(baseAlignSize), C.size_t(count)) != 0 {
		fmt.Printf("Error allocating aligned memory\n")
		return 0, fmt.Errorf("error allocating aligned memory")
	}
	defer C.free(unsafe.Pointer(readBuffer))

	// Call the C function to read with O_DIRECT
	ret := C.directRead(C.int(src.Fd()), readBuffer, C.size_t(count), C.off_t(offset))
	if ret < 0 {
		fmt.Printf("Error reading data: %v\n", ret)
		return 0, fmt.Errorf("error reading data")
	}

	// Copy the C data into the Go buffer
	C.memcpy(unsafe.Pointer(&buf[0]), readBuffer, C.size_t(ret))

	return int(ret), nil
}

func getSourceVolSize(src *os.File) (uint64, error) {

	var srcSize uint64
	srcInfo, err := src.Stat()
	if err != nil {
		return 0, err
	}

	if srcInfo.Mode().IsRegular() {
		// file size should not be negative, directly return as uint64
		return uint64(srcInfo.Size()), nil
	}

	if (srcInfo.Mode() & os.ModeDevice) != 0 {
		_, _, err := syscall.Syscall(
			syscall.SYS_IOCTL,
			src.Fd(),
			BLKGETSIZE64,
			uintptr(unsafe.Pointer(&srcSize)),
		)
		if err != 0 {
			return 0, fmt.Errorf("error getting file size: %v", err)
		}
		return srcSize, nil
	}

	return 0, fmt.Errorf("unsupported file type: %v", srcInfo.Mode())
}

func ioQflusher(ioQueue chan Content) {
	for {
		_, got := <-ioQueue
		if !got {
			return
		}
	}
}

func ioWorker(dst *os.File, ioQueue chan Content, ioWG, workerWG *sync.WaitGroup, errChan chan error, completedChan chan struct{}, faultInject error) {
	defer workerWG.Done()
	for {
		select {
		case _, open := <-completedChan:
			if !open {
				return
			}
		case obj, got := <-ioQueue:
			if !got {
				return
			}
			_, err := PWrite(dst, obj.buf, len(obj.buf), obj.offset)
			ioWG.Done()
			if err != nil || faultInject != nil {
				if err == nil {
					err = faultInject
				}
				errChan <- err
				return
			}
		}
	}
}
