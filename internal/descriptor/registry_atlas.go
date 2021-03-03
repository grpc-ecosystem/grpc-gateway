package descriptor

// IsAtlasPatch whether generation is followed by atlas-patch changes.
func (r *Registry) IsAtlasPatch() bool {
	return r.atlasPatch
}

func (r *Registry) SetAtlasPatch(atlas bool) {
	r.atlasPatch = atlas
}

// IsWithPrivateOperations if true exclude all operations with tag "private"
func (r *Registry) IsWithPrivateOperations() bool {
	return r.withPrivateOperations
}

func (r *Registry) SetWithPrivateOperations(withPrivateOperations bool) {
	r.withPrivateOperations = withPrivateOperations
}

func (r *Registry) SetWithCustomAnnotations(custom bool) {
	r.withCustomAnnotations = custom
}

func (r *Registry) IsWithCustomAnnotations() bool {
	return r.withCustomAnnotations
}
