package cartesi_machine_step

import "errors"

func readUint64(state State, paddr uint64) (uint64, error) {
	if paddr&7 != 0 {
		return 0, errors.New("misaligned readUint64 address")
	}
	return state.readWord(paddr), nil
}

func readUint32(state State, paddr uint64) (uint32, error) {
	if paddr&3 != 0 {
		return 0, errors.New("misaligned readUint32 address")
	}
	palign := paddr & ^uint64(7)
	bitoffset := uint32ShiftLeft(uint32(paddr)&7, 3)
	val64, err := readUint64(state, palign)
	if err != nil {
		return 0, err
	}
	return uint32(uint64ShiftRight(val64, bitoffset)), nil
}

func readUint16(state State, paddr uint64) (uint16, error) {
	if paddr&1 != 0 {
		return 0, errors.New("misaligned readUint16 address")
	}
	palign := paddr & ^uint64(7)
	bitoffset := uint32ShiftLeft(uint32(paddr)&7, 3)
	val64, err := readUint64(state, palign)
	if err != nil {
		return 0, err
	}
	return uint16(uint64ShiftRight(val64, bitoffset)), nil
}

func readUint8(state State, paddr uint64) (uint8, error) {
	palign := paddr & ^uint64(7)
	bitoffset := uint32ShiftLeft(uint32(paddr)&7, 3)
	val64, err := readUint64(state, palign)
	if err != nil {
		return 0, err
	}
	return uint8(uint64ShiftRight(val64, bitoffset)), nil
}

func writeUint64(state State, paddr uint64, val uint64) error {
	if paddr&7 != 0 {
		return errors.New("misaligned writeUint64 address")
	}
	state.writeWord(paddr, val)
	return nil
}

func copyBits(from uint32, count uint32, to uint64, offset uint32) uint64 {
	if offset+count > 64 {
		panic("copyBits count exceeds limit of 64")
	}
	eraseMask := uint64ShiftLeft(1, count) - 1
	eraseMask = ^uint64ShiftLeft(eraseMask, offset)
	return uint64ShiftLeft(uint64(from), offset) | (to & eraseMask)
}

func writeUint32(state State, paddr uint64, val uint32) error {
	if paddr&3 != 0 {
		return errors.New("misaligned writeUint32 address")
	}
	palign := paddr & ^uint64(7)
	bitoffset := uint32ShiftLeft(uint32(paddr)&7, 3)
	oldval64, err := readUint64(state, palign)
	if err != nil {
		return err
	}
	newval64 := copyBits(uint32(val), 32, oldval64, bitoffset)
	return writeUint64(state, paddr, newval64)
}

func writeUint16(state State, paddr uint64, val uint16) error {
	if paddr&1 != 0 {
		return errors.New("misaligned writeUint16 address")
	}
	palign := paddr & ^uint64(7)
	bitoffset := uint32ShiftLeft(uint32(paddr)&7, 3)
	oldval64, err := readUint64(state, palign)
	if err != nil {
		return err
	}
	newval64 := copyBits(uint32(val), 16, oldval64, bitoffset)
	return writeUint64(state, palign, newval64)
}

func writeUint8(state State, paddr uint64, val uint8) error {
	palign := paddr & ^uint64(7)
	bitoffset := uint32ShiftLeft(uint32(paddr)&7, 3)
	oldval64, err := readUint64(state, palign)
	if err != nil {
		return err
	}
	newval64 := copyBits(uint32(val), 8, oldval64, bitoffset)
	return writeUint64(state, palign, newval64)
}

func operandRd(insn uint32) uint8 {
	return uint8(uint32ShiftRight(uint32ShiftLeft(insn, 20), 27))
}

func operandRs1(insn uint32) uint8 {
	return uint8(uint32ShiftRight(uint32ShiftLeft(insn, 12), 27))
}

func operandRs2(insn uint32) uint8 {
	return uint8(uint32ShiftRight(uint32ShiftLeft(insn, 7), 27))
}

