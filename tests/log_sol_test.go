//  Copyright 2020 The Example Authors
//  This file is part of the evm library.
//
//  The evm library is free software: you can redistribute it and/or modify
//  it under the terms of the GNU Lesser General Public License as published by
//  the Free Software Foundation, either version 3 of the License, or
//  (at your option) any later version.
//
//  The evm library is distributed in the hope that it will be useful,/
//  but WITHOUT ANY WARRANTY; without even the implied warranty of
//  MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
//  GNU Lesser General Public License for more details.
//
//  You should have received a copy of the GNU Lesser General Public License
//  along with the evm library. If not, see <http://www.gnu.org/licenses/>.
//

package tests

import (
	"github.com/sunvim/evm"
	"github.com/sunvim/evm/db"
	"github.com/sunvim/evm/example"
	"github.com/sunvim/evm/util"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

var (
	logAbi     = "sols/Log_sol_Log.abi"
	logBin     = "sols/Log_sol_Log.bin"
	logCode    []byte
	logAddress evm.Address
)

func TestLogSol(t *testing.T) {
	var err error
	binBytes, err := util.ReadBinFile(logBin)
	require.NoError(t, err)
	bc := example.NewBlockchain()
	memoryDB := db.NewMemory(bc.NewAccount)
	var origin = example.HexToAddress("6ac7ea33f8831ea9dcc53393aaa88b25a785dbf0")
	var exceptCode = `60806040523480156100115760006000fd5b50600436106100305760003560e01c806336b899bb1461003657610030565b60006000fd5b6101016004803603604081101561004d5760006000fd5b810190808035906020019064010000000081111561006b5760006000fd5b82018360208201111561007e5760006000fd5b803590602001918460018302840111640100000000831117156100a15760006000fd5b91908080601f016020809104026020016040519081016040528093929190818152602001838380828437600081840152601f19601f8201169050808301925050505050505090909192909091929080359060200190929190505050610103565b005b7f5c3eb8a1a7a5c7ef4c1a21921deb90256ddeb181f2665b8868ebfe35b8c9b9a782826040518080602001838152602001828103825284818151815260200191508051906020019080838360005b8381101561016d5780820151818401525b602081019050610151565b50505050905090810190601f16801561019a5780820380516001836020036101000a031916815260200191505b50935050505060405180910390a15b505056fea264697066735822122086a125077f6cec359aa907d52abb0e9a93cedbf5635c51e521fd7b9c83820e3c64736f6c63430006000033`
	var exceptAddress = `cd234a471b72ba2f1ccf0a70fcaba648a5eecd8d`
	logCode, logAddress = deployContract(t, memoryDB, bc, origin, binBytes, exceptAddress, exceptCode, 96759)
	// then call the contract with appendEntry function
	callWithPayload(t, memoryDB, bc, origin, logAddress, mustPack(logAbi, "appendEntry", "money", "10"), 2779, 0)
	require.Len(t, memoryDB.GetLog(), 1)
	entry := memoryDB.GetLog()[0]
	require.Equal(t, exceptAddress, fmt.Sprintf("%x", entry.Address.Bytes()))
}
