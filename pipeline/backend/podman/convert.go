package podman

import (
	"net"

	"github.com/containers/podman/v3/pkg/specgen"
	"github.com/rs/zerolog/log"

	backend "github.com/woodpecker-ci/woodpecker/pipeline/backend/types"
)

// returns specgen.SpecGenerator container config
func toSpecGenerator(proc *backend.Step) (*specgen.SpecGenerator, error) {
	var err error
	specGen := specgen.NewSpecGenerator(proc.Image, false)
	specGen.Terminal = true

	specGen.Image = proc.Image
	specGen.Name = proc.Name
	specGen.Labels = proc.Labels
	specGen.WorkDir = proc.WorkingDir

	if len(proc.Environment) > 0 {
		specGen.Env = proc.Environment
	}
	if len(proc.Command) > 0 {
		specGen.Command = proc.Command
	}

	log.Trace().Msgf("proc.Entrypoint: %v", proc.Entrypoint)

	if len(proc.Entrypoint) > 0 {
		specGen.Entrypoint = proc.Entrypoint
	}
	log.Trace().Msgf("proc.volumes: %v", proc.Volumes)
	if len(proc.Volumes) > 0 {
		for _, v := range proc.Volumes {
			log.Trace().Msgf("proc.vol: %v", v)
		}
		_, vols, _, err := specgen.GenVolumeMounts(proc.Volumes)
		if err != nil {
			return nil, err
		}
		for _, vol := range vols {
			log.Trace().Msgf("specgen.volume: %v", vol)
			specGen.Volumes = append(specGen.Volumes, vol)
		}
	}

	specGen.LogConfiguration = &specgen.LogConfig{
		//Driver: "json-file",
	}
	// TODO: Privileged seems to be required
	specGen.Privileged = true
	specGen.ShmSize = new(int64)
	*specGen.ShmSize = proc.ShmSize
	specGen.Sysctl = proc.Sysctls

	if len(proc.IpcMode) > 0 {
		if specGen.IpcNS, err = specgen.ParseNamespace(proc.IpcMode); err != nil {
			return nil, err
		}
	}
	if len(proc.DNS) > 0 {
		for _, dns := range proc.DNS {
			if ip := net.ParseIP(dns); ip != nil {
				specGen.DNSServers = append(specGen.DNSServers, ip)
			}
		}
	}
	if len(proc.DNSSearch) > 0 {
		specGen.DNSSearch = proc.DNSSearch
	}
	if len(proc.ExtraHosts) > 0 {
		specGen.HostAdd = proc.ExtraHosts
	}

	return specGen, err
}
