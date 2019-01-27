package requests

import (
	"geo-observers-blockchain/core/common/types/transactions"
	"geo-observers-blockchain/core/network/communicator/geo/api/v0/common"
	"geo-observers-blockchain/core/network/communicator/geo/api/v0/requests"
	"geo-observers-blockchain/core/network/communicator/geo/api/v0/responses"
	testsCommon "geo-observers-blockchain/tests/network/geo"
	"testing"
)

const (
	TSLGetRequestID = 68
)

func TestTSLGetRequestID(t *testing.T) {
	if //noinspection GoBoolExpressions
	TSLGetRequestID != common.ReqTSLGet {
		t.Fatal()
	}
}

func TestTSLGet(t *testing.T) {
	{
		// TSL that is ABSENT in chain.
		response := requestTSLGet(t, transactions.NewTxID())
		if response.IsPresent {
			t.Error()
		}

		if response.TSL != nil {
			t.Error()
		}
	}

	{
		// TSL that is present in chain
		// todo: implement
		// todo: add observers clusters
	}
}

func requestTSLGet(t *testing.T, TxID *transactions.TxID) *responses.TSLGet {
	conn := testsCommon.ConnectToObserver(t)
	defer conn.Close()

	request := requests.NewTSLGet(TxID)
	testsCommon.SendRequest(t, request, conn)

	response := &responses.TSLGet{}
	testsCommon.GetResponse(t, response, conn)
	return response
}