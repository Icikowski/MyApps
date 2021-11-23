package types

import "time"

// Deployment represents details about a deployment.
type Deployment struct {
	Repository  string    `json:"repository" yaml:"repository"`
	Application string    `json:"application" yaml:"application"`
	InstalledOn time.Time `json:"installed_on" yaml:"installed_on"`
	UpdatedOn   time.Time `json:"updated_on" yaml:"updated_on"`
}

// String returns the string representation of a deployment.
func (d Deployment) String() string {
	return d.Repository + "/" + d.Application
}

// Deployments is a list of deployments.
type Deployments []Deployment

// Exists check if a deployment exists.
func (d Deployments) Exists(repository, application string) bool {
	for _, deployment := range d {
		if deployment.Repository == repository && deployment.Application == application {
			return true
		}
	}
	return false
}

// Find returns the deployment with the given repository and application.
func (d Deployments) Find(repository, application string) (Deployment, bool) {
	for _, deployment := range d {
		if deployment.Repository == repository && deployment.Application == application {
			return deployment, true
		}
	}
	return Deployment{}, false
}

// Add adds a deployment to the list.
func (d Deployments) Add(deployment Deployment) Deployments {
	return append(d, deployment)
}

// Update updates a deployment on the list.
func (d Deployments) Update(deployment Deployment) Deployments {
	for i, depl := range d {
		if depl.Repository == deployment.Repository && depl.Application == deployment.Application {
			d[i] = deployment
			return d
		}
	}
	return d
}

// Delete deletes a deployment from the list.
func (d Deployments) Delete(repository, application string) Deployments {
	for i, deployment := range d {
		if deployment.Repository == repository && deployment.Application == application {
			d = append(d[:i], d[i+1:]...)
			return d
		}
	}
	return d
}
