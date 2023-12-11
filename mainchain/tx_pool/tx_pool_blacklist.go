/*
 *  Copyright 2023 KardiaChain
 *  This file is part of the go-kardia library.
 *
 *  The go-kardia library is free software: you can redistribute it and/or modify
 *  it under the terms of the GNU Lesser General Public License as published by
 *  the Free Software Foundation, either version 3 of the License, or
 *  (at your option) any later version.
 *
 *  The go-kardia library is distributed in the hope that it will be useful,
 *  but WITHOUT ANY WARRANTY; without even the implied warranty of
 *  MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
 *  GNU Lesser General Public License for more details.
 *
 *  You should have received a copy of the GNU Lesser General Public License
 *  along with the go-kardia library. If not, see <http://www.gnu.org/licenses/>.
 */

package tx_pool

import (
	"bufio"
	"io"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/kardiachain/go-kardia/lib/common"
	"github.com/kardiachain/go-kardia/lib/log"
)

const (
	// UpdateBlacklistInterval blocks since last update
	UpdateBlacklistInterval        uint64 = 50
	blacklistURL                          = "https://raw.githubusercontent.com/kardiachain/consensus/main/notes"
	InitialBlacklistRequestTimeout        = 1 * time.Second
	BlacklistRequestTimeout               = 2 * time.Second
)

// Blacklisted contains the blacklisted senders
var Blacklisted = make(map[string]bool)

// UpdateBlacklistLocal read and update the current blacklist
// from local file
func (pool *TxPool) UpdateBlacklistLocal() error {
	filePath := pool.config.BlacklistPath

	// Check file exists
	if _, err := os.Stat(filePath); err != nil {
		log.Warn("blacklist.txt does not exist", "err", err)
		return err
	}

	// open the file
	file, err := os.Open(filePath)
	if err != nil {
		log.Warn("Cannot open blacklist file", "err", err)
		return err
	}

	// read the file line-by-line
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		Blacklisted[common.HexToAddress(scanner.Text()).Hex()] = true
	}

	if err := scanner.Err(); err != nil {
		log.Warn("Cannot get address in blacklist file", "err", err)
		return err
	}

	return file.Close()
}

// UpdateBlacklistRemote fetch and update the current blacklist
// from remote URL
func UpdateBlacklistRemote(timeout time.Duration) error {
	httpClient := http.Client{Timeout: timeout}
	resp, err := httpClient.Get(blacklistURL)
	if err != nil {
		log.Warn("Cannot get blacklisted addresses", "err", err)
		return err
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Warn("Cannot import blacklisted addresses", "err", err)
		return err
	}

	blacklisted := strings.Split(string(body), "\n")
	for _, str := range blacklisted {
		Blacklisted[common.HexToAddress(str).Hex()] = true
	}

	return nil
}
