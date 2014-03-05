package helpers

import (
	"os"
)

const HYDRA_URI string = "http://127.0.0.1:8082"

func RunHydraInStandaloneAndReturnProcess() *os.Process {
	procAttr := new(os.ProcAttr)
	procAttr.Files = []*os.File{nil, os.Stdout, os.Stderr}
	// TODO: Testing force
	args := []string{"hydra", "-name=node1", "-f", "-data-dir=/tmp/node1"}

	process, err := os.StartProcess(HydraBinPath, args, procAttr)
	// It("should be running successfully", func() {
	// 	Expect(err).NotTo(HaveOccurred())
	// })
	if err != nil {
		// TODO
		// GinkgoT().Fatal("start process failed:" + err.Error())
		panic("start process failed:" + err.Error())
	}

	return process
}

func KillHydraProcess(process *os.Process) {
	process.Kill()
}
