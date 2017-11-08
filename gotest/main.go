package main

import (
	"fmt"
	"os/exec"
)

func main() {

	cmd := exec.Command("/usr/local/bin/go" ,"test", "-cover", "./..." ,".")
	if err := cmd.Run(); err != nil {
		fmt.Printf("failed:%v",err)
		return
	}
	fmt.Println("succeeded")
}
