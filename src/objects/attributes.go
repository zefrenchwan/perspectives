package objects

import "github.com/zefrenchwan/perspectives.git/periods"

// AttributeDetails represents the metadata details of the attribute.
// It contains information about the attribute's name, type, validity, and instance activity.
type AttributeDetails struct {
	// AttributeName is the actual name of the attribute
	AttributeName string
	// AttributeType is the actual type of the attribute
	AttributeType string
	// AttributeValidity is the validity period of the attribute
	AttributeValidity periods.Period
	// InstanceActivity is the activity period of the instance
	InstanceActivity periods.Period
}
