// Package cgoref a way to tell Go to not move a particular variable around
// in memory. It is not advised to generally use this package, as it may
// memory fragmentation, high memory usage, and other funky things.
package cgoref

import (
	"fmt"
	"unsafe"
)

/*
 #include <pthread.h>

 void *hold(pthread_mutex_t *mut, void *ptr)
 {
     pthread_mutex_lock(mut);
     pthread_mutex_lock(mut);
     return ptr;
 }

*/
import "C"

type CRef struct {
	mut C.pthread_mutex_t
}

// Ref will hold a given pointer inside a CGO call until UnRef is called.
// In Go 1.6, this should prevent Go from moving it around in memory
// and changing its address until it has been unreffed. This will, of course,
// completey hose Go's garbage collector's ability to be effective. You
// may end up with OOMs or very fragmented memory, or other wacky things.
//
// Mostly untested.
func Ref(ptr unsafe.Pointer) (*CRef, error) {
	ret := new(CRef)

	// The mutex should also not move during this time,
	// allowing UnRef to work.
	cret := C.pthread_mutex_init(&ret.mut, nil)
	if int(cret) != 0 {
		return nil, fmt.Errorf("Failed to init mutex!")
	}

	// Function will block until UnRef is called, keeping
	// the pointer inside a CGO call.
	go C.hold(&ret.mut, ptr)

	return ret, nil
}

// UnRef will allow the pointer to leave the CGO call, letting
// the Go garbage collector clean/move it.
func (this *CRef) UnRef() {
	C.pthread_mutex_unlock(&this.mut)
	C.pthread_mutex_destroy(&this.mut)
}
