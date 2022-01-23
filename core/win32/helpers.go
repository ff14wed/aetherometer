//go:build windows
// +build windows

package win32

import (
	"encoding/binary"
	"fmt"
	"io"
	"net"
	"path/filepath"
	"strings"
	"time"
	"unicode/utf16"
	"unsafe"

	winio "github.com/Microsoft/go-winio"
	"golang.org/x/sys/windows"
)

type Provider struct{}

var pe32Size = uint32(unsafe.Sizeof(windows.ProcessEntry32{})) // nolint: gosec
var me32Size = uint32(unsafe.Sizeof(ModuleEntry32{}))          // nolint: gosec

// EnumerateProcesses enumerates all processes running on the system and
// returns a map of ProcessID -> ProcessName
//
// API Methods used:
// CreateToolhelp32Snapshot https://docs.microsoft.com/en-us/windows/desktop/api/tlhelp32/nf-tlhelp32-createtoolhelp32snapshot
// Process32First https://docs.microsoft.com/en-us/windows/desktop/api/tlhelp32/nf-tlhelp32-process32firstw
// Process32Next https://docs.microsoft.com/en-us/windows/desktop/api/tlhelp32/nf-tlhelp32-process32nextw
// CloseHandle https://msdn.microsoft.com/en-us/library/windows/desktop/ms724211(v=vs.85).aspx
func (p Provider) EnumerateProcesses() (map[uint32]string, error) {
	snapshotHandle, err := windows.CreateToolhelp32Snapshot(windows.TH32CS_SNAPPROCESS, 0)
	if err != nil {
		return nil, fmt.Errorf("CreateToolhelp32Snapshot error: %s", err.Error())
	}
	defer func() {
		_ = windows.CloseHandle(snapshotHandle)
	}()
	var pe32 windows.ProcessEntry32
	pe32.Size = pe32Size
	err = windows.Process32First(snapshotHandle, &pe32)
	if err != nil {
		return nil, fmt.Errorf("Process32First error: %s", err.Error())
	}

	procMap := make(map[uint32]string)
	for err == nil {
		processName := string(utf16.Decode(pe32.ExeFile[:]))
		procMap[pe32.ProcessID] = processName
		err = windows.Process32Next(snapshotHandle, &pe32)
	}
	return procMap, nil
}

// EnumerateProcessModules enumerates all modules loaded inside a process and
// returns a list of modules
//
// API Methods used:
// CreateToolhelp32Snapshot https://docs.microsoft.com/en-us/windows/desktop/api/tlhelp32/nf-tlhelp32-createtoolhelp32snapshot
// Module32First https://docs.microsoft.com/en-us/windows/win32/api/tlhelp32/nf-tlhelp32-module32firstw
// Module32Next https://docs.microsoft.com/en-us/windows/desktop/api/tlhelp32/nf-tlhelp32-module32nextw
// CloseHandle https://msdn.microsoft.com/en-us/library/windows/desktop/ms724211(v=vs.85).aspx
func (p Provider) EnumerateProcessModules(pid uint32) ([]string, error) {
	snapshotHandle, err := windows.CreateToolhelp32Snapshot(windows.TH32CS_SNAPMODULE, pid)
	if err != nil {
		return nil, fmt.Errorf("CreateToolhelp32Snapshot error: %s", err.Error())
	}
	defer func() {
		_ = windows.CloseHandle(snapshotHandle)
	}()
	var me32 ModuleEntry32
	me32.Size = me32Size
	err = Module32First(snapshotHandle, &me32)
	if err != nil {
		return nil, fmt.Errorf("Module32First error: %s", err.Error())
	}

	var modules []string
	for err == nil {
		moduleName := string(utf16.Decode(me32.ModuleName[:]))
		moduleName = strings.Trim(moduleName, "\x00")
		modules = append(modules, moduleName)
		err = Module32Next(snapshotHandle, &me32)
	}
	return modules, nil
}

func getWString(s string) []byte {
	wChars := utf16.Encode(append([]rune(s), 0))
	sBytes := make([]byte, len(wChars)*2)
	for i := 0; i < len(wChars); i++ {
		binary.LittleEndian.PutUint16(sBytes[i*2:], wChars[i])
	}
	return sBytes
}

func getRunningTime(handle windows.Handle) (time.Duration, error) {
	var creationTime, exitTime, kernelTime, userTime windows.Filetime
	err := windows.GetProcessTimes(handle, &creationTime, &exitTime, &kernelTime, &userTime)
	if err != nil {
		return 0, err
	}
	createdAt := time.Unix(0, creationTime.Nanoseconds())
	return time.Since(createdAt), nil
}

type ErrDLLAlreadyInjected struct {
	DLLName string
}

func (e ErrDLLAlreadyInjected) Error() string {
	return fmt.Sprintf("DLL %s already injected.", e.DLLName)
}

func (ErrDLLAlreadyInjected) IsDLLAlreadyInjectedError() {}