func operandImm12(insn uint32) int32 {
	return int32ShiftRight(int32(insn), 20)
}

func operandImm20(insn uint32) int32 {
	return int32(uint32ShiftLeft(uint32ShiftRight(insn, 12), 12))
}

func operandJimm20(insn uint32) int32 {
	a := int32(uint32ShiftLeft(uint32(int32ShiftRight(int32(insn), 31)), 20))
	b := uint32ShiftLeft(uint32ShiftRight(uint32ShiftLeft(insn, 1), 22), 1)
	c := uint32ShiftLeft(uint32ShiftRight(uint32ShiftLeft(insn, 11), 31), 11)
	d := uint32ShiftLeft(uint32ShiftRight(uint32ShiftLeft(insn, 12), 24), 12)
	return int32(uint32(a) | b | c | d)
}

func operandShamt5(insn uint32) int32 {
	return int32(uint32ShiftRight(uint32ShiftLeft(insn, 7), 27))
}

func operandShamt6(insn uint32) int32 {
	return int32(uint32ShiftRight(uint32ShiftLeft(insn, 6), 26))
}

func operandSbimm12(insn uint32) int32 {
	a := int32(uint32ShiftLeft(uint32(int32ShiftRight(int32(insn), 31)), 12))
	b := uint32ShiftLeft(uint32ShiftRight(uint32ShiftLeft(insn, 1), 26), 5)
	c := uint32ShiftLeft(uint32ShiftRight(uint32ShiftLeft(insn, 20), 28), 1)
	d := uint32ShiftLeft(uint32ShiftRight(uint32ShiftLeft(insn, 24), 31), 11)
	return int32(uint32(a) | b | c | d)
}

func operandSimm12(insn uint32) int32 {
	return int32(uint32ShiftLeft(uint32(int32ShiftRight(int32(insn), 25)), 5) |
		uint32ShiftRight(uint32ShiftLeft(insn, 20), 27))
}

func advancePc(state State, pc uint64) {
	newPc := uint64AddUint64(pc, 4)
	state.writePc(newPc)
}

func branch(state State, pc uint64) {
	state.writePc(pc)
}

func executeLUI(state State, insn uint32, pc uint64) {
	rd := operandRd(insn)
	imm := operandImm20(insn)
	if rd != 0 {
		state.writeX(rd, int32ToUint64(imm))
	}
	advancePc(state, pc)
}

func executeAUIPC(state State, insn uint32, pc uint64) {
	imm := operandImm20(insn)
	rd := operandRd(insn)
	if rd != 0 {
		state.writeX(rd, uint64AddInt32(pc, imm))
	}
	advancePc(state, pc)
}

func executeJAL(state State, insn uint32, pc uint64) {
	imm := operandJimm20(insn)
	rd := operandRd(insn)
	if rd != 0 {
		state.writeX(rd, uint64AddUint64(pc, 4))
	}
	branch(state, uint64AddInt32(pc, imm))
}

func executeJALR(state State, insn uint32, pc uint64) {
	imm := operandImm12(insn)
	rd := operandRd(insn)
	rs1 := operandRs1(insn)
	rs1val := state.readX(rs1)
	if rd != 0 {
		state.writeX(rd, uint64AddUint64(pc, 4))
	}
	branch(state, uint64AddInt32(rs1val, imm)&(^uint64(1)))
}

func executeBEQ(state State, insn uint32, pc uint64) {
	imm := operandSbimm12(insn)
	rs1 := operandRs1(insn)
	rs2 := operandRs2(insn)
	rs1val := state.readX(rs1)
	rs2val := state.readX(rs2)
	if rs1val == rs2val {
		branch(state, uint64AddInt32(pc, imm))
	} else {
		advancePc(state, pc)
	}
}

