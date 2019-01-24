package geo

import (
	"bufio"
	"encoding"
	"fmt"
	"geo-observers-blockchain/core/common"
	"geo-observers-blockchain/core/utils"
	"geo-observers-blockchain/tests"
	"net"
	"testing"
	"time"
)

func connectToObserver(t *testing.T) (conn net.Conn) {
	conn, err := net.Dial("tcp", fmt.Sprint(tests.ObserverHost, ":", tests.ObserverPort))
	if err != nil {
		t.Fatal("could not connect to observer: ", err)
	}

	return
}

func getResponse(t *testing.T, response encoding.BinaryUnmarshaler, conn net.Conn) {
	_ = conn.SetReadDeadline(time.Now().Add(time.Second * 3))
	reader := bufio.NewReader(conn)

	messageSizeBinary := []byte{0, 0, 0, 0}
	bytesRead, err := reader.Read(messageSizeBinary)
	if err != nil {
		t.Fatal(err)
	}

	if bytesRead != common.Uint32ByteSize {
		t.Fatal("Invalid message header received")
	}

	messageSize, err := utils.UnmarshalUint32(messageSizeBinary)
	if err != nil {
		t.Fatal(err)
	}

	_, _ = reader.Discard(4)
	var offset uint32 = 0
	data := make([]byte, messageSize, messageSize)
	for {
		bytesReceived, err := reader.Read(data[offset:])
		if err != nil {
			t.Fatal(err)
		}

		offset += uint32(bytesReceived)
		if offset == messageSize {
			err := response.UnmarshalBinary(data)
			if err != nil {
				t.Error()
				return
			}

			return
		}
	}
}

func sendRequest(t *testing.T, request encoding.BinaryMarshaler, conn net.Conn) {
	requestBinary, err := request.MarshalBinary()
	if err != nil {
		t.Error()
	}

	sendData(t, conn, requestBinary)
}

func sendData(t *testing.T, conn net.Conn, data []byte) {
	data = append([]byte{0}, data...) // protocol header

	var (
		dataLength       = uint64(len(data))
		dataLengthBinary = utils.MarshalUint64(dataLength)
	)

	_, err := conn.Write(dataLengthBinary)
	if err != nil {
		t.Error("cant send payload: ", err)
	}

	_, err = conn.Write(data)
	if err != nil {
		t.Error("cant send payload: ", err)
	}
}