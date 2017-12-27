// midi_windows
package main

import (
	"errors"
	"syscall"
	"unsafe"
)

var (
	libwinmm uintptr

	midiOutOpen     uintptr
	midiOutClose    uintptr
	midiOutShortMsg uintptr
)

func init() {
	libwinmm = MustLoadLibrary("winmm.dll")

	midiOutOpen = MustGetProcAddress(libwinmm, "midiOutOpen")
	midiOutClose = MustGetProcAddress(libwinmm, "midiOutClose")
	midiOutShortMsg = MustGetProcAddress(libwinmm, "midiOutShortMsg")
}

func MustLoadLibrary(name string) uintptr {
	lib, err := syscall.LoadLibrary(name)
	if err != nil {
		panic(err)
	}

	return uintptr(lib)
}

func MustGetProcAddress(lib uintptr, name string) uintptr {
	addr, err := syscall.GetProcAddress(syscall.Handle(lib), name)
	if err != nil {
		panic(err)
	}

	return uintptr(addr)
}

func MidiOutOpen(
	lphmo *uintptr,
	uDeviceID uintptr,
) bool {
	ret, _, _ := syscall.Syscall6(midiOutOpen, 5,
		uintptr(unsafe.Pointer(lphmo)),
		uintptr(uDeviceID),
		0,
		0,
		0,
		0)

	return ret != 0
}

func MidiOutClose(hmo uintptr) bool {
	ret, _, _ := syscall.Syscall(midiOutClose, 1,
		hmo,
		0,
		0)

	return ret != 0
}

func MidiOutShortMsg(hmo uintptr, dwMsg uint32) bool {
	ret, _, _ := syscall.Syscall(midiOutShortMsg, 2,
		hmo,
		uintptr(dwMsg),
		0)

	return ret != 0
}

// Midi impl:
type midiImpl struct {
	hDevice uintptr
}

// Constructor:
func NewMidi() (midi Midi, err error) {
	impl := &midiImpl{}

	// Open device 1:
	if !MidiOutOpen(&impl.hDevice, 1) {
		err = syscall.GetLastError()
		return
	}

	midi = impl
	return
}

func (m *midiImpl) Close() error {
	if !MidiOutClose(m.hDevice) {
		return syscall.GetLastError()
	}
	return nil
}

func (m *midiImpl) CC(channel uint8, controller uint8, value uint8) error {
	msg := uint32(0xB0|channel) | uint32(controller<<8) | uint32(value<<16)
	if !MidiOutShortMsg(m.hDevice, msg) {
		return syscall.GetLastError()
	}
	return nil
}

func (m *midiImpl) PC(channel uint8, program uint8) error {
	msg := uint32(0xC0|channel) | uint32(program<<8)
	if !MidiOutShortMsg(m.hDevice, msg) {
		return syscall.GetLastError()
	}
	return nil
}

func (m *midiImpl) Sysex(data []byte) error {
	panic(errors.New("midi sysex not implemented yet!"))
}