func executeBNE(state State, insn uint32, pc uint64) {
	imm := operandSbimm12(insn)
	rs1 := operandRs1(insn)
	rs2 := operandRs2(insn)
	rs1val := state.readX(rs1)
	rs2val := state.readX(rs2)
	if rs1val != rs2val {
		branch(state, uint64AddInt32(pc, imm))
	} else {
		advancePc(state, pc)
	}
}

func executeBLT(state State, insn uint32, pc uint64) {
	imm := operandSbimm12(insn)
	rs1 := operandRs1(insn)
	rs2 := operandRs2(insn)
	rs1val := int64(state.readX(rs1))
	rs2val := int64(state.readX(rs2))
	if rs1val < rs2val {
		branch(state, uint64AddInt32(pc, imm))
	} else {
		advancePc(state, pc)
	}
}

func executeBGE(state State, insn uint32, pc uint64) {
	imm := operandSbimm12(insn)
	rs1 := operandRs1(insn)
	rs2 := operandRs2(insn)
	rs1val := int64(uarch_compat.ReadX(state, rs1))
	rs2val := int64(uarch_compat.ReadX(state, rs2))
	if rs1val >= rs2val {
		branch(state, uint64AddInt32(pc, imm))
	} else {
		advancePc(state, pc)
	}
}

func executeBLTU(state State, insn uint32, pc uint64) {
	imm := operandSbimm12(insn)
	rs1 := operandRs1(insn)
	rs2 := operandRs2(insn)
	rs1val := state.readX(rs1)
	rs2val := state.readX(rs2)
	if rs1val < rs2val {
		branch(state, uint64AddInt32(pc, imm))
	} else {
		advancePc(state, pc)
	}
}

func executeBGEU(state State, insn uint32, pc uint64) {
	imm := operandSbimm12(insn)
	rs1 := operandRs1(insn)
	rs2 := operandRs2(insn)
	rs1val := state.readX(rs1)
	rs2val := state.readX(rs2)
	if rs1val >= rs2val {
		branch(state, uint64AddInt32(pc, imm))
	} else {
		advancePc(state, pc)
	}
}

func executeLB(state State, insn uint32, pc uint64) {
	imm := operandImm12(insn)
	rd := operandRd(insn)
	rs1 := operandRs1(insn)
	rs1val := state.readX(rs1)
	i8 := int8(readUint8(state, uint64AddInt32(rs1val, imm)))
	if rd != 0 {
		state.writeX(rd, int8ToUint64(i8))
	}
	advancePc(state, pc)
}

func executeLHU(state State, insn uint32, pc uint64) {
	imm := operandImm12(insn)
	rd := operandRd(insn)
	rs1 := operandRs1(insn)
	rs1val := state.readX(rs1)
	u16 := readUint16(uint64AddInt32(rs1val, imm))
	if rd != 0 {
		state.writeX(rd, u16)
	}
	advancePc(state, pc)
}

func executeLH(state State, insn uint32, pc uint64) {
	imm := operandImm12(insn)
	rd := operandRd(insn)
	rs1 := operandRs1(insn)
	rs1val := state.readX(rs1)
	i16 := int16(readUint16(state, uint64AddInt32(rs1val, imm)))
	if rd != 0 {
		state.writeX(rd, int16ToUint64(i16))
	}
	advancePc(state, pc)
}

func executeLW(state State, insn uint32, pc uint64) {
	imm := operandImm12(insn)
	rd := operandRd(insn)
	rs1 := operandRs1(insn)
	rs1val := state.readX(rs1)
	i32 := int32(readUint32(state, uint64AddInt32(rs1val, imm)))
	if rd != 0 {
		state.writeX(rd, int32ToUint64(i32))
	}
	advancePc(state, pc)
}

func executeLBU(state State, insn uint32, pc uint64) {
	imm := operandImm12(insn)
	rd := operandRd(insn)
	rs1 := operandRs1(insn)
	rs1val := state.readX(rs1)
	u8 := readUint8(state, uint64AddInt32(rs1val, imm))
	if rd != 0 {
		state.writeX(rd, u8)
	}
	advancePc(state, pc)
}
