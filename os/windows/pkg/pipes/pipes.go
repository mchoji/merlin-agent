// +build windows

// Merlin is a post-exploitation command and control framework.
// This file is part of Merlin.
// Copyright (C) 2022  Russel Van Tuyl

// Merlin is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// any later version.

// Merlin is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU General Public License for more details.

// You should have received a copy of the GNU General Public License
// along with Merlin.  If not, see <http://www.gnu.org/licenses/>.

package pipes

import (
	"fmt"
	"golang.org/x/sys/windows"
)

// CreateAnonymousPipes creates and returns a handle for STDIN, STDOUT, and STDERR
func CreateAnonymousPipes() (stdInRead, stdInWrite, stdOutRead, stdOutWrite, stdErrRead, stdErrWrite windows.Handle, err error) {
	// Create anonymous pipe for STDIN
	err = windows.CreatePipe(&stdInRead, &stdInWrite, &windows.SecurityAttributes{InheritHandle: 1}, 0)
	if err != nil {
		err = fmt.Errorf("error creating the STDIN pipe:\r\n%s", err)
		return
	}

	// Create anonymous pipe for STDOUT
	err = windows.CreatePipe(&stdOutRead, &stdOutWrite, &windows.SecurityAttributes{InheritHandle: 1}, 0)
	if err != nil {
		err = fmt.Errorf("error creating the STDOUT pipe:\r\n%s", err)
		return
	}

	// Create anonymous pipe for STDERR
	err = windows.CreatePipe(&stdErrRead, &stdErrWrite, &windows.SecurityAttributes{InheritHandle: 1}, 0)
	if err != nil {
		err = fmt.Errorf("error creating the STDERR pipe:\r\n%s", err)
		return
	}

	err = nil
	return
}

// ClosePipes closes the handle for all the passed in STDIN, STDOUT, and STDERR read and write handles
func ClosePipes(stdInRead, stdInWrite, stdOutRead, stdOutWrite, stdErrRead, stdErrWrite windows.Handle) (err error) {
	// STDIN - Read
	if stdInRead != 0 {
		err = windows.CloseHandle(stdInRead)
		if err != nil {
			err = fmt.Errorf("error closing the STDIN read pipe handle: %s", err)
			return
		}
	}

	// STDIN - Write
	if stdInWrite != 0 {
		err = windows.CloseHandle(stdInWrite)
		if err != nil {
			err = fmt.Errorf("error closing the STDIN write pipe handle: %s", err)
			return
		}
	}

	// STDOUT - Read
	if stdOutRead != 0 {
		err = windows.CloseHandle(stdOutRead)
		if err != nil {
			err = fmt.Errorf("error closing the STDOUT read pipe handle: %s", err)
			return
		}
	}

	// STDOUT - Write
	if stdOutWrite != 0 {
		err = windows.CloseHandle(stdOutWrite)
		if err != nil {
			err = fmt.Errorf("error closing the STDOUT write pipe handle: %s", err)
			return
		}
	}

	// STDERR - Read
	if stdErrRead != 0 {
		err = windows.CloseHandle(stdErrRead)
		if err != nil {
			err = fmt.Errorf("error closing the STDERR read pipe handle: %s", err)
			return
		}
	}

	// STDERR - Write
	if stdErrWrite != 0 {
		err = windows.CloseHandle(stdErrWrite)
		if err != nil {
			err = fmt.Errorf("error closing the STDERR write pipe handle: %s", err)
			return
		}
	}

	err = nil
	return
}

// ReadPipes reads data from the passed in STDIN, STDOUT, and STDERR pipes and returns it as a string
func ReadPipes(stdInRead, stdOutRead, stdErrRead windows.Handle) (stdin, stdout, stderr string, err error) {
	// Read STDOUT from child process
	/*
		BOOL ReadFile(
		HANDLE       hFile,
		LPVOID       lpBuffer,
		DWORD        nNumberOfBytesToRead,
		LPDWORD      lpNumberOfBytesRead,
		LPOVERLAPPED lpOverlapped
		);
	*/
	nNumberOfBytesToRead := make([]byte, 1)

	// STDIN
	if stdInRead != 0 {
		// Read STDIN
		var stdInBuffer []byte
		var stdInDone uint32
		var stdInOverlapped windows.Overlapped

		for {
			errReadFileStdErr := windows.ReadFile(stdInRead, nNumberOfBytesToRead, &stdInDone, &stdInOverlapped)
			if errReadFileStdErr != nil && errReadFileStdErr.Error() != "The pipe has been ended." {
				stderr = fmt.Sprintf("error reading from STDIN pipe: %s", errReadFileStdErr)
				return
			}
			if int(stdInDone) == 0 {
				break
			}
			for _, b := range nNumberOfBytesToRead {
				stdInBuffer = append(stdInBuffer, b)
			}
		}
		stdin = string(stdInBuffer)
	}

	// STDOUT
	if stdOutRead != 0 {
		var stdOutBuffer []byte
		var stdOutDone uint32
		var stdOutOverlapped windows.Overlapped

		// ReadFile on STDOUT pipe
		for {
			errReadFileStdOut := windows.ReadFile(stdOutRead, nNumberOfBytesToRead, &stdOutDone, &stdOutOverlapped)
			if errReadFileStdOut != nil && errReadFileStdOut.Error() != "The pipe has been ended." {
				stderr = fmt.Sprintf("error reading from STDOUT pipe: %s", errReadFileStdOut)
				return
			}
			if int(stdOutDone) == 0 {
				break
			}
			for _, b := range nNumberOfBytesToRead {
				stdOutBuffer = append(stdOutBuffer, b)
			}
		}
		stdout = string(stdOutBuffer)
	}

	// STDERR
	if stdErrRead != 0 {
		// Read STDERR
		var stdErrBuffer []byte
		var stdErrDone uint32
		var stdErrOverlapped windows.Overlapped

		for {
			errReadFileStdErr := windows.ReadFile(stdErrRead, nNumberOfBytesToRead, &stdErrDone, &stdErrOverlapped)
			if errReadFileStdErr != nil && errReadFileStdErr.Error() != "The pipe has been ended." {
				stderr = fmt.Sprintf("error reading from STDOUT pipe: %s", errReadFileStdErr)
				return
			}
			if int(stdErrDone) == 0 {
				break
			}
			for _, b := range nNumberOfBytesToRead {
				stdErrBuffer = append(stdErrBuffer, b)
			}
		}
		stderr = string(stdErrBuffer)
	}

	err = nil
	return
}
