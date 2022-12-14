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
	abi "github.com/sunvim/evm/abi"
	"github.com/sunvim/evm/db"
	"github.com/sunvim/evm/example"
	"github.com/sunvim/evm/util"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

var (
	balanceBin       = "sols/Balance_sol_Balance.bin"
	balanceAbi       = "sols/Balance_sol_Balance.abi"
	balanceCode      []byte
	balanceAddress   evm.Address
	benckmarkBalance = true
)

func TestBalanceSol(t *testing.T) {
	// first create the contract
	binBytes, err := util.ReadBinFile(balanceBin)
	require.NoError(t, err)
	bc := example.NewBlockchain()
	memoryDB := db.NewMemory(bc.NewAccount)
	var origin = example.HexToAddress("6ac7ea33f8831ea9dcc53393aaa88b25a785dbf0")
	var exceptCode = `60806040523480156100115760006000fd5b506004361061005c5760003560e01c80631003e2d21461006257806327ee58a6146100a5578063370158ea146100e857806360fe47b1146101395780636d4ce63c146101805761005c565b60006000fd5b61008f600480360360208110156100795760006000fd5b810190808035906020019092919050505061019e565b6040518082815260200191505060405180910390f35b6100d2600480360360208110156100bc5760006000fd5b81019080803590602001909291905050506101c6565b6040518082815260200191505060405180910390f35b6100f06101ee565b604051808373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020018281526020019250505060405180910390f35b610166600480360360208110156101505760006000fd5b8101908080359060200190929190505050610209565b604051808215151515815260200191505060405180910390f35b610188610225565b6040518082815260200191505060405180910390f35b6000816000600082828250540192505081909090555060006000505490506101c1565b919050565b6000816000600082828250540392505081909090555060006000505490506101e9565b919050565b600060003360006000505481915091509150610205565b9091565b600081600060005081909090555060019050610220565b919050565b60006000600050549050610234565b9056fea26469706673582212206d1f7e72f2d26816fe48ff60de6fa42d7b6fb40fc1603494b8c63221cd3c2c2364736f6c63430006000033`
	var exceptAddress = `cd234a471b72ba2f1ccf0a70fcaba648a5eecd8d`
	balanceCode, balanceAddress = deployContract(t, memoryDB, bc, origin, binBytes, exceptAddress, exceptCode, 144400)
	// then call the contract with get function
	result := callBalance(t, memoryDB, bc, origin, mustPack(balanceAbi, "get"), 1096) // except 10
	require.EqualValues(t, []string{"10"}, mustUnpack(balanceAbi, "get", result))
	// then set value to 20
	result = callBalance(t, memoryDB, bc, origin, mustPack(balanceAbi, "set", "20"), 5393) // except true
	require.EqualValues(t, []string{"true"}, mustUnpack(balanceAbi, "set", result))
	// then get
	result = callBalance(t, memoryDB, bc, origin, mustPack(balanceAbi, "get"), 1096) // except 20
	require.EqualValues(t, []string{"20"}, mustUnpack(balanceAbi, "get", result))
	// then add
	callBalance(t, memoryDB, bc, origin, mustPack(balanceAbi, "add", "10"), 6939) // except 30
	// then get
	result = callBalance(t, memoryDB, bc, origin, mustPack(balanceAbi, "get"), 1096) // except 30
	require.EqualValues(t, []string{"30"}, mustUnpack(balanceAbi, "get", result))
	// info
	result = callBalance(t, memoryDB, bc, origin, mustPack(balanceAbi, "info"), 1105) // except "6ac7ea33f8831ea9dcc53393aaa88b25a785dbf0", "30"
	require.EqualValues(t, []string{"6ac7ea33f8831ea9dcc53393aaa88b25a785dbf0", "30"}, mustUnpack(balanceAbi, "info", result))
	// define temporary address for testing
	var temporarySender = RandomAddress()
	var temporaryBC = NewBlockchain()
	abi.SetAddressParser(temporarySender.Length(), func(bytes []byte) string {
		return BytesToAddress(bytes).String()
	}, nil)
	result = callBalance(t, memoryDB, temporaryBC, temporarySender, mustPack(balanceAbi, "info"), 1105)
	require.EqualValues(t, []string{temporarySender.String(), "30"}, mustUnpack(balanceAbi, "info", result))
	abi.ResetAddressParser()
}

// you can set gasCost to 0 if you do not want to compare gasCost
func callBalance(t *testing.T, db evm.DB, bc evm.Blockchain, caller evm.Address, payload []byte, gasCost uint64) []byte {
	var gasQuota uint64 = 1000000
	var gas = gasQuota
	output, err := evm.New(bc, db, &evm.Context{
		Input: payload,
		Value: 0,
		Gas:   &gas,
	}).Call(caller, balanceAddress, balanceCode)
	require.NoError(t, err)
	if gasCost != 0 {
		require.EqualValues(t, gasCost, gasQuota-gas, fmt.Sprintf("Except gas cost %d other than %d", gasCost, gasQuota-gas))
	}
	return output
}
