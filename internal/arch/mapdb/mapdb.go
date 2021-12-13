/*
 * Copyright © 2021 Paulo Villela. All rights reserved.
 * Use of this source code is governed by the Apache 2.0 license
 * that can be found in the LICENSE file.
 */

package mapdb

import (
	"github.com/pvillela/go-foa-realworld/internal/arch/db"
	"github.com/pvillela/go-foa-realworld/internal/arch/errx"
	"github.com/pvillela/go-foa-realworld/internal/arch/util"
	"github.com/sirupsen/logrus"
	"sync"
)

type Any = interface{}

type MapDb struct {
	Name  string
	Store *sync.Map
}

func NewMapDb(name string, store *sync.Map) MapDb {
	return MapDb{name, store}
}

var globalTxnLock *sync.Mutex
var globalTxnToken util.Uuid
var globalTxnTokenRwLock *sync.RWMutex

type TxnMapDb struct {
	context string
	token   util.Uuid
}

// Validate is part of db.Txn interface implementation.
func (txn TxnMapDb) Validate() error {
	globalTxnTokenRwLock.RLock()
	defer globalTxnTokenRwLock.RUnlock()
	if txn.token != globalTxnToken {
		return ErrInvalidTransaction.Make(nil, txn.context, txn.token, globalTxnToken)
	}
	return nil
}

func BeginTxn(context string) db.Txn {
	globalTxnLock.Lock()
	globalTxnTokenRwLock.Lock()
	defer globalTxnTokenRwLock.Unlock()
	globalTxnToken = util.NewUuid()
	return TxnMapDb{context, globalTxnToken}
}

// End is part of db.Txn interface implementation.
func (t TxnMapDb) End() {
	err := t.Validate()
	if err != nil {
		logrus.Error(err)
		return
	}
	globalTxnLock.Unlock()
	return
}

var (
	ErrInvalidTransaction = errx.NewKind("txn context %v - invalid token %v != %v")
	ErrDuplicateKey       = errx.NewKind("database %v - duplicate key \"%v\"")
	ErrRecordNotFound     = errx.NewKind("database %v - record not found with key \"%v\"")
)

func (s MapDb) Create(key Any, value Any, txn db.Txn) error {
	if err := txn.Validate(); err != nil {
		return err
	}
	if _, loaded := s.store.LoadOrStore(key, value); loaded {
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

func (s MapDb) Update(key Any, value Any, txn db.Txn) error {
	if err := txn.Validate(); err != nil {
		return err
	}

	// make sure record exists
	if _, err := s.Read(key); err != nil {
		return err
	}

	s.store.Store(key, value)
	return nil
}

func (s MapDb) Delete(key Any, txn db.Txn) error {
	if err := txn.Validate(); err != nil {
		return err
	}

	// make sure record exists
	if _, err := s.Read(key); err != nil {
		return err
	}

	s.store.Delete(key)
	return nil
}

func (s MapDb) Range(f func(key, value interface{}) bool) {
	s.store.Range(f)
}

func (s MapDb) FindFirst(pred func(Any, Any) bool) (retVal Any, found bool) {
	f := func(key, value Any) bool {
		if pred(key, value) {
			retVal = value
			found = true
			return false
		}
		return true
	}
	s.store.Range(f)
	return
}

func (s MapDb) FindAll(pred func(Any, Any) bool) []Any {
	retVals := make([]Any, 10)
	f := func(key, value Any) bool {
		if pred(key, value) {
			retVals = append(retVals, value)
		}
		return true
	}
	s.store.Range(f)
	return retVals
}
