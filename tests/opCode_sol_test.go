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
	opCodeBin     = "sols/Opcode_sol_OpCode.bin"
	opCodeAbi     = "sols/Opcode_sol_OpCode.abi"
	opCode        []byte
	opCodeAddress evm.Address
)

func TestOpCOde(t *testing.T) {
	binBytes, err := util.ReadBinFile(opCodeBin)
	require.NoError(t, err)
	bc := example.NewBlockchain()
	memoryDB := db.NewMemory(bc.NewAccount)
	var origin = example.HexToAddress("6ac7ea33f8831ea9dcc53393aaa88b25a785dbf0")
	//	var exceptCode = `60806040523480156100115760006000fd5b50610017565b6102db806100266000396000f3fe60806040523480156100115760006000fd5b50600436106100985760003560e01c8063ab70fd6911610067578063ab70fd6914610142578063b6baffe314610160578063d1a82a9d1461017e578063df1f29ee146101c8578063f2c9ecd81461021257610098565b806312065fe01461009e578063188ec356146100bc57806338cc4831146100da578063a16963b31461012457610098565b60006000fd5b6100a6610230565b6040518082815260200191505060405180910390f35b6100c461023d565b6040518082815260200191505060405180910390f35b6100e261024a565b604051808273ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200191505060405180910390f35b61012c610257565b6040518082815260200191505060405180910390f35b61014a610264565b6040518082815260200191505060405180910390f35b610168610271565b6040518082815260200191505060405180910390f35b61018661027e565b604051808273ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200191505060405180910390f35b6101d061028b565b604051808273ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200191505060405180910390f35b61021a610298565b6040518082815260200191505060405180910390f35b600047905061023a565b90565b6000429050610247565b90565b6000309050610254565b90565b6000459050610261565b90565b60003a905061026e565b90565b600044905061027b565b90565b6000419050610288565b90565b6000329050610295565b90565b60004390506102a2565b9056fea264697066735822122015c6e882c3fdc1443e57fbd751c159a2579310479893e31b21797f9e2579ce4b64736f6c63430006000033`
	var exceptAddress = `cd234a471b72ba2f1ccf0a70fcaba648a5eecd8d`
	opCode, opCodeAddress = deployContract(t, memoryDB, bc, origin, binBytes, exceptAddress, "", 178432)
	callOp(t, memoryDB, bc, origin, mustPack(opCodeAbi, "signExtend"), 263)
	callOp(t, memoryDB, bc, origin, mustPack(opCodeAbi, "signExtend2"), 245)
	callOp(t, memoryDB, bc, origin, mustPack(opCodeAbi, "codeSize"), 237)
	callOp(t, memoryDB, bc, origin, mustPack(opCodeAbi, "returnDataSize"), 280)
	callOp(t, memoryDB, bc, origin, mustPack(opCodeAbi, "blockHash"), 299)
	callOp(t, memoryDB, bc, origin, mustPack(opCodeAbi, "mStore8"), 271)
	callOp(t, memoryDB, bc, origin, mustPack(opCodeAbi, "testPC"), 280)
	callOp(t, memoryDB, bc, origin, mustPack(opCodeAbi, "testMSize"), 302)
	callOp(t, memoryDB, bc, origin, mustPack(opCodeAbi, "testGas"), 281)
	callOp(t, memoryDB, bc, origin, mustPack(opCodeAbi, "testStop"), 223)
	callOpIgnoreError(t, memoryDB, bc, origin, mustPack(opCodeAbi, "testRevert"), 186)
	callOpIgnoreError(t, memoryDB, bc, origin, mustPack(opCodeAbi, "testInvalid"), 10000)

}

func callOp(t *testing.T, db evm.DB, bc evm.Blockchain, caller evm.Address, payload []byte, gasCost uint64) {
	var gasQuota uint64 = 10000
	var gas = gasQuota
	output, err := evm.New(bc, db, &evm.Context{
		Input: payload,
		Value: 0,
		Gas:   &gas,
	}).Call(caller, opCodeAddress, opCode)
	require.NoError(t, err)
	if gasCost != 0 {
		require.EqualValues(t, gasCost, gasQuota-gas)
	}
	t.Log(output)
	t.Log(gasQuota - gas)
}

func callOpIgnoreError(t *testing.T, db evm.DB, bc evm.Blockchain, caller evm.Address, payload []byte, gasCost uint64) {
	var gasQuota uint64 = 10000
	var gas = gasQuota
	output, err := evm.New(bc, db, &evm.Context{
		Input: payload,
		Value: 0,
		Gas:   &gas,
	}).Call(caller, opCodeAddress, opCode)
	t.Log(err)
	if gasCost != 0 {
		require.EqualValues(t, gasCost, gasQuota-gas)
	}
	t.Log(output)
	t.Log(gasQuota - gas)
}
