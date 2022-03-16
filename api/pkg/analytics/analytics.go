package analytics

import "getsturdy.com/api/pkg/users"

type CaptureOptions struct {
	DistinctId string
	Groups     map[string]any
	Properties map[string]any
}

type CaptureOption func(*CaptureOptions)

func Property(key string, value any) CaptureOption {
	return func(o *CaptureOptions) {
		if o.Properties == nil {
			o.Properties = make(map[string]any)
		}
		o.Properties[key] = value
	}
}

func UserID(id users.ID) CaptureOption {
	return func(o *CaptureOptions) {
		o.DistinctId = id.String()
	}
}

func DistinctID(id string) CaptureOption {
	return func(o *CaptureOptions) {
		o.DistinctId = id
	}
}

func CodebaseID(id string) CaptureOption {
	return func(o *CaptureOptions) {
		if o.Properties == nil {
			o.Properties = make(map[string]any)
		}
		o.Properties["codebase_id"] = id
		if o.Groups == nil {
			o.Groups = make(map[string]any)
		}
		o.Groups["codebase"] = id
	}
}

func OrganizationID(id string) CaptureOption {
	return func(o *CaptureOptions) {
		if o.Properties == nil {
			o.Properties = make(map[string]any)
		}
		o.Properties["organization_id"] = id
		if o.Groups == nil {
			o.Groups = make(map[string]any)
		}
		o.Groups["organization"] = id
	}
}
