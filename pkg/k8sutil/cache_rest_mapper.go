package k8sutil

import (
	"sync"

	"golang.org/x/time/rate"
	"k8s.io/apimachinery/pkg/api/meta"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/rest"
	"sigs.k8s.io/controller-runtime/pkg/client/apiutil"
	"sigs.k8s.io/controller-runtime/pkg/runtime/log"
)

// NewCachedRESTMapper .
func NewCachedRESTMapper(config *rest.Config) (meta.RESTMapper, error) {
	c := &Cached{
		limiter: rate.NewLimiter(rate.Limit(1), 2),
		factory: func() (meta.RESTMapper, error) {
			return apiutil.NewDiscoveryRESTMapper(config)
		},
	}
	c.flush()
	return c, nil
}

// Cached .
type Cached struct {
	sync.Mutex

	limiter *rate.Limiter
	factory func() (meta.RESTMapper, error)
	mapper  meta.RESTMapper
}

func (c *Cached) flush() error {
	c.Lock()
	defer c.Unlock()

	var err error
	if c.mapper == nil || c.limiter.Allow() {
		c.mapper, err = c.factory()
	}
	return err
}

func (c *Cached) shouldFlushOn(err error) bool {
	switch err.(type) {
	case *meta.NoKindMatchError:
		return true
	}
	return false
}

func (c *Cached) onError(err error) bool {
	if !c.shouldFlushOn(err) {
		return false
	}
	if err := c.flush(); err != nil {
		log.Log.Error(err, "failed to reload RESTMapper")
		return false
	}
	return true
}

// KindFor .
func (c *Cached) KindFor(resource schema.GroupVersionResource) (schema.GroupVersionKind, error) {
	gvk, err := c.mapper.KindFor(resource)
	if c.onError(err) {
		gvk, err = c.mapper.KindFor(resource)
	}
	return gvk, err
}

// KindsFor .
func (c *Cached) KindsFor(resource schema.GroupVersionResource) ([]schema.GroupVersionKind, error) {
	gvks, err := c.mapper.KindsFor(resource)
	if c.onError(err) {
		gvks, err = c.mapper.KindsFor(resource)
	}
	return gvks, err
}

// ResourceFor .
func (c *Cached) ResourceFor(input schema.GroupVersionResource) (schema.GroupVersionResource, error) {
	gvr, err := c.mapper.ResourceFor(input)
	if c.onError(err) {
		gvr, err = c.mapper.ResourceFor(input)
	}
	return gvr, err
}

// ResourcesFor .
func (c *Cached) ResourcesFor(input schema.GroupVersionResource) ([]schema.GroupVersionResource, error) {
	gvrs, err := c.mapper.ResourcesFor(input)
	if c.onError(err) {
		gvrs, err = c.mapper.ResourcesFor(input)
	}
	return gvrs, err
}

// RESTMapping .
func (c *Cached) RESTMapping(gk schema.GroupKind, versions ...string) (*meta.RESTMapping, error) {
	m, err := c.mapper.RESTMapping(gk, versions...)
	if c.onError(err) {
		m, err = c.mapper.RESTMapping(gk, versions...)
	}
	return m, err
}

// RESTMappings .
func (c *Cached) RESTMappings(gk schema.GroupKind, versions ...string) ([]*meta.RESTMapping, error) {
	ms, err := c.mapper.RESTMappings(gk, versions...)
	if c.onError(err) {
		ms, err = c.mapper.RESTMappings(gk, versions...)
	}
	return ms, err
}

// ResourceSingularizer .
func (c *Cached) ResourceSingularizer(resource string) (singular string, err error) {
	s, err := c.mapper.ResourceSingularizer(resource)
	if c.onError(err) {
		s, err = c.mapper.ResourceSingularizer(resource)
	}
	return s, err
}
