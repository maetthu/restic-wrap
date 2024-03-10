package profile

import "fmt"

type Profile struct {
	Env      []EnvVar
	Backends []*Backend
	Stages   []Stage
	Notify   []string
}

func (p *Profile) Backend(name string) (*Backend, error) {
	for _, b := range p.Backends {
		if b.Name == name {
			return b, nil
		}
	}

	return nil, fmt.Errorf("Backend %s not found in profile", name)
}

func (p *Profile) BuildEnv(backend *Backend) []string {
	env := []string{}

	for _, e := range p.Env {
		env = append(env, fmt.Sprintf("%s=%s", e.Name, e.Value))
	}

	for _, e := range backend.Env {
		env = append(env, fmt.Sprintf("%s=%s", e.Name, e.Value))
	}

	env = append(env, fmt.Sprintf("RESTIC_REPOSITORY=%s", backend.Repository))
	env = append(env, fmt.Sprintf("RESTIC_PASSWORD=%s", backend.Password))

	return env
}

type EnvVar struct {
	Name  string
	Value string
}

type Backend struct {
	Name       string
	Repository string
	Password   string
	Env        []EnvVar
}

type Stage struct {
	Command string
	Args    []string
}
