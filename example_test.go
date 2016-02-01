package cgoref_test

func ExampleRef() {
	buf := make([]byte, 1024)

	var CThing C.struct_some_thing

	// Make sure it won't move during processing.
	ref, err := cgoref.Ref(unsafe.Pointer(&buf[0]))
	if err != nil {
		return nil, err
	}

	// Stash it in a C struct for multiple uses. Evil!
	CThing.buf = unsafe.Pointer(&buf[0])

	length := int(C.some_init(&CThing))

	for i := 0; i < length; i++ {
		C.some_process(&CThing)
	}

	// Hand Go back control.
	ref.UnRef()

	return buf, nil
}
