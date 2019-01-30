// +build windows

package win32

// nolint
const (
	PROCESS_CREATE_THREAD = 0x0002
	PROCESS_VM_OPERATION  = 0x0008
	PROCESS_VM_WRITE      = 0x0020

	MEM_COMMIT     = 0x00001000
	PAGE_READWRITE = 0x04

	MEM_RELEASE = 0x8000
)
