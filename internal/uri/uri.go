package uri

import (
	"errors"
	"fmt"
)

type URI struct {
	Raw     string
	Entity  EntityType
	Context *string
	Domain  *string
	Slug    string
}

func Parse(raw string) (*URI, error) {
	for _, p := range uriPatterns {
		m := p.re.FindStringSubmatch(raw)
		if m == nil {
			continue
		}

		var context, domain *string

		if p.hasContext {
			context = &m[1]

			if p.hasDomain {
				domain = &m[2]
			}
		}

		slug := m[len(m)-1]

		uri := URI{
			Raw:     raw,
			Entity:  p.entityType,
			Context: context,
			Domain:  domain,
			Slug:    slug,
		}

		return &uri, nil
	}

	return nil, fmt.Errorf("%q does not match any known scio:// URI pattern", raw)
}

func (uri *URI) String() string {
	return uri.Raw
}

func (uri *URI) ParentURI() (string, error) {
	if uri.Context == nil {
		return "", errors.New("entity doesn't have a parent")
	}

	if uri.Entity == EntityTypeConcept {
		return fmt.Sprintf("scio://contexts/%s/domains/%s", *uri.Context, *uri.Domain), nil
	}

	return fmt.Sprintf("scio://contexts/%s", *uri.Context), nil
}
