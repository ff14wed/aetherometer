//go:build windows
// +build windows

package win32

import (
	"errors"
	"fmt"
	"path/filepath"
	"sync"
	"syscall"
	"unsafe"

	"github.com/ff14wed/xivnet/v3"
	"golang.org/x/sys/windows"
)

type ffxivOnce struct {
	sync.Once
	handle                                      syscall.Handle
	offsetOodleNetwork1UDP_State_Size           uintptr
	offsetOodleNetwork1_Shared_Size             uintptr
	offsetOodleNetwork1_Shared_SetWindow        uintptr
	offsetOodleNetwork1UDP_Train_State_Counting uintptr
	offsetOodleNetwork1UDP_Decode               uintptr
	offset_match_set_from_histo_normalized      uintptr
	offset_nomatch_set_from_histo_normalized    uintptr

	err error
}

type OodleFactory struct {
	ffxivOnce
}

func (f *OodleFactory) New(processID uint32) (xivnet.OodleImpl, error) {
	f.ffxivOnce.Do(func() {
		processHandle, err := windows.OpenProcess(windows.PROCESS_QUERY_LIMITED_INFORMATION, false, processID)
		if err != nil {
			f.ffxivOnce.err = fmt.Errorf("open process pid %d: %s", processID, err)
			return
		}
		ffxivPathChars := make([]uint16, windows.MAX_PATH)
		var pathLen uint32 = windows.MAX_PATH
		err = windows.QueryFullProcessImageName(processHandle, 0, &ffxivPathChars[0], &pathLen)
		if err != nil {
			f.ffxivOnce.err = err
			return
		}
		_ = windows.CloseHandle(processHandle)
		ffxivPath := syscall.UTF16ToString(ffxivPathChars[:pathLen])
		ffxivDX11Path := filepath.Join(filepath.Dir(ffxivPath), "ffxiv_dx11.exe")
		ffxivHandle, err := syscall.LoadLibrary(ffxivDX11Path)
		if err != nil {
			f.ffxivOnce.err = err
			return
		}
		f.ffxivOnce.handle = ffxivHandle
		f.ffxivOnce.offsetOodleNetwork1UDP_State_Size = 0x15eda10
		f.ffxivOnce.offsetOodleNetwork1_Shared_Size = 0x15ef390
		f.ffxivOnce.offsetOodleNetwork1_Shared_SetWindow = 0x15ef260
		f.ffxivOnce.offsetOodleNetwork1UDP_Train_State_Counting = 0x15e8260
		f.ffxivOnce.offsetOodleNetwork1UDP_Decode = 0x15ed370
		f.ffxivOnce.offset_match_set_from_histo_normalized = 0x15ecc80
		f.ffxivOnce.offset_nomatch_set_from_histo_normalized = 0x15ec890
	})
	if f.ffxivOnce.err != nil {
		return nil, f.ffxivOnce.err
	}

	stateSize, err := callWrapper(f.handle, f.ffxivOnce.offsetOodleNetwork1UDP_State_Size)
	if err != nil {
		return nil, err
	}

	sharedDictSize, err := callWrapper(f.handle, f.ffxivOnce.offsetOodleNetwork1_Shared_Size, 0x13)
	if err != nil {
		return nil, err
	}

	d := &decompressor{
		ffxivOnce:  &f.ffxivOnce,
		state:      make([]byte, stateSize),
		sharedDict: make([]byte, sharedDictSize),
		initDict:   make([]byte, 0x8000),
	}

	_, err = callWrapper(
		f.handle,
		f.ffxivOnce.offsetOodleNetwork1_Shared_SetWindow,
		uintptr(unsafe.Pointer(&d.sharedDict[0])),
		0x13,
		uintptr(unsafe.Pointer(&d.initDict[0])),
		uintptr(len(d.initDict)),
	)

	if err != nil {
		return nil, err
	}

	if err := d.Train(); err != nil {
		return nil, err
	}

	return d, nil
}

type decompressor struct {
	*ffxivOnce
	state      []byte
	sharedDict []byte
	initDict   []byte
}

func (d *decompressor) Train() error {
	countingState := make([]byte, 0x4A8400)
	_, err := callWrapper(
		d.handle,
		d.ffxivOnce.offsetOodleNetwork1UDP_Train_State_Counting,
		uintptr(unsafe.Pointer(&countingState[0])),
		uintptr(unsafe.Pointer(&d.sharedDict[0])),
		0,
		0,
		0,
	)
	if err != nil {
		return err
	}
	// state->SetFromCounting(countingState)
	state_m_total_match := 0x2E7000
	cs_m_total_match := 0x4A400
	m_total_match_len := 0x4000
	m_total_nomatch_len := 0x400
	num_bytes_to_copy := m_total_match_len + m_total_nomatch_len
	copy(
		d.state[state_m_total_match:state_m_total_match+num_bytes_to_copy],
		countingState[cs_m_total_match:cs_m_total_match+num_bytes_to_copy],
	)
	statePtr := uintptr(unsafe.Pointer(&d.state[0]))
	csPtr := uintptr(unsafe.Pointer(&countingState[0]))
	for c := 0; c < 4096; c++ {
		_, err := callWrapper(d.handle, d.ffxivOnce.offset_match_set_from_histo_normalized, statePtr, csPtr)
		if err != nil {
			return err
		}
		statePtr += 706 // sizeof t_o0coder_match
		csPtr += 1124   // sizeof Histo_match::counts (281*4)
	}
	for c := 0; c < 256; c++ {
		_, err := callWrapper(d.handle, d.ffxivOnce.offset_nomatch_set_from_histo_normalized, statePtr, csPtr)
		if err != nil {
			return err
		}
		statePtr += 592 // sizeof t_o0coder_nomatch
		csPtr += 1024   // sizeof Histo_nomatch::counts (256*4)
	}
	countingState = nil
	return nil
}

func (d *decompressor) State() []byte {
	return append([]byte(nil), d.state...)
}

func (d *decompressor) SharedDict() []byte {
	return append([]byte(nil), d.sharedDict...)
}

func (d *decompressor) Decompress(input []byte, outputLength int64) ([]byte, error) {
	output := make([]byte, outputLength)

	res, _, err := syscall.SyscallN(
		uintptr(d.ffxivOnce.handle)+d.ffxivOnce.offsetOodleNetwork1UDP_Decode,
		uintptr(unsafe.Pointer(&d.state[0])),
		uintptr(unsafe.Pointer(&d.sharedDict[0])),
		uintptr(unsafe.Pointer(&input[0])),
		uintptr(len(input)),
		uintptr(unsafe.Pointer(&output[0])),
		uintptr(outputLength),
	)
	if err != 0 {
		return nil, err
	}
	if res == 0 {
		return nil, errors.New("unable to decompress")
	}
	return output, nil
}

func callWrapper(handle syscall.Handle, offset uintptr, args ...uintptr) (uintptr, error) {
	procAddr := uintptr(handle) + offset
	res, _, err := syscall.SyscallN(procAddr, args...)
	if err != 0 {
		return 0, err
	}
	return res, nil
}
