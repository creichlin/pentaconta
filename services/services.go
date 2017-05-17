package services

import (
	"gitlab.com/creichlin/pentaconta/logger"
)

type Services struct {
	Logs        logger.Logger
	Executors   map[string]*Executor
	FSListeners map[string]*FSListener
}
