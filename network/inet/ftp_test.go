/*
 * Copyright (C) 2026 Franklin D. Amador
 *
 * This software is dual-licensed under:
 * - EPL v2.0
 * - Commercial
 *
 * You may choose to use this software under the terms of either license.
 * See the LICENSE files in the project root for full license text.
 */

package inet

import (
	"bufio"
	"fmt"
	"io"
	"net"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/apiarytech/royaljelly/iec"
)

// ftpServer is an FTP server of the files in memory, in the passive mode,
// on 127.0.0.1:21.
type ftpServer struct {
	mu    sync.Mutex
	files map[string][]byte
	log   []string
}

func startFTP(t *testing.T) *ftpServer {
	ln, err := net.Listen("tcp4", "127.0.0.1:21")
	if err != nil {
		t.Skipf("port 21: %v", err)
	}
	t.Cleanup(func() { ln.Close() })
	s := &ftpServer{files: map[string][]byte{}}
	go func() {
		for {
			c, err := ln.Accept()
			if err != nil {
				return
			}
			go s.serve(c)
		}
	}()
	return s
}

func (s *ftpServer) serve(c net.Conn) {
	defer c.Close()
	r := bufio.NewReader(c)
	reply := func(f string, a ...any) { fmt.Fprintf(c, f+"\r\n", a...) }
	reply("220 ready")
	var data net.Listener
	for {
		line, err := r.ReadString('\n')
		if err != nil {
			return
		}
		cmd, arg, _ := strings.Cut(strings.TrimSpace(line), " ")
		s.mu.Lock()
		s.log = append(s.log, cmd)
		s.mu.Unlock()
		switch cmd {
		case "USER":
			reply("331 password")
		case "PASS":
			reply("230 logged in")
		case "TYPE":
			reply("200 ok")
		case "SIZE":
			s.mu.Lock()
			reply("213 %d", len(s.files[arg]))
			s.mu.Unlock()
		case "PASV":
			data, _ = net.Listen("tcp4", "127.0.0.1:0")
			p := data.Addr().(*net.TCPAddr).Port
			reply("227 Entering Passive Mode (127,0,0,1,%d,%d)", p>>8, p&0xFF)
		case "STOR":
			reply("150 ok")
			dc, err := data.Accept()
			if err == nil {
				b, _ := io.ReadAll(dc)
				dc.Close()
				s.mu.Lock()
				s.files[arg] = b
				s.mu.Unlock()
			}
			data.Close()
			reply("226 done")
		case "RETR":
			reply("150 ok")
			dc, err := data.Accept()
			if err == nil {
				s.mu.Lock()
				dc.Write(s.files[arg])
				s.mu.Unlock()
				dc.Close()
			}
			data.Close()
			reply("226 done")
		case "DELE":
			s.mu.Lock()
			delete(s.files, arg)
			s.mu.Unlock()
			reply("250 deleted")
		case "QUIT":
			reply("221 bye")
			return
		default:
			reply("502 no")
		}
	}
}

func TestFTP(t *testing.T) {
	srv := startFTP(t)
	dir := t.TempDir()
	local := filepath.Join(dir, "data.csv")
	content := strings.Repeat("1;2;3\r\n", 2000) // more than a buffer
	if err := os.WriteFile(local, []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}

	var f FTP_CLIENT
	f.INIT()
	f.FILENAME = iec.STRING(local)
	f.FTP_URL = "ftp://anna:secret@127.0.0.1/up/"
	f.TIMEOUT = iec.TIME(5 * time.Second)
	f.ACTIVATE = true
	run := func() {
		t.Helper()
		end := time.Now().Add(15 * time.Second)
		for f.step == 0 && time.Now().Before(end) {
			f.Execute(time.Now())
		}
		for bool(f.BUSY) && time.Now().Before(end) {
			f.Execute(time.Now())
			time.Sleep(time.Millisecond)
		}
		if f.BUSY {
			t.Fatalf("still busy at step %d, reply %d", f.step, f.c1.rcvState)
		}
	}
	run()
	if !f.DONE || f.ERROR_T != 0 {
		t.Fatalf("upload: DONE %v, ERROR %08X/%d", f.DONE, f.ERROR_C, f.ERROR_T)
	}
	srv.mu.Lock()
	got := string(srv.files["/up/data.csv"])
	srv.mu.Unlock()
	if got != content {
		t.Fatalf("the server has %d bytes, want %d", len(got), len(content))
	}

	// Get it back, and delete it on the server.
	f.ACTIVATE = false
	f.Execute(time.Now())
	back := filepath.Join(dir, "back.csv")
	f.FILENAME, f.FTP_URL, f.FTP_DOWNLOAD, f.FILE_DELETE = iec.STRING(back), "ftp://127.0.0.1/up/data.csv", true, true
	f.ACTIVATE = true
	run()
	if !f.DONE || f.ERROR_T != 0 {
		t.Fatalf("download: DONE %v, ERROR %08X/%d", f.DONE, f.ERROR_C, f.ERROR_T)
	}
	if b, _ := os.ReadFile(back); string(b) != content {
		t.Errorf("the file read back has %d bytes, want %d", len(b), len(content))
	}
	srv.mu.Lock()
	_, there := srv.files["/up/data.csv"]
	srv.mu.Unlock()
	if there {
		t.Error("the file was not deleted on the server")
	}
}
