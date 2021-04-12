/*
 *  Copyright © 2021 Paulo Villela. All rights reserved.
 *  Use of this source code is governed by the Apache 2.0 license
 *  that can be found in the LICENSE file.
 */

package crypto

import (
	"crypto/sha256"
	"fmt"
)

func Hash(salt []byte, text string) string {
	h := sha256.New()
	bytes := make([]byte, len(salt)+len(text))
	bytes = append(bytes, salt...)
	bytes = append(bytes, []byte(text)...)
	h.Write(bytes)
	return fmt.Sprintf("%x", h.Sum(nil))
}
