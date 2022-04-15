package main

import "syscall"

const ACCESS = syscall.STANDARD_RIGHTS_READ | syscall.PROCESS_QUERY_INFORMATION | syscall.SYNCHRONIZE
const STILL_ALIVE = 259
