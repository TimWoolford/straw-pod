package status


type Probe struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Status      string `json:"status"`
}

type Status struct {
	ApplicationVersion string  `json:"applicationVersion"`
	Hostname           string  `json:"hostname"`
	OverallStatus      string  `json:"overallStatus"`
	Probes             []Probe `json:"probes"`
	Namespace          string
	Port               int
}

