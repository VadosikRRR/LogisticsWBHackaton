package entity

type OfficeCapacityPolicy struct {
	VehicleCapacity      float64
	SafetyBufferFraction float64
	MaxVehiclesPerSlot   int
}

type CapacityPolicySet struct {
	DefaultPolicy OfficeCapacityPolicy
	ByOffice      map[int64]OfficeCapacityPolicy
}

func (p CapacityPolicySet) ForOffice(officeID int64) OfficeCapacityPolicy {
	if p.ByOffice == nil {
		return p.DefaultPolicy
	}
	if policy, ok := p.ByOffice[officeID]; ok {
		return policy
	}
	return p.DefaultPolicy
}
