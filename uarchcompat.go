package cartesi_machine_step

func (state *State) readWord(address uint64) uint64 {
	var res = state.StateI.ReadWord(state.AccessLogs, address)
	state.AccessLogs.Current += 1
	return res
}

func (state *State) readPc() uint64 {
	var res = state.StateI.ReadPc(state.AccessLogs)
	state.AccessLogs.Current += 1
	return res
}

func (state *State) readHaltFlag() bool {
	var res = state.StateI.ReadHaltFlag(state.AccessLogs)
	state.AccessLogs.Current += 1
	return res
}

func (state *State) readCycle() uint64 {
	var res = state.StateI.ReadCycle(state.AccessLogs)
	state.AccessLogs.Current += 1
	return res
}

func (state *State) writeCycle(value uint64) {
	state.StateI.WriteCycle(state.AccessLogs, value)
	state.AccessLogs.Current += 1
}

func (state *State) readX(index uint64) uint64 {
	var res = state.StateI.ReadX(state.AccessLogs, index)
	state.AccessLogs.Current += 1
	return res
}

func (state *State) writeWord(address uint64, value uint64) {
	state.StateI.writeWord(state.AccessLogs, address, value)
	state.AccessLogs.Current += 1
}

func (state *State) writeX(index uint64, value uint64) {
	state.StateI.WriteX(state.AccessLogs, index, value)
	state.AccessLogs.Current += 1
}

func (state *State) writePc(value uint64) {
	state.StateI.WritePc(state.AccessLogs, value)
	state.AccessLogs.Current += 1
}

func uint64ShiftRight(v uint64, count uint32) uint64 {
	return v >> (count & 0x3f)
}

func uint64ShiftLeft(v uint64, count uint32) uint64 {
	return v << (count & 0x3f)
}

func int64ShiftRight(v int64, count uint32) int64 {
	return v >> (count & 0x3f)
}

func uint32ShiftRight(v uint32, count uint32) uint32 {
	return v >> (count & 0x1f)
}

func uint32ShiftLeft(v uint32, count uint32) uint32 {
	return v << (count & 0x1f)
}

func int32ShiftRight(v int32, count uint32) int32 {
	return v >> (count & 0x1f)
}
