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
	name  string
	store *sync.Map
}

func NewMapDb(name string, store *sync.Map) MapDb {
	return MapDb{name, store}
}

var globalTxnLock *sync.Mutex
var globalTxnToken util.Uuid
var globalTxnTokenRwLock *sync.RWMutex

type Txn struct {
	context string
	token   util.Uuid
}

func (s Txn) invalidTokenError() error {
	return ErrInvalidTransaction.Make(nil, s.context, s.token, globalTxnToken)
}

type Any = interface{}

var (
	ErrInvalidTransaction = util.NewErrKind("txn context %v - invalid token %v != %v")
	ErrDuplicateKey       = util.NewErrKind("database %v - duplicate key \"%v\"")
	ErrRecordNotFound     = util.NewErrKind("database %v - record not found with key \"%v\"")
)

func BeginTxn(context string) Txn {
	globalTxnLock.Lock()
	globalTxnTokenRwLock.Lock()
	defer globalTxnTokenRwLock.Unlock()
	globalTxnToken = util.NewUuid()
	return Txn{context, globalTxnToken}
}

func EndTxn(txn Txn) error {
	globalTxnTokenRwLock.RLock()
	defer globalTxnTokenRwLock.RUnlock()
	if txn.token != globalTxnToken {
		return txn.invalidTokenError()
	}
	globalTxnLock.Unlock()
	return nil
}

func (s MapDb) Create(key Any, value Any, txn Txn) error {
	if txn.token != globalTxnToken {
		return txn.invalidTokenError()
	}
	_, loaded := s.store.LoadOrStore(key, value)
	if loaded {
		return ErrDuplicateKey.Make(nil, key)
	}
	return nil
}

func (s MapDb) Read(key Any) (Any, error) {
	value, loaded := s.store.Load(key)
	if loaded {
		return nil, ErrRecordNotFound.Make(nil, s.name, key)
	}
	return value, nil
}

func (s MapDb) Update(key Any, value Any, txn Txn) error {
	if txn.token != globalTxnToken {
		return txn.invalidTokenError()
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
		return txn.invalidTokenError()
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
