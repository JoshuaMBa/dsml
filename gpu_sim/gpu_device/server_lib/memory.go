package server_lib

import (
	"fmt"
	"os"
	"syscall"
	"unsafe"
)

// MemorySpace interface defines the methods for a memory space
type MemorySpace interface {
	Allocate(size uint64) error
	Read(addr uint64, length uint64) ([]byte, error)
	Write(addr uint64, data []byte) error
	MinAddr() uint64
	MaxAddr() uint64
	Close() error
}

// DeviceMemory implements MemorySpace using a file-backed memory
type DeviceMemory struct {
	File    *os.File
	minAddr uint64
	maxAddr uint64
	Mapped  []byte
}

// Allocate sets up the memory-backed file and maps it
func (dm *DeviceMemory) Allocate(size uint64) error {
	if dm.File == nil {
		return fmt.Errorf("file not initialized")
	}

	// Allocate file size
	err := dm.File.Truncate(int64(size))
	if err != nil {
		return err
	}

	// Memory map the file
	mapped, err := syscall.Mmap(int(dm.File.Fd()), 0, int(size), syscall.PROT_READ|syscall.PROT_WRITE, syscall.MAP_SHARED)
	if err != nil {
		return err
	}

	dm.minAddr = uint64(uintptr(unsafe.Pointer(&mapped[0])))
	dm.maxAddr = dm.minAddr + size
	dm.Mapped = mapped
	return nil
}

// Write writes data to the specified address
func (dm *DeviceMemory) Write(addr uint64, data []byte) error {
	if addr < dm.minAddr || addr+uint64(len(data)) > dm.maxAddr {
		return fmt.Errorf("address out of range")
	}
	copy(dm.Mapped[addr-dm.minAddr:], data)
	return nil
}

// Read reads data from the specified address
func (dm *DeviceMemory) Read(addr uint64, length uint64) ([]byte, error) {
	if addr < dm.minAddr || addr+length > dm.maxAddr {
		return nil, fmt.Errorf("address out of range")
	}
	return dm.Mapped[addr-dm.minAddr : length], nil
}

// Close releases the mapped memory and closes the file
func (dm *DeviceMemory) Close() error {
	if dm.Mapped != nil {
		if err := syscall.Munmap(dm.Mapped); err != nil {
			return err
		}
	}
	if dm.File != nil {
		return dm.File.Close()
	}
	return nil
}

func (dm *DeviceMemory) MinAddr() uint64 {
	return dm.minAddr
}

func (dm *DeviceMemory) MaxAddr() uint64 {
	return dm.maxAddr
}

// NewDeviceMemory initializes a DeviceMemory instance
func NewDeviceMemory(fileName string) (MemorySpace, error) {
	file, err := os.Create(fileName)
	if err != nil {
		return nil, err
	}
	return &DeviceMemory{File: file}, nil
}
