// Copyright (C) 2024 Opsmate, Inc.
//
// This Source Code Form is subject to the terms of the Mozilla
// Public License, v. 2.0. If a copy of the MPL was not distributed
// with this file, You can obtain one at http://mozilla.org/MPL/2.0/.
//
// This software is distributed WITHOUT A WARRANTY OF ANY KIND.
// See the Mozilla Public License for details.

package monitor

import (
	"context"
	"time"

	"software.sslmate.com/src/certspotter/ct"
	"software.sslmate.com/src/certspotter/loglist"
	"software.sslmate.com/src/certspotter/merkletree"
)

type LogState struct {
	DownloadPosition *merkletree.CollapsedTree `json:"download_position"`
	VerifiedPosition *merkletree.CollapsedTree `json:"verified_position"`
	VerifiedSTH      *ct.SignedTreeHead        `json:"verified_sth"`
	LastSuccess      time.Time                 `json:"last_success"`
}

type StateProvider interface {
	// Initialize the state.  Called before any other method in this interface.
	// Idempotent: returns nil if the state is already initialized.
	Prepare(context.Context) error

	// Initialize the state for the given log.  Called before any other method
	// with the log ID.  Idempotent: returns nil if log state already initialized.
	PrepareLog(context.Context, LogID) error

	// Store log state for retrieval by LoadLogState.
	StoreLogState(context.Context, LogID, *LogState) error

	// Load log state that was previously stored with StoreLogState.
	// Returns nil, nil if StoreLogState has not been called yet for this log.
	LoadLogState(context.Context, LogID) (*LogState, error)

	// Store STH for retrieval by LoadSTHs.  If an STH with the same
	// timestamp and root hash is already stored, this STH can be ignored.
	StoreSTH(context.Context, LogID, *ct.SignedTreeHead) error

	// Load all STHs for this log previously stored with StoreSTH.
	// The returned slice must be sorted by tree size.
	LoadSTHs(context.Context, LogID) ([]*ct.SignedTreeHead, error)

	// Remove an STH so it is no longer returned by LoadSTHs.
	RemoveSTH(context.Context, LogID, *ct.SignedTreeHead) error

	// Called when a certificate matching the watch list is discovered.
	NotifyCert(context.Context, *DiscoveredCert) error

	// Called when certspotter fails to parse a log entry.
	NotifyMalformedEntry(context.Context, *LogEntry, error) error

	// Called when a health check fails.  The log is nil if the
	// feailure is not associated with a log.
	NotifyHealthCheckFailure(context.Context, *loglist.Log, HealthCheckFailure) error

	// Called when a non-fatal error occurs.  The log is nil if the error is
	// not associated with a log.  Note that most errors are transient, and
	// certspotter will retry the failed operation later.
	NotifyError(context.Context, *loglist.Log, error) error
}
