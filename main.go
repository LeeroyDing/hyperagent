package main

import (
"log"
"github.com/LeeroyDing/hyperagent/internal/cmd"
)

func main() {
if err := cmd.Execute(); err != nil {
log.Fatal(err)
}
}
