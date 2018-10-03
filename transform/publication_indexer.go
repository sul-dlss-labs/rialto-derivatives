package transform

import (
	"github.com/sul-dlss-labs/rialto-derivatives/models"
	"github.com/vanng822/go-solr/solr"
)

// PublicationIndexer transforms publication resource types into solr Documents
type PublicationIndexer struct {
}

// Index adds fields from the resource to the Solr Document
func (m *PublicationIndexer) Index(resource *models.Publication, doc solr.Document) solr.Document {
	doc.Set("type_ssi", "Publication")
	doc.Set("title_tesi", resource.Title)
	doc.Set("created_ssim", resource.Created)
	doc.Set("identifier_ssim", resource.Identifier)

	if resource.DOI != nil {
		doc.Set("doi_ssim", *resource.DOI)
	}

	if resource.Abstract != nil {
		doc.Set("abstract_tesim", *resource.Abstract)
	}

	var authors = []string{}
	var authorLabels = []string{}
	for _, author := range resource.Authors {
		authors = append(authors, author.URI)
		authorLabels = append(authorLabels, author.Label)
	}

	doc.Set("authors_ssim", authors)
	doc.Set("author_labels_tsim", authorLabels)

	if resource.Description != nil {
		doc.Set("description_tesim", *resource.Description)
	}

	if resource.Publisher != nil {
		doc.Set("publisher_ssim", *resource.Publisher)
	}

	// TODO Fields still to map:
	// "cites":            "cites_ssim",
	// "link":             "link_ssim",
	// "fundedBy":         "funded_by_ssim",
	// "sponsor":          "sponsor_label_tsim",   // TODO: Needs URI lookup
	// "hasInstrument":    "has_instrument_ssim",
	// "sameAs":           "same_as_ssim",
	// "journalIssue":     "journal_issue_ssim",
	// "subject":          "subject_label_ssim", // TODO: Needs URI
	// "alternativeTitle": "alternative_title_tesim",

	// TODO: complex lookups
	// Profiles confirmed 	vivo:relatedBy vivo:Authorship dcterms:source 	"Profiles" string-literal 	[0,1] 	If the authorship relationship has been confirmed by the Author in Profiles. Can be reused for any relationship needed (i.e. Editorship, Advising Relationship, etc.)
	// editor 	vivo:relatedBy vivo:Editorship vivo:relates 	URI for foaf:Agent 	[0,n] 	Editor of the publication.

	return doc
}
