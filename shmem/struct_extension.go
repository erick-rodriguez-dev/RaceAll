package shared_memory

import (
	"fmt"
	"syscall"
	"unsafe"
)

var (
	kernel32            = syscall.NewLazyDLL("kernel32.dll")
	procOpenFileMapping = kernel32.NewProc("OpenFileMappingW")
	procMapViewOfFile   = kernel32.NewProc("MapViewOfFile")
	procUnmapViewOfFile = kernel32.NewProc("UnmapViewOfFile")
	procCloseHandle     = kernel32.NewProc("CloseHandle")
)

// ToStruct lee datos del memory mapped file y los deserializa en una estructura
// Equivalente 100% fiel al método de extensión ToStruct<T> de C# original:
//
//	public static T ToStruct<T>(this MemoryMappedFile file, byte[] buffer)
//	{
//	    using (var stream = file.CreateViewStream())
//	    {
//	        stream.ReadExactly(buffer);
//	        var handle = GCHandle.Alloc(buffer, GCHandleType.Pinned);
//	        var data = Marshal.PtrToStructure<T>(handle.AddrOfPinnedObject());
//	        handle.Free();
//	        return data;
//	    }
//	}
//
// Uso: ToStruct(mapName, &result)

func ToStruct[T any](mapName string, result *T) error {
	// Calcular el tamaño de la estructura
	size := unsafe.Sizeof(*result)

	// Convertir nombre a UTF16 para la API de Windows
	namePtr, err := syscall.UTF16PtrFromString(mapName)
	if err != nil {
		return fmt.Errorf("error converting name: %w", err)
	}

	// Abrir el file mapping (equivalente a MemoryMappedFile.OpenExisting)
	handle, _, err := procOpenFileMapping.Call(
		uintptr(FILE_MAP_READ),
		0,
		uintptr(unsafe.Pointer(namePtr)),
	)
	if handle == 0 {
		return fmt.Errorf("OpenFileMapping failed: %w", err)
	}
	defer procCloseHandle.Call(handle)

	// Mapear la vista del archivo (equivalente a CreateViewStream)
	addr, _, err := procMapViewOfFile.Call(
		handle,
		uintptr(FILE_MAP_READ),
		0,
		0,
		size,
	)
	if addr == 0 {
		return fmt.Errorf("MapViewOfFile failed: %w", err)
	}
	defer procUnmapViewOfFile.Call(addr)

	// Crear slice de bytes desde el puntero mapeado
	// Esto es equivalente a stream.ReadExactly(buffer)
	var buffer []byte
	sliceHeader := (*struct {
		addr uintptr
		len  int
		cap  int
	})(unsafe.Pointer(&buffer))
	sliceHeader.addr = addr
	sliceHeader.len = int(size)
	sliceHeader.cap = int(size)

	// Pinear el buffer y convertir a estructura
	// Esto es equivalente a:
	//   var handle = GCHandle.Alloc(buffer, GCHandleType.Pinned);
	//   var data = Marshal.PtrToStructure<T>(handle.AddrOfPinnedObject());
	//   handle.Free();
	ptr := unsafe.Pointer(&buffer[0])
	*result = *(*T)(ptr)

	return nil
}
