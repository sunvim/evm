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

package evm

import (
	"github.com/sunvim/evm/crypto"
	"github.com/sunvim/evm/rlp"
)

// This file defines some default funcion if the user do not want to implement it by themself.

// defaultCreateAddress is the default implementation of CreateAddress
func defaultCreateAddress(caller Address, nonce uint64, toAddressFunc func(bytes []byte) Address) Address {
	data, _ := rlp.EncodeToBytes([]interface{}{caller, nonce})
	bytes := crypto.Keccak256(data)[12:]
	return toAddressFunc(bytes)
}

func defaultCreate2Address(caller Address, salt, code []byte, toAddressFunc func(bytes []byte) Address) Address {
	bytes := crypto.Keccak256([]byte{0xff}, caller.Bytes(), salt[:], crypto.Keccak256(code))[12:]
	return toAddressFunc(bytes)
}
