package model

// import (
// 	""
// )

// type Application struct {
// 	AppId struct{
// 		LocalStrategyEvents map[uint]uint8
// 		CloudStrategyEvents map[uint]uint8
// 		Servers map[string]struct{
// 			Server string
// 			Cloud string
// 			Cost uint8
// 			Status struct {
// 				MemLoad float32
// 				CpuLoad float32
// 				TimeStamp uint
// 				StateEvents map[uint]uint8
// 			}
// 		}
// 	}
// }

type Application map[string]struct {
	LocalStrategyEvents map[uint]uint8
	CloudStrategyEvents map[uint]uint8
	Servers             map[string]struct {
		Server string
		Cloud  string
		Cost   uint8
		Status struct {
			MemLoad   float32
			CpuLoad   float32
			TimeStamp uint
			// StateEvents map[uint]uint8
		}
	}
}

func NewApplication() *Application {
	a := new(Application)
	return a
}
