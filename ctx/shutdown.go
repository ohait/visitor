package ctx

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"runtime/pprof"
	"syscall"
)

var Shutdown = initShutdown()

func initShutdown() chan struct{} {
	return make(chan struct{})
}

func WaitForSignal(cancel context.CancelFunc) {
	sigchan := make(chan os.Signal, 10)
	signal.Notify(sigchan, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)

	sig := <-sigchan
	fmt.Fprintf(os.Stderr, "SIG %q", sig.String())
	close(Shutdown)

	sig = <-sigchan
	fmt.Fprintf(os.Stderr, "SIG %q", sig.String())
	pprof.Lookup("goroutine").WriteTo(os.Stdout, 1)
	cancel()

	sig = <-sigchan
	fmt.Fprintf(os.Stderr, "SIG %q", sig.String())
	pprof.Lookup("goroutine").WriteTo(os.Stdout, 1)
	os.Exit(-1)
}
