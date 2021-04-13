/*
 *  Copyright © 2021 Paulo Villela. All rights reserved.
 *  Use of this source code is governed by the Apache 2.0 license
 *  that can be found in the LICENSE file.
 */

package mapdb

import (
	"github.com/pvillela/go-foa-realworld/internal/arch/util"
	"sync"
)

type MapDb struct {
	store *sync.Map
}

var globalTxnLock *sync.Mutex
var globalTxnToken util.Uuid

type Txn struct {
	token util.Uuid
}

type Any = interface{}

var (
	ErrInvalidTransaction = util.NewErrKind("invalid transaction token %v != %v")
	ErrDuplicateKey       = util.NewErrKind("duplicate key \"%v\"")
	ErrRecordNotFound     = util.NewErrKind("record not found with key \"%v\"")
)

func BeginTxn() Txn {
	globalTxnLock.Lock()
	token := util.NewUuid()
	return Txn{token}
}

func EndTxn(txn Txn) error {
	if txn.token != globalTxnToken {
		return ErrInvalidTransaction.Make(txn.token, globalTxnToken)
	}
	globalTxnLock.Unlock()
	return nil
}

func (s MapDb) Create(key Any, value Any, txn Txn) error {
	if txn.token != globalTxnToken {
		return ErrInvalidTransaction.Make(txn.token, globalTxnToken)
	}
	_, loaded := s.store.LoadOrStore(key, value)
	if loaded {
		return ErrDuplicateKey.Make(key)
	}
	return nil
}

func (s MapDb) Read(key Any) (Any, error) {
	value, loaded := s.store.Load(key)
	if loaded {
		return nil, ErrRecordNotFound.Make(key)
	}
	return value, nil
}

func (s MapDb) Update(key Any, value Any, txn Txn) error {
	if txn.token != globalTxnToken {
		return ErrInvalidTransaction.Make(txn.token, globalTxnToken)
	}

	_, err := s.Read(key) // make sure record exists
	if err != nil {
		return err
	}

	s.store.Store(key, value)
	return nil
}

func (s MapDb) Delete(key Any, txn Txn) error {
	if txn.token != globalTxnToken {
		return ErrInvalidTransaction.Make(txn.token, globalTxnToken)
	}

	_, err := s.Read(key) // make sure record exists
	if err != nil {
		return err
	}

	s.store.Delete(key)
	return nil
}

func (s MapDb) Range(f func(key, value interface{}) bool) {
	s.store.Range(f)
}
