go-schroot
==========

`import "pault.ag/go/schroot"`

```go
package main

import (
	"log"

	"pault.ag/go/schroot"
)

func main() {
	chroot, err := schroot.NewSchroot("unstable")
	if err != nil {
		log.Fatalf("%s\n", err)
	}
	defer chroot.End()
	cmd, err := chroot.Command("ls", "/var")
	if err != nil {
		log.Fatalf("%s\n", err)
	}
	out, err := cmd.Output()
	if err != nil {
		log.Fatalf("%s\n", err)
		return
	}
	log.Printf("%s\n", out)
```
