package cartesi_machine_step

type IUArchStep interface {
	Step() (IUArchState, bool)
}
