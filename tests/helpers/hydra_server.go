package helpers

import (
	"os"
)

// const PRIVATE_HYDRA_URI string = "127.0.0.1:8082"

func RunHydraInStandaloneAndReturnProcess(privateAddr string) *os.Process {
	procAttr := new(os.ProcAttr)
	procAttr.Files = []*os.File{nil, os.Stdout, os.Stderr}
	args := []string{"hydra", "-f", "-name=node0", "-private-addr=" + privateAddr, "-data-dir=/tmp/node0"}

	process, err := os.StartProcess(HydraBinPath, args, procAttr)
	if err != nil {
		panic("start process failed:" + err.Error())
	}

	return process
}

func KillHydraProcess(process *os.Process) {
	process.Kill()
}
