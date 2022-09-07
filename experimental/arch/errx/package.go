/*
 *  Copyright © 2021 Paulo Villela. All rights reserved.
 *  Use of this source code is governed by the Apache 2.0 license
 *  that can be found in the LICENSE file.
 */

// This package defines and implements an error interface that includes:
// definition of an error kind that can be referenced in error handlers;
// a chain of errors that are the causes of the current error;
// a stack trace.
// The stack trace is based on runtime/debug/Stack.
package errx
