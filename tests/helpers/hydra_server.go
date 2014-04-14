package helpers

import (
	"os"
)

// const PRIVATE_HYDRA_URI string = "127.0.0.1:8082"
const DATA_DIR_PATH = "/tmp/hydra_test/"

func RunHydraInStandaloneAndReturnProcess(args []string) *os.Process {
	procAttr := new(os.ProcAttr)
	procAttr.Files = []*os.File{nil, os.Stdout, os.Stderr}
	dataDirExits, err := existsPath(DATA_DIR_PATH)
	if err == nil && dataDirExits {
		os.RemoveAll(DATA_DIR_PATH)
	}
	args = append([]string{"hydra", "-f"}, args...)

	process, err := os.StartProcess(HydraBinPath, args, procAttr)
	if err != nil {
		panic("start process failed:" + err.Error())
	}

	return process
}

// func RunHydraInStandaloneAndReturnProcess(privateAddr string) *os.Process {
// 	procAttr := new(os.ProcAttr)
// 	procAttr.Files = []*os.File{nil, os.Stdout, os.Stderr}
// 	dataDirExits, err := existsPath(DATA_DIR_PATH)
// 	if err == nil && dataDirExits {
// 		os.RemoveAll(DATA_DIR_PATH)
// 	}
// 	args := []string{"hydra", "-f", "-name=node0", "-private-addr=" + privateAddr, "-data-dir=" + DATA_DIR_PATH, "-apps-file=fixtures/apps.empty.json"}

// 	process, err := os.StartProcess(HydraBinPath, args, procAttr)
// 	if err != nil {
// 		panic("start process failed:" + err.Error())
// 	}

// 	return process
// }

func KillHydraProcess(process *os.Process) {
	process.Kill()
}

func existsPath(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}
