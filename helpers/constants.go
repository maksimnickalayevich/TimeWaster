package helpers

import "syscall"

const ACCESS = syscall.STANDARD_RIGHTS_READ | syscall.PROCESS_QUERY_INFORMATION | syscall.SYNCHRONIZE
const StillAlive = 259

type StorageType int

const (
	JsonFile StorageType = iota
	Database StorageType = iota
)
