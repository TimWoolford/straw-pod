package status

import (
	"os"
	"strconv"
)

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
	namespace          string
	port               int
}

func MakeMe() Status {
	port, _ := strconv.Atoi(os.Getenv("STATUS_PORT"))
	return Status{
		ApplicationVersion: "1.0",
		Hostname:           os.Getenv("HOSTNAME"),
		OverallStatus:      "OK",
		Probes:				make([]Probe, 0),
		namespace:          os.Getenv("POD_NAMESPACE"),
		port:               port,
	}
}

func (s *Status) Port() int {
	return s.port
}

func (s *Status) Namespace() string {
	return s.namespace
}
