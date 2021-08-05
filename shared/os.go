//go:build linux && cgo
// +build linux,cgo

package shared

import (
	"fmt"
	"unsafe"

	"golang.org/x/sys/unix"
)

/*
#ifndef _GNU_SOURCE
#define _GNU_SOURCE 1
#endif
#include <errno.h>
#include <fcntl.h>
#include <grp.h>
#include <limits.h>
#include <poll.h>
#include <pty.h>
#include <pwd.h>
#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <unistd.h>
#include <sys/stat.h>
#include <sys/socket.h>
#include <sys/types.h>
#include <sys/un.h>
*/
import "C"

// GroupID is an adaption from https://codereview.appspot.com/4589049.
func GroupID(name string) (int, error) {
	var grp C.struct_group
	var result *C.struct_group

	bufSize := C.sysconf(C._SC_GETGR_R_SIZE_MAX)
	if bufSize < 0 {
		bufSize = 4096
	}

	buf := C.malloc(C.size_t(bufSize))
	if buf == nil {
		return -1, fmt.Errorf("allocation failed")
	}

	cname := C.CString(name)
	defer C.free(unsafe.Pointer(cname))

again:
	rv, errno := C.getgrnam_r(cname,
		&grp,
		(*C.char)(buf),
		C.size_t(bufSize),
		&result)
	if rv != 0 {
		// OOM killer will take care of us if we end up doing this too
		// often.
		if errno == unix.ERANGE {
			bufSize *= 2
			tmp := C.realloc(buf, C.size_t(bufSize))
			if tmp == nil {
				return -1, fmt.Errorf("allocation failed")
			}
			buf = tmp
			goto again
		}

		C.free(buf)
		return -1, fmt.Errorf("failed group lookup: %s", unix.Errno(rv))
	}
	C.free(buf)

	if result == nil {
		return -1, fmt.Errorf("unknown group %s", name)
	}

	return int(C.int(result.gr_gid)), nil
}
