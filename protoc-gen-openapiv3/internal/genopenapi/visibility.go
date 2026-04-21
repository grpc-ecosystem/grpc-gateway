package genopenapi

import (
	"strings"

	"github.com/grpc-ecosystem/grpc-gateway/v2/internal/descriptor"
	"google.golang.org/genproto/googleapis/api/visibility"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/descriptorpb"
)

// isVisible checks whether an element with the given VisibilityRule should be
// included in the generated output. Elements without an annotation are always
// visible. When an annotation is present, at least one of its comma-separated
// restriction labels must appear in the registry's configured selectors.
func isVisible(r *visibility.VisibilityRule, reg *descriptor.Registry) bool {
	if r == nil {
		return true
	}

	selectors := reg.GetVisibilityRestrictionSelectors()

	restrictions := strings.Split(strings.TrimSpace(r.Restriction), ",")
	if len(restrictions) == 0 {
		return true
	}

	for _, restriction := range restrictions {
		if selectors[strings.TrimSpace(restriction)] {
			return true
		}
	}

	return false
}

func fieldVisibility(fd *descriptor.Field) *visibility.VisibilityRule {
	if fd.Options == nil {
		return nil
	}
	if !proto.HasExtension(fd.Options, visibility.E_FieldVisibility) {
		return nil
	}
	opts, ok := proto.GetExtension(fd.Options, visibility.E_FieldVisibility).(*visibility.VisibilityRule)
	if !ok {
		return nil
	}
	return opts
}

func serviceVisibility(svc *descriptor.Service) *visibility.VisibilityRule {
	if svc.Options == nil {
		return nil
	}
	if !proto.HasExtension(svc.Options, visibility.E_ApiVisibility) {
		return nil
	}
	opts, ok := proto.GetExtension(svc.Options, visibility.E_ApiVisibility).(*visibility.VisibilityRule)
	if !ok {
		return nil
	}
	return opts
}

func methodVisibility(m *descriptor.Method) *visibility.VisibilityRule {
	if m.Options == nil {
		return nil
	}
	if !proto.HasExtension(m.Options, visibility.E_MethodVisibility) {
		return nil
	}
	opts, ok := proto.GetExtension(m.Options, visibility.E_MethodVisibility).(*visibility.VisibilityRule)
	if !ok {
		return nil
	}
	return opts
}

func enumValueVisibility(v *descriptorpb.EnumValueDescriptorProto) *visibility.VisibilityRule {
	if v.Options == nil {
		return nil
	}
	if !proto.HasExtension(v.Options, visibility.E_ValueVisibility) {
		return nil
	}
	opts, ok := proto.GetExtension(v.Options, visibility.E_ValueVisibility).(*visibility.VisibilityRule)
	if !ok {
		return nil
	}
	return opts
}