// InjectDLL injects a library into another process on the system
//
// API Methods used:
// OpenProcess https://docs.microsoft.com/en-us/windows/desktop/api/processthreadsapi/nf-processthreadsapi-openprocess
// VirtualAllocEx https://msdn.microsoft.com/en-us/library/windows/desktop/aa366890(v=vs.85).aspx
// WriteProcessMemory https://msdn.microsoft.com/en-us/library/windows/desktop/ms681674(v=vs.85).aspx
// LoadLibraryW https://docs.microsoft.com/en-us/windows/desktop/api/libloaderapi/nf-libloaderapi-loadlibraryw
// CreateRemoteThread https://docs.microsoft.com/en-us/windows/desktop/api/processthreadsapi/nf-processthreadsapi-createremotethread
// WaitForSingleObject https://docs.microsoft.com/en-us/windows/desktop/api/synchapi/nf-synchapi-waitforsingleobject
// VirtualFreeEx https://msdn.microsoft.com/en-us/library/windows/desktop/aa366894(v=vs.85).aspx
// CloseHandle https://msdn.microsoft.com/en-us/library/windows/desktop/ms724211(v=vs.85).aspx
func (p Provider) InjectDLL(processID uint32, payloadPath string) error {
	pathBytes := getWString(payloadPath)
	lenPathBytes := int64(len(pathBytes))

	sd, err := windows.GetSecurityInfo(windows.CurrentProcess(), windows.SE_KERNEL_OBJECT, windows.DACL_SECURITY_INFORMATION)
	if err != nil {
		return fmt.Errorf("GetSecurityInfo: %s", err)
	}

	baseHandle, err := windows.OpenProcess(
		windows.PROCESS_QUERY_LIMITED_INFORMATION|
			windows.WRITE_DAC|
			windows.READ_CONTROL,
		false, processID)
	if err != nil {
		return fmt.Errorf("open process pid %d: %s", processID, err)
	}

	runningTime, err := getRunningTime(baseHandle)
	if err != nil {
		return fmt.Errorf("get runningTime time for %d: %s", processID, err)
	}

	if runningTime < 5*time.Second {
		time.Sleep(5*time.Second - runningTime)
	}

	dacl, _, err := sd.DACL()
	if err != nil {
		return fmt.Errorf("getting DACL for current process: %s", err)
	}
	windows.SetSecurityInfo(baseHandle, windows.SE_KERNEL_OBJECT, windows.DACL_SECURITY_INFORMATION|windows.UNPROTECTED_DACL_SECURITY_INFORMATION, nil, nil, dacl, nil)

	_ = windows.CloseHandle(baseHandle)

	modules, err := p.EnumerateProcessModules(processID)
	if err != nil {
		return fmt.Errorf("EnumerateProcessModules: %s", err)
	}
	payloadName := filepath.Base(payloadPath)
	var alreadyExists bool
	for _, m := range modules {
		if m == payloadName {
			alreadyExists = true
		}
	}
	if alreadyExists {
		return ErrDLLAlreadyInjected{
			DLLName: payloadName,
		}
	}

	hProcess, err := windows.OpenProcess(
		windows.PROCESS_QUERY_INFORMATION|
			windows.PROCESS_CREATE_THREAD|
			windows.PROCESS_VM_OPERATION|
			windows.PROCESS_VM_WRITE,
		false,
		processID,
	)
	if err != nil {
		return fmt.Errorf("open process pid %d: %s", processID, err)
	}

	remotePathAddr, err := VirtualAllocEx(hProcess, 0, lenPathBytes, MEM_COMMIT, PAGE_READWRITE)
	if err != nil {
		return fmt.Errorf("reserving memory in target pid %d: %s", processID, err)
	}
	err = WriteProcessMemory(hProcess, remotePathAddr, pathBytes, lenPathBytes)
	if err != nil {
		return fmt.Errorf("writing data into target pid %d memory: %s", processID, err)
	}
	// Get the address of LoadLibraryW in our own loaded kernel32... but it should
	// be the same as any other process
	loadLibraryAddr := procLoadLibraryW.Addr()
	hThread, err := CreateRemoteThread(hProcess, 0, 0, loadLibraryAddr, remotePathAddr, 0)
	if err != nil {
		return fmt.Errorf("running LoadLibrary in target pid %d: %s", processID, err)
	}

	// Wait for remote thread to terminate
	_, err = windows.WaitForSingleObject(hThread, windows.INFINITE)
	if err != nil {
		return fmt.Errorf("waiting for remote thread in target pid %d: %s", processID, err)
	}

	if remotePathAddr != 0 {
		err = VirtualFreeEx(hProcess, remotePathAddr, 0, MEM_RELEASE)
		if err != nil {
			return fmt.Errorf("freeing reserved memory in target pid %d: %s", processID, err)
		}
	}
	if hThread != 0 {
		_ = windows.CloseHandle(hThread)
	}
	if hProcess != 0 {
		_ = windows.CloseHandle(hProcess)
	}

	return nil
}

// DialPipe dials a named pipe on the system for interprocess communication
func (p Provider) DialPipe(path string, timeout *time.Duration) (net.Conn, error) {
	return winio.DialPipe(path, timeout)
}

// IsPipeClosed returns whether or the error is due to the named pipe connection
// being closed or if it's another error
func (p Provider) IsPipeClosed(err error) bool {
	return err == io.EOF || err == winio.ErrFileClosed
}
