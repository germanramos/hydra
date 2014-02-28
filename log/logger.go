package log

import (
	golog "github.com/innotech/hydra/vendors/github.com/coreos/go-log/log"
)

var logger *golog.Logger := log.New("hydra", false,
		log.CombinedSink(os.Stdout, "[%s] %s %-9s | %s\n", []string{"prefix", "time", "priority", "message"}))
