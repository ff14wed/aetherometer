// +build windows

package win32

import (
	"fmt"
	"syscall"
	"unsafe"

	"golang.org/x/sys/windows"
)

var (
	modkernel32 = windows.NewLazySystemDLL("kernel32.dll")

	procVirtualAllocEx = modkernel32.NewProc("VirtualAllocEx")
	procVirtualFreeEx  = modkernel32.NewProc("VirtualFreeEx")

	procWriteProcessMemory = modkernel32.NewProc("WriteProcessMemory")
	procLoadLibraryW       = modkernel32.NewProc("LoadLibraryW")
	// procFreeLibrary        = modkernel32.NewProc("FreeLibrary")
	procCreateRemoteThread = modkernel32.NewProc("CreateRemoteThread")
)

// nolint: golint
/*
LPVOID WINAPI VirtualAllocEx(
  _In_     HANDLE hProcess,
  _In_opt_ LPVOID lpAddress,
  _In_     SIZE_T dwSize,
  _In_     DWORD  flAllocationType,
  _In_     DWORD  flProtect
);
*/
func VirtualAllocEx(hProcess windows.Handle, lpAddress uintptr, dwSize int64, flAllocationType, flProtect uint32) (addr uintptr, err error) {
	addr, _, e1 := syscall.Syscall6(procVirtualAllocEx.Addr(), 5,
		uintptr(hProcess),
		lpAddress,
		uintptr(dwSize),
		uintptr(flAllocationType),
		uintptr(flProtect),
		0,
	)
	if addr == 0 {
		err = fmt.Errorf("VirtualAllocEx failed: code %d", e1)
	}
	return
}

// nolint: golint
/*
BOOL WINAPI VirtualFreeEx(
  _In_ HANDLE hProcess,
  _In_ LPVOID lpAddress,
  _In_ SIZE_T dwSize,
  _In_ DWORD  dwFreeType
);
*/
func VirtualFreeEx(hProcess windows.Handle, lpAddress uintptr, dwSize int64, dwFreeType uint32) (err error) {
	ret, _, e1 := syscall.Syscall6(procVirtualFreeEx.Addr(), 4,
		uintptr(hProcess),
		lpAddress,
		uintptr(dwSize),
		uintptr(dwFreeType),
		0,
		0,
	)
	if ret == 0 {
		err = fmt.Errorf("VirtualFreeEx failed: code %d", e1)
	}
	return

}

// nolint: golint
/*
BOOL WINAPI WriteProcessMemory(
  _In_  HANDLE  hProcess,
  _In_  LPVOID  lpBaseAddress,
  _In_  LPCVOID lpBuffer,
  _In_  SIZE_T  nSize,
  _Out_ SIZE_T  *lpNumberOfBytesWritten
);
*/
func WriteProcessMemory(hProcess windows.Handle, lpBaseAddress uintptr, buffer []byte, nSize int64) (err error) {
	/* #nosec */
	lpBuffer := uintptr(unsafe.Pointer(&buffer[0]))
	ret, _, e1 := syscall.Syscall6(procWriteProcessMemory.Addr(), 5,
		uintptr(hProcess),
		lpBaseAddress,
		lpBuffer,
		uintptr(nSize),
		0,
		0,
	)
	if ret == 0 {
		err = fmt.Errorf("WriteProcessMemory failed: code %d", e1)
	}
	return
}

// nolint: golint
/*
HANDLE WINAPI CreateRemoteThread(
  _In_  HANDLE                 hProcess,
  _In_  LPSECURITY_ATTRIBUTES  lpThreadAttributes,
  _In_  SIZE_T                 dwStackSize,
  _In_  LPTHREAD_START_ROUTINE lpStartAddress,
  _In_  LPVOID                 lpParameter,
  _In_  DWORD                  dwCreationFlags,
  _Out_ LPDWORD                lpThreadId
);
*/
func CreateRemoteThread(hProcess windows.Handle, lpThreadAttributes, dwStackSize, lpStartAddress, lpParameter, dwCreationFlags uintptr) (hThread windows.Handle, err error) {
	ret, _, e1 := syscall.Syscall9(procCreateRemoteThread.Addr(), 7,
		uintptr(hProcess),
		lpThreadAttributes,
		dwStackSize,
		lpStartAddress,
		lpParameter,
		dwCreationFlags,
		0,
		0,
		0,
	)
	hThread = windows.Handle(ret)
	if ret == 0 {
		err = fmt.Errorf("CreateRemoteThread failed: code %d", e1)
	}
	return
}
