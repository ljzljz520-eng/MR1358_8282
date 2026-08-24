package policy

type Allocation struct {
	Capacity  int
	Reserved  int
	Requested int
}

func (a Allocation) Remaining() int {
	remaining := a.Capacity - a.Reserved
	if remaining < 0 {
		return 0
	}
	return remaining
}

func (a Allocation) CanReserve() bool { return a.Requested > 0 && a.Requested <= a.Remaining() }

func (a Allocation) AfterReserve() Allocation {
	if !a.CanReserve() {
		return a
	}
	result := a
	result.Reserved += a.Requested
	result.Requested = 0
	return result
}

func (a Allocation) IsFull() bool { return a.Remaining() == 0 }

func (a Allocation) Summary() string {
	if a.IsFull() {
		return "full"
	}
	return "open"
}

func SplitAllocation(capacity, reserved, requested int) (Allocation, Allocation) {
	allocation := Allocation{Capacity: capacity, Reserved: reserved, Requested: requested}
	return allocation, allocation.AfterReserve()
}

func ReserveMany(capacity int, requests []int) (int, []int) {
	reserved := 0
	accepted := make([]int, 0, len(requests))
	for _, request := range requests {
		if request > 0 && reserved+request <= capacity {
			reserved += request
			accepted = append(accepted, request)
		}
	}
	return reserved, accepted
}
