package response

type PlayerEnergy struct {
	CurrentEnergy int   `json:"current_energy"`
	MaxEnergy     int   `json:"max_energy"`
	NextEnergyAt  int64 `json:"next_energy_at"`
}
