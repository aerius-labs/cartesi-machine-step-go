package cartesi_machine_step

import (
	"bytes"
	"encoding/binary"
	"fmt"
)

type AccessType int

const (
	READ AccessType = iota
	WRITE
)

type Access struct {
	Position uint64
	Val      [8]byte
	Access   AccessType
}

type AccessLogs struct {
	Logs    []Access
	Current int
}

type IMemoryAccessLog interface {
	ReadWord(address uint64) uint64
	WriteWord(a AccessLogs, address uint64, value uint64)
}

func (a *AccessLogs) readWord(address uint64) (uint64, error) {
	readVal, err := a.accessManager(address, READ)
	if err != nil {
		return 0, err
	}
	return binary.LittleEndian.Uint64(readVal[:]), nil
}

func (a *AccessLogs) writeWord(address uint64, val uint64) error {
	var bytesVal [8]byte
	binary.LittleEndian.PutUint64(bytesVal[:], val)
	expected, err := a.accessManager(address, WRITE)
	if err != nil {
		return err
	}
	if !bytes.Equal(expected[:], bytesVal[:]) {
		return fmt.Errorf("written value mismatch")
	}
	return nil
}

func (a *AccessLogs) accessManager(addr uint64, accessType AccessType) ([8]byte, error) {
	if a.Current >= len(a.Logs) {
		return [8]byte{}, fmt.Errorf("too many accesses")
	}

	access := a.Logs[a.Current]

	if access.Access != accessType {
		return [8]byte{}, fmt.Errorf("access type mismatch")
	}

	if access.Position != addr {
		return [8]byte{}, fmt.Errorf("position and access address mismatch")
	}

	return access.Val, nil
}
