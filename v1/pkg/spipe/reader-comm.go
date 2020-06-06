package spipe

import "io"

type reader interface {
	io.Reader
	hasNext() bool
	nextInput() io.Reader
	popInput() error
}

func internalRead(r reader, p []byte) (totalRead int, err error) {
	// If we have no more available readers, return an EOF.
	if !r.hasNext() {
		return 0, io.EOF
	}

	ln := len(p)
	pos := 0

	// Read the current input until it EOFs or throws some other error.
	for totalRead < ln {
		n, e := r.nextInput().Read(p[pos:])
		totalRead += n

		// If the last read returned an error
		if e != nil {
			// And that error was an EOF, the stream is dead, skip out of the loop and
			// continue.
			if e == io.EOF {
				break
			}

			// And that error was not an EOF, return it and halt.
			err = e
			return
		}

		// Move the current position up by the number of bytes read
		pos += n
	}

	// If the loop filled the buffer, then we have nothing more to do.
	if totalRead >= ln {
		return
	}

	// if the last read resulted in fewer bytes read than len(p), pop the dead
	// reader out of the queue and try filling the remainder with the next reader
	// (if any exist).
	if err = r.popInput(); err != nil {
		return
	}

	// Try a read using the unwritten part of the input buffer.
	n, err := r.Read(p[pos:])

	// If we got an EOF from our last read, then we have no input readers left
	// to use to fill the input buffer.  If we have also read more than 0 bytes
	// overall, clear the error for this return, they will get it on the next
	// Read call (if one is made).
	if n > 0 && err == io.EOF {
		err = nil
	}

	// Append the number of additional bytes read from the recursive call to the
	// overall read count.
	totalRead += n

	return
}
