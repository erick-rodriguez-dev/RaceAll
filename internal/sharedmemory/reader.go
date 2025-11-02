package sharedmemory

import (
	"syscall"
	"unsafe"

	"RaceAll/internal/errors"
)

const moduleName = "shared-memory"

func NewError(op string, err error) error {
	return errors.NewError(moduleName, op, err)
}

func NewErrorWithContext(op string, err error, ctx string) error {
	return errors.NewErrorWithContext(moduleName, op, err, ctx)
}

var (
	kernel32            = syscall.NewLazyDLL("kernel32.dll")
	procOpenFileMapping = kernel32.NewProc("OpenFileMappingW")
	procMapViewOfFile   = kernel32.NewProc("MapViewOfFile")
	procUnmapViewOfFile = kernel32.NewProc("UnmapViewOfFile")
	procCloseHandle     = kernel32.NewProc("CloseHandle")
)

const (
	FILE_MAP_READ = 0x0004
)

type SharedMemoryReader struct {
	physicsHandle  syscall.Handle
	graphicsHandle syscall.Handle
	staticHandle   syscall.Handle

	physicsAddr  uintptr
	graphicsAddr uintptr
	staticAddr   uintptr
}

func NewSharedMemoryReader() *SharedMemoryReader {
	return &SharedMemoryReader{}
}

func (r *SharedMemoryReader) Connect() error {
	var err error

	// Open Physics shared memory
	r.physicsHandle, err = openFileMapping(AccSharedMemoryName)
	if err != nil {
		return NewError("Connect", errors.ErrSharedMemoryNotFound)
	}

	r.physicsAddr, err = mapViewOfFile(r.physicsHandle, PhysicsPageFileSize)
	if err != nil {
		r.Disconnect()
		return NewError("Connect", errors.ErrSharedMemoryMap)
	}

	// Open Graphics shared memory
	r.graphicsHandle, err = openFileMapping(AccGraphicsMemoryName)
	if err != nil {
		r.Disconnect()
		return NewError("Connect", errors.ErrSharedMemoryNotFound)
	}

	r.graphicsAddr, err = mapViewOfFile(r.graphicsHandle, GraphicsPageFileSize)
	if err != nil {
		r.Disconnect()
		return NewError("Connect", errors.ErrSharedMemoryMap)
	}

	// Open Static shared memory
	r.staticHandle, err = openFileMapping(AccStaticMemoryName)
	if err != nil {
		r.Disconnect()
		return NewError("Connect", errors.ErrSharedMemoryNotFound)
	}

	r.staticAddr, err = mapViewOfFile(r.staticHandle, StaticPageFileSize)
	if err != nil {
		r.Disconnect()
		return NewError("Connect", errors.ErrSharedMemoryMap)
	}

	return nil
}

// Disconnect closes all shared memory handles
func (r *SharedMemoryReader) Disconnect() {
	if r.physicsAddr != 0 {
		unmapViewOfFile(r.physicsAddr)
		r.physicsAddr = 0
	}
	if r.physicsHandle != 0 {
		closeHandle(r.physicsHandle)
		r.physicsHandle = 0
	}

	if r.graphicsAddr != 0 {
		unmapViewOfFile(r.graphicsAddr)
		r.graphicsAddr = 0
	}
	if r.graphicsHandle != 0 {
		closeHandle(r.graphicsHandle)
		r.graphicsHandle = 0
	}

	if r.staticAddr != 0 {
		unmapViewOfFile(r.staticAddr)
		r.staticAddr = 0
	}
	if r.staticHandle != 0 {
		closeHandle(r.staticHandle)
		r.staticHandle = 0
	}
}

func (r *SharedMemoryReader) ReadPhysics() (*Physics, error) {
	if r.physicsAddr == 0 {
		return nil, NewError("ReadPhysics", errors.ErrNotConnected)
	}

	physics := (*Physics)(unsafe.Pointer(r.physicsAddr))
	return physics, nil
}

func (r *SharedMemoryReader) ReadGraphics() (*Graphics, error) {
	if r.graphicsAddr == 0 {
		return nil, NewError("ReadGraphics", errors.ErrNotConnected)
	}

	graphics := (*Graphics)(unsafe.Pointer(r.graphicsAddr))
	return graphics, nil
}

func (r *SharedMemoryReader) ReadStatic() (*Static, error) {
	if r.staticAddr == 0 {
		return nil, NewError("ReadStatic", errors.ErrNotConnected)
	}

	static := (*Static)(unsafe.Pointer(r.staticAddr))
	return static, nil
}

// IsConnected returns true if all shared memory handles are connected
func (r *SharedMemoryReader) IsConnected() bool {
	return r.physicsHandle != 0 && r.graphicsHandle != 0 && r.staticHandle != 0
}


func openFileMapping(name string) (syscall.Handle, error) {
	namePtr, err := syscall.UTF16PtrFromString(name)
	if err != nil {
		return 0, err
	}

	handle, _, err := procOpenFileMapping.Call(
		uintptr(FILE_MAP_READ),
		0,
		uintptr(unsafe.Pointer(namePtr)),
	)

	if handle == 0 {
		return 0, errors.ErrSharedMemoryAccess
	}

	return syscall.Handle(handle), nil
}

func mapViewOfFile(handle syscall.Handle, size int) (uintptr, error) {
	addr, _, err := procMapViewOfFile.Call(
		uintptr(handle),
		uintptr(FILE_MAP_READ),
		0,
		0,
		uintptr(size),
	)

	if addr == 0 {
		return 0, err
	}

	return addr, nil
}

func unmapViewOfFile(addr uintptr) {
	procUnmapViewOfFile.Call(addr)
}

func closeHandle(handle syscall.Handle) {
	procCloseHandle.Call(uintptr(handle))
}
