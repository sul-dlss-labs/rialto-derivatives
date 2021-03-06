package repository

import (
	"fmt"

	"github.com/knakk/rdf"
	"github.com/knakk/sparql"
	"github.com/sul-dlss/rialto-derivatives/models"
)

// Repository is an interface that rialto-derivatives reads from as its source
type Repository interface {
	SubjectsToResources(subjects []string) ([]models.Resource, error)
	AllResources(func([]models.Resource) error) error
}

// Service is the Neptune implementation of the repository
type Service struct {
	reader Reader
}

// NewService creates a new Service instance
func NewService(reader Reader) Repository {
	return &Service{reader: reader}
}

// SubjectsToResources takes a subject string and returns a resource
func (m *Service) SubjectsToResources(subjects []string) ([]models.Resource, error) {
	responseSet, err := m.reader.QueryByIDs(subjects)
	results := []models.Resource{}
	if err != nil {
		return nil, err
	}
	for n, response := range responseSet {
		list := m.toResourceList(response.Solutions())
		if len(list) == 0 {
			return nil, fmt.Errorf("Record not found: %s", subjects[n])
		}
		results = append(results, list[0])
	}

	return results, nil
}

// AllResources returns a full list of resources
func (m *Service) AllResources(f func([]models.Resource) error) error {
	err := m.reader.QueryEverything(func(response *sparql.Results) error {
		// Solutions look like this:
		// [map[o:AA00 s:http://rialto.stanford.edu/stanford p:http://rialto.stanford.edu/vocab/organizationCodes]]
		// log.Printf("Solutions %s", response.Solutions())
		list := m.toResourceList(response.Solutions())
		return f(list)
	})
	return err
}

func (m *Service) toResourceList(solutions []map[string]rdf.Term) []models.Resource {
	list := []models.Resource{}
	for _, solution := range solutions {
		resource := models.NewResource(solution)
		if v, ok := resource.(*models.Person); ok {
			m.toPerson(v)
		} else if v, ok := resource.(*models.Publication); ok {
			m.toPublication(v)
		} else if v, ok := resource.(*models.Grant); ok {
			m.toGrant(v)
		}
		list = append(list, resource)
	}
	return list
}

func (m *Service) toPublication(v *models.Publication) {
	// Publications need to be informed of their authors.
	results, err := m.reader.GetAuthorInfo(v.Subject())
	if err != nil {
		panic(err)
	}
	v.SetAuthorInfo(results)
	// Publications need to be informed of their concepts.
	results, err = m.reader.GetConceptInfo(v.Subject())
	if err != nil {
		panic(err)
	}
	v.SetConceptInfo(results)

	// Publications need to be informed of their grants.
	results, err = m.reader.GetGrantInfo(v.Subject())
	if err != nil {
		panic(err)
	}
	v.SetGrantInfo(results)

	// Publications need to be informed of their identifiers.
	results, err = m.reader.GetIdentifierInfo(v.Subject())
	if err != nil {
		panic(err)
	}
	v.SetIdentifierInfo(results)
}

func (m *Service) toGrant(v *models.Grant) {
	// Grants need to be informed of their identifiers.
	results, err := m.reader.GetGrantIdentifierInfo(v.Subject())
	if err != nil {
		panic(err)
	}
	v.SetIdentifierInfo(results)
}

func (m *Service) toPerson(v *models.Person) {
	// People also need to be informed of their organization membership.
	results, err := m.reader.GetPositionOrganizationInfo(v.Subject())
	if err != nil {
		panic(err)
	}
	v.SetPositionOrganizationInfo(results)
	// People also need to be informed of their countries.
	results, err = m.reader.GetCountriesInfo(v.Subject())
	if err != nil {
		panic(err)
	}
	v.SetCountriesInfo(results)
	// People also need to be informed of their subtypes.
	results, err = m.reader.GetPersonSubtypesInfo(v.Subject())
	if err != nil {
		panic(err)
	}
	v.SetPersonSubtypesInfo(results)
}
