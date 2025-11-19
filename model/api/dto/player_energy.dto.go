package dto

type PlayerEnergy struct {
	CurrentEnergy int32 `json:"current_energy"`
	MaxEnergy     int32 `json:"max_energy"`
	NextEnergyAt  int64 `json:"next_energy_at"`
}
