package types

import "fmt"

var ErrServerShuttingDown = fmt.Errorf("server is shutting down")

var ErrInvalidUserKeySecret = fmt.Errorf("User key or secret is invalid")

var ErrJobStoppedByUser = fmt.Errorf("user stopped job")

var ErrProcessingRequest = fmt.Errorf("error processing request")
