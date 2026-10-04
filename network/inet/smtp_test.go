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
	"crypto/hmac"
	"crypto/md5"
	gobase64 "encoding/base64"
	"encoding/hex"
	"fmt"
	"net"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	td "github.com/apiarytech/beebread/basic/time_date"
	"github.com/apiarytech/royaljelly/iec"
)

// smtpServer takes one mail, with the authentication auth, and sends what
// it got on the channel.
type smtpServer struct {
	ln   net.Listener
	auth string
	got  chan smtpMail
}

type smtpMail struct {
	user, pass, from string
	to               []string
	data             string
}

func startSMTP(t *testing.T, auth string) *smtpServer {
	ln, err := net.Listen("tcp4", "127.0.0.1:0")
	if err != nil {
		t.Skip(err)
	}
	t.Cleanup(func() { ln.Close() })
	s := &smtpServer{ln: ln, auth: auth, got: make(chan smtpMail, 1)}
	go func() {
		c, err := ln.Accept()
		if err != nil {
			return
		}
		defer c.Close()
		s.serve(c)
	}()
	return s
}

func (s *smtpServer) serve(c net.Conn) {
	r := bufio.NewReader(c)
	line := func() string { l, _ := r.ReadString('\n'); return strings.TrimRight(l, "\r\n") }
	reply := func(f string, a ...any) { fmt.Fprintf(c, f+"\r\n", a...) }
	b64 := func(s string) string { b, _ := gobase64.StdEncoding.DecodeString(s); return string(b) }
	var m smtpMail
	reply("220 test ESMTP")
	for {
		l := line()
		cmd, arg, _ := strings.Cut(l, " ")
		switch strings.ToUpper(cmd) {
		case "EHLO":
			reply("250-test")
			reply("250 AUTH %s", s.auth)
		case "HELO":
			reply("250 test")
		case "AUTH":
			switch {
			case strings.HasPrefix(arg, "PLAIN"):
				reply("334 ")
				parts := strings.Split(b64(line()), "\x00")
				m.user, m.pass = parts[1], parts[2]
			case strings.HasPrefix(arg, "LOGIN"):
				reply("334 VXNlcm5hbWU6")
				m.user = b64(line())
				reply("334 UGFzc3dvcmQ6")
				m.pass = b64(line())
			case strings.HasPrefix(arg, "CRAM-MD5"):
				challenge := "<12345.67890@test>"
				reply("334 %s", gobase64.StdEncoding.EncodeToString([]byte(challenge)))
				user, digest, _ := strings.Cut(b64(line()), " ")
				mac := hmac.New(md5.New, []byte("secret"))
				mac.Write([]byte(challenge))
				m.user = user
				if digest == hex.EncodeToString(mac.Sum(nil)) {
					m.pass = "secret"
				}
			}
			reply("235 ok")
		case "MAIL":
			m.from = arg
			reply("250 ok")
		case "RCPT":
			m.to = append(m.to, arg)
			reply("250 ok")
		case "DATA":
			reply("354 go on")
			var data strings.Builder
			for {
				l := line()
				if l == "." {
					break
				}
				data.WriteString(l + "\n")
			}
			m.data = data.String()
			reply("250 queued")
		case "QUIT":
			reply("221 bye")
			s.got <- m
			return
		default:
			reply("502 no")
		}
	}
}

func sendMail(t *testing.T, srv *smtpServer, files iec.STRING) smtpMail {
	t.Helper()
	port := srv.ln.Addr().(*net.TCPAddr).Port
	server := iec.STRING(fmt.Sprintf("smtp://anna:secret@127.0.0.1:%d", port))
	from, to := iec.STRING("plc@example.org;PLC 1"), iec.STRING("a@example.org;b@example.org#c@example.org#d@example.org")
	subject, body := iec.STRING("hello"), iec.STRING("the body")
	s := SMTP_CLIENT{SERVER: &server, MAILFROM: &from, MAILTO: &to, SUBJECT: &subject, BODY: &body, FILES: &files}
	s.INIT()
	s.DTI = td.SET_DT(2024, 5, 9, 12, 30, 0)
	s.ACTIVATE = true
	end := time.Now().Add(20 * time.Second)
	s.Execute(time.Now())
	for bool(s.BUSY) && time.Now().Before(end) {
		s.Execute(time.Now())
		time.Sleep(time.Millisecond)
	}
	if !s.DONE {
		t.Fatalf("DONE %v, ERROR %08X/%d, step %d", s.DONE, s.ERROR_C, s.ERROR_T, s.step)
	}
	select {
	case m := <-srv.got:
		return m
	case <-time.After(5 * time.Second):
		t.Fatal("the server got no mail")
	}
	return smtpMail{}
}

func TestSMTP(t *testing.T) {
	dir := t.TempDir()
	att := filepath.Join(dir, "report.txt")
	if err := os.WriteFile(att, []byte("attached data"), 0o644); err != nil {
		t.Fatal(err)
	}
	for _, auth := range []string{"PLAIN", "LOGIN", "CRAM-MD5"} {
		t.Run(auth, func(t *testing.T) {
			m := sendMail(t, startSMTP(t, auth), iec.STRING(att))
			if m.user != "anna" || m.pass != "secret" {
				t.Errorf("authentication: %q %q", m.user, m.pass)
			}
			if m.from != "FROM: <plc@example.org>" {
				t.Errorf("MAIL %q", m.from)
			}
			if strings.Join(m.to, ",") != "TO: <a@example.org>,TO: <b@example.org>,TO: <c@example.org>,TO: <d@example.org>" {
				t.Errorf("RCPT %q", m.to)
			}
			for _, want := range []string{
				`From: "PLC 1" <plc@example.org>`,
				"To: <a@example.org>\nTo: <b@example.org>\nCc: <c@example.org>\n",
				"Subject: hello\n",
				"Date: Thu, 09 May 2024 12:30:00 +0000\n",
				"the body\n",
				`filename="report.txt"`,
				gobase64.StdEncoding.EncodeToString([]byte("attached data")),
			} {
				if !strings.Contains(m.data, want) {
					t.Errorf("no %q in\n%s", want, m.data)
				}
			}
			if strings.Contains(m.data, "d@example.org") {
				t.Error("the blind copy is in the header")
			}
		})
	}
}
