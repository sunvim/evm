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
	"testing"

	"github.com/stretchr/testify/require"
)

var (
	createBin = "sols/Create_sol_C.bin"
	createAbi = "sols/Create_sol_C.abi"
	CCode     []byte
	CAddress  evm.Address
)

func TestCreateSol(t *testing.T) {
	binBytes, err := util.ReadBinFile(createBin)
	require.NoError(t, err)
	bc := example.NewBlockchain()
	memoryDB := db.NewMemory(bc.NewAccount)
	var origin = example.HexToAddress("6ac7ea33f8831ea9dcc53393aaa88b25a785dbf0")
	var exceptAddress = `cd234a471b72ba2f1ccf0a70fcaba648a5eecd8d`
	CCode, CAddress = deployContract(t, memoryDB, bc, origin, binBytes, exceptAddress, "6080604052600436106100385760003560e01c80638dcd64cc1461003e57806395fe0e65146100b7578063afb6e5a1146100f457610038565b60006000fd5b610075600480360360408110156100555760006000fd5b81019080803590602001909291908035906020019092919050505061014f565b604051808273ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200191505060405180910390f35b3480156100c45760006000fd5b506100f2600480360360208110156100dc5760006000fd5b8101908080359060200190929190505050610221565b005b3480156101015760006000fd5b50610139600480360360408110156101195760006000fd5b8101908080359060200190929190803590602001909291905050506102e4565b6040518082815260200191505060405180910390f35b60006000828460405161016190610321565b808281526020019150506040518091039082f0905080158015610189573d600060003e3d6000fd5b5090508073ffffffffffffffffffffffffffffffffffffffff16630c55699c6040518163ffffffff1660e01b815260040160206040518083038186803b1580156101d35760006000fd5b505afa1580156101e8573d600060003e3d6000fd5b505050506040513d60208110156101ff5760006000fd5b8101908080519060200190929190505050508091505061021b56505b92915050565b60008160405161023090610321565b80828152602001915050604051809103906000f080158015610257573d600060003e3d6000fd5b5090508073ffffffffffffffffffffffffffffffffffffffff16630c55699c6040518163ffffffff1660e01b815260040160206040518083038186803b1580156102a15760006000fd5b505afa1580156102b6573d600060003e3d6000fd5b505050506040513d60208110156102cd5760006000fd5b810190808051906020019092919050505050505b50565b600060006102f8848461014f63ffffffff16565b90508073ffffffffffffffffffffffffffffffffffffffff163191505061031b56505b92915050565b6101258061032f8339019056fe6080604052604051610125380380610125833981810160405260208110156100275760006000fd5b81019080805190602001909291905050505b806000600050819090905550602060006000375b50610053565b60c4806100616000396000f3fe608060405234801560105760006000fd5b506004361060365760003560e01c80630c55699c14603c578063993a04b7146058576036565b60006000fd5b60426074565b6040518082815260200191505060405180910390f35b605e607d565b6040518082815260200191505060405180910390f35b60006000505481565b60006000600050549050608b565b9056fea2646970667358221220f1b371a69e1448dc014dd9e059a37e408373f64b097e12934a2dbe7b8ac3d25864736f6c63430006000033a264697066735822122013adea62c3a2f1f0b55929f719c8b8310cfd590b0dc49d4d281714e3697c820b64736f6c63430006000033", 344952)
	callCreate(t, memoryDB, bc, origin, mustPack(createAbi, "createAndGetBalance", "44", "0"), 95396)
}

func callCreate(t *testing.T, db evm.DB, bc evm.Blockchain, caller evm.Address, payload []byte, gasCost uint64) {
	var gasQuota uint64 = 1000000
	var gas = gasQuota
	output, err := evm.New(bc, db, &evm.Context{
		Input: payload,
		Value: 0,
		Gas:   &gas,
	}).Call(caller, CAddress, CCode)
	require.NoError(t, err)
	if gasCost != 0 {
		require.EqualValues(t, gasCost, gasQuota-gas)
	}

	t.Log(output)
	t.Log(gasQuota - gas)
}
