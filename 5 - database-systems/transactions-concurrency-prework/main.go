package main

import (
	"errors"
)

type lockManager struct {
	// keyed on the names of the objects being locked (e.g. table name)
	locks map[string]*lockEntry

	// txs that have started unlocking (i.e. entered "shrinking" phase of 2PL) and can no longer acquire a lock
	shrinking map[string]bool
}

func newLockManager() *lockManager {
	return &lockManager{
		locks:     make(map[string]*lockEntry),
		shrinking: make(map[string]bool),
	}
}

type lockEntry struct {
	// names of txs granted access to the lock; shared lock can be granted to multiple txs
	granted map[string]bool

	mode lockMode

	// names of txs (in order of arrival) waiting to acquire a lock for this object
	queue []*lockRequest
}

type lockRequest struct {
	txId string
	mode lockMode
}

type lockMode uint8

const (
	None lockMode = iota
	Shared
	Exclusive
)

func (lm *lockManager) lock(transaction string, object string, mode lockMode) error {
	if mode == None || mode > Exclusive {
		return errors.New("invalid lock mode")

	}

	if lm.shrinking[transaction] {
		return errors.New("transaction no longer allowed to acquire locks")
	}

	entry, ok := lm.locks[object]
	// if there has never been a lock (or there are no longer any locks) on the object
	if !ok || entry.mode == None {
		lm.locks[object] = newLockEntry(transaction, mode)
		return nil
	}

	// if a new transaction is requesting another Shared lock, grant it
	if mode == Shared && entry.mode == Shared {
		entry.granted[transaction] = true
		return nil
	}

	// if requesting mode or existing mode is exclusive, add to wait queue and block
	entry.queue = append(entry.queue, newLockRequest(transaction, mode))
	// Block thread??
	return nil

}

func newLockRequest(transaction string, mode lockMode) *lockRequest {
	return &lockRequest{
		txId: transaction,
		mode: mode,
	}
}

func newLockEntry(transaction string, mode lockMode) *lockEntry {
	return &lockEntry{
		granted: map[string]bool{transaction: true},
		mode:    mode,
		queue:   []*lockRequest{},
	}
}

func (lm *lockManager) unlock(transaction string, object string) error {
	entry, ok := lm.locks[object]
	if !ok {
		return errors.New("there is nothing to unlock for this object")
	}

	// remove matching transaction from granted set
	for tx := range entry.granted {
		if tx == transaction {
			delete(entry.granted, tx)
			break
		}
	}

	// if there are no more transactions with the lock, update mode accordingly
	if len(entry.granted) == 0 {
		entry.mode = None
	}

	return nil
}
