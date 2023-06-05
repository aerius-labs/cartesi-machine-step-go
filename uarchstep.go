package cartesi_machine_step

type IUArchStep interface {
	Step() (IUArchState, bool)
}

type UArchStep struct {
}

func (s *UArchStep) Step(state State) (uint64, bool) {
	ucycle := state.readCycle()

	if state.readHaltFlag() {
		if state.AccessLogs.Current == uint64(len(state.AccessLogs.Logs)) {
			panic("access pointer should match accesses length when halt")
		}
		return ucycle, true
	}

	if ucycle == ^uint64(0) {
		if state.AccessLogs.Current == uint64(len(state.AccessLogs.Logs)) {
			panic("access pointer should match accesses length when cycle is uint64.max")
		}
		return ucycle, false
	}

	upc := state.readPc()
	insn := readUint32(state, upc)
	uArchExecuteInsn(state, insn, upc)

	ucycle++
	state.writeCycle(ucycle)

	if state.AccessLogs.Current == uint64(len(state.AccessLogs.Logs)) {
		panic("access pointer should match accesses length")
	}

	return ucycle, false
}
