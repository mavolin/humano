<div align="center">
<h1>humano</h1>

[![Go Reference](https://pkg.go.dev/badge/github.com/mavolin/humano.svg)](https://pkg.go.dev/github.com/mavolin/humano)
[![GitHub Workflow Status (branch)](https://img.shields.io/github/workflow/status/mavolin/adam/Test/develop?label=tests)](https://github.com/mavolin/humano/actions)
[![codecov](https://codecov.io/gh/mavolin/adam/branch/develop/graph/badge.svg?token=3qRIAudu4r)](https://codecov.io/gh/mavolin/adam)
[![Go Report Card](https://goreportcard.com/badge/github.com/mavolin/adam)](https://goreportcard.com/report/github.com/mavolin/adam)
[![License](https://img.shields.io/github/license/mavolin/dismock)](https://github.com/mavolin/dismock/blob/v2/LICENSE)
</div>

---

## About

Humano is a small utility for [adam](https://github.com/mavolin/adam), that allows you to imitate typing, by delaying messages.

## Example

```go
package main

import (
	"log"
	"os"
	"os/signal"

	"github.com/mavolin/adam/pkg/bot"

	"github.com/mavolin/humano/pkg/replier"
)

func main() {
	b, err := bot.New(bot.Options{
		Token: os.Getenv("DISCORD_BOT_TOKEN"),
	})
	if err != nil {
		log.Fatal(err)
	}

	b.MustAddMiddleware(replier.NewMiddleware(&replier.Options{}))

	// add commands

	if err := b.Open(); err != nil {
		log.Fatal(err)
	}

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, os.Interrupt)

	<-sig

	if err := b.Close(); err != nil {
		log.Fatal(err)
	}
}
```

```go
package mycommand

...

func (c *MyCommand) Invoke(s *state.State, ctx *plugin.Context) (interface{}, error) {
	_, err := replier.Reply(ctx, "Wumpus!")
	return nil, err
}
```
