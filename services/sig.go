package services

import (
  "syscall"
  "fmt"
)

func parseSignal(sig string) (syscall.Signal, error) {
  if sign, has := signals[sig]; has {
    return syscall.Signal(sign), nil
  }
  return syscall.Signal(0), fmt.Errorf("Could not parse signal %v", sig)
}