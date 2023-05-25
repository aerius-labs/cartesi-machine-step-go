package cartesi_machine_step

type State struct {
	StateI     IUArchState
	AccessLogs AccessLogs
}

type IUArchState interface {
	ReadCycle(a AccessLogs) uint64
	ReadHaltFlag(a AccessLogs) bool
	ReadPc(a AccessLogs) uint64
	ReadWord(a AccessLogs, address uint64) uint64
	ReadX(a AccessLogs, reg uint64) uint64
	WriteCycle(a AccessLogs, value uint64)
	WritePc(a AccessLogs, value uint64)
	writeWord(a AccessLogs, address uint64, value uint64)
	WriteX(a AccessLogs, index uint64, value uint64)
}
