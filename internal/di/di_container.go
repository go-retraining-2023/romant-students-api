package di

import (
	"errors"
	"sync"
)

type Container struct {
	dependencies map[string]interface{}
	mu           sync.Mutex
}

var applicationContainer *Container

func GetAppContainer() *Container {
	if applicationContainer == nil {
		applicationContainer = NewContainer()
	}

	return applicationContainer
}

func NewContainer() *Container {
	return &Container{
		dependencies: make(map[string]interface{}),
	}
}

// Register adds a dependency to the container.
func (c *Container) Register(name string, dependency interface{}) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.dependencies[name] = dependency
}

// Resolve retrieves a dependency from the container.
func (c *Container) Resolve(name string) (interface{}, error) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if dep, ok := c.dependencies[name]; ok {
		return dep, nil
	}

	return nil, errors.New("dependency not found")
}
