// Copyright 2023 Shannon Wynter
//
// This software may be modified and distributed under the terms
// of the MIT license.  See the LICENSE file for details.

package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/cenkalti/backoff"
)

type client struct {
	*http.Client
	token   string
	retries uint64
}

func (c *client) fetch(url string) (jsonResp Response, err error) {
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return jsonResp, fmt.Errorf("failed to create request for %s: %w", url, err)
	}

	req.Header.Add("Authorization", c.token)

	err = backoff.RetryNotify(func() error {
		resp, err := c.Client.Do(req)
		if err != nil {
			return err
		}
		defer resp.Body.Close()

		if err := json.NewDecoder(resp.Body).Decode(&jsonResp); err != nil {
			return backoff.Permanent(fmt.Errorf("unable to parse json: %w", err))
		}

		return nil
	}, backoff.WithMaxRetries(backoff.NewExponentialBackOff(), c.retries), notifyFunc)

	return jsonResp, err
}

func notifyFunc(err error, dur time.Duration) {
	fmt.Println("Retrying request in", dur, "due to failure", err.Error())
}
