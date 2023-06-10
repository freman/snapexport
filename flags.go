// Copyright 2023 Shannon Wynter
//
// This software may be modified and distributed under the terms
// of the MIT license.  See the LICENSE file for details.

package main

import (
	"flag"
	"fmt"
	"os"
	"time"
)

func strEnvFlag(name, env, desc string) *string {
	return flag.String(name, os.Getenv(env), fmt.Sprintf("%s (env: %s)", desc, env))
}

func parseDateFlag(name, in string) (resp time.Time) {
	resp, err := time.Parse(dateFormat, in)
	if err != nil {
		fmt.Println("Unable to parse", name, "date:", err.Error())
		flag.PrintDefaults()
		os.Exit(1)
	}

	return
}

func chkRequiredFlags(args ...string) {
	l := len(args)
	if l%2 != 0 {
		panic("arguments to be divisible by 2 (value, name)")
	}

	var errd bool
	for i := 0; i < l; i += 2 {
		if args[i] == "" {
			fmt.Println("Flag", args[i+1], "is required")
			errd = true
		}
	}

	if errd {
		flag.PrintDefaults()
		os.Exit(1)
	}
}
