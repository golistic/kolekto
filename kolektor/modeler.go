// Copyright (c) 2022, Geert JM Vanderkelen

package kolektor

// Modeler defines methods which must be implemented by structs,
// so they are considered to be models. Note that this is actually
// done by embedding the Model-struct.
type Modeler interface {
	CollectionName() string
	SetMeta(m *Meta)
	GetID() int64
	GetUID() string
}

// Model is to be embedded by structs. It implements all methods of the
// Modeler-interface except CollectionName().
type Model struct {
	Meta *Meta `json:"Meta,omitempty"`
}

// CollectionName returns the name of the Collection in which this Model's
// objects will be stored. Note that this will be used for, for example, the
// SQL table. It is advised to use the plural, for example, 'books'.
func (m *Model) CollectionName() string {
	panic("models must implement Modeller by embedding kolektor.Model and implementing CollectionName")
}

// SetMeta stores metadata meta.
func (m *Model) SetMeta(meta *Meta) {
	m.Meta = meta
}

// GetID returns the data stores primary key, which is always an int64.
func (m *Model) GetID() int64 {
	if m.Meta == nil {
		m.Meta = &Meta{}
	}
	return m.Meta.ID
}

// GetUID returns the unique identifier of the object.
func (m *Model) GetUID() string {
	if m.Meta == nil {
		m.Meta = &Meta{}
	}
	return m.Meta.UID
}
