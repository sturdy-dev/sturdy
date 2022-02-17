package analytics

type CaptureOptions struct {
	DistinctId string
	Groups     map[string]interface{}
	Properties map[string]interface{}
}

type CaptureOption func(*CaptureOptions)

func Property(key string, value interface{}) CaptureOption {
	return func(o *CaptureOptions) {
		if o.Properties == nil {
			o.Properties = make(map[string]interface{})
		}
		o.Properties[key] = value
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
			o.Properties = make(map[string]interface{})
		}
		o.Properties["codebase_id"] = id
		if o.Groups == nil {
			o.Groups = make(map[string]interface{})
		}
		o.Groups["codebase"] = id
	}
}

func OrganizationID(id string) CaptureOption {
	return func(o *CaptureOptions) {
		if o.Properties == nil {
			o.Properties = make(map[string]interface{})
		}
		o.Properties["organization_id"] = id
		if o.Groups == nil {
			o.Groups = make(map[string]interface{})
		}
		o.Groups["organization"] = id
	}
}
