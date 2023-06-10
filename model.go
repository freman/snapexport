// Copyright 2023 Shannon Wynter
//
// This software may be modified and distributed under the terms
// of the MIT license.  See the LICENSE file for details.

package main

import (
	"encoding/json"
	"time"

	"github.com/tidwall/gjson"
)

type Response struct {
	Cases                  []Case `json:"cases"`
	LinkToNextSetOfResults string `json:"linkToNextSetOfResults"`
}

type Case struct {
	json.RawMessage
	ID      string
	Created time.Time
}

func (c *Case) UnmarshalJSON(in []byte) error {
	c.RawMessage = in
	c.ID = gjson.GetBytes(in, "id").String()
	c.Created = time.Unix(gjson.GetBytes(in, "created_at_seconds").Int(), 0)

	return nil
}

func (c Case) MarshalJSON() ([]byte, error) {
	return c.RawMessage, nil
}
