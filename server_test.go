package main

import (
	. "github.com/smartystreets/goconvey/convey"
	"net"
	"testing"
	"reflect"
)

func TestCreateServer(t *testing.T) {
	Convey("Server should be created without errors", t, func() {
		s := &Server{
			priority: "1",
			joins:    make(chan net.Conn),
		}
		So(s.priority, ShouldEqual, "1")
	})
}

func TestParseRequest(t *testing.T) {
	Convey("Test parse request message", t, func() {
		s := &Server{}
		msg := "[3]ALIVE?"
		expectedReq := Message{
			msg: "ALIVE?",
			priority: "3",
		}
		req, _ := s.parseProtocol(msg)
		So(reflect.DeepEqual(req, expectedReq), ShouldBeTrue)
	})
}