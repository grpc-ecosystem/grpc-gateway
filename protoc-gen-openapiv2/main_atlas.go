package main

import (
	"flag"

	"github.com/grpc-ecosystem/grpc-gateway/v2/internal/descriptor"
)

var (
	atlasPatch            = flag.Bool("atlas_patch", false, "if set, generation will be proceded with atlas-patch changes")
	withPrivate           = flag.Bool("with_private", false, "if unset, generate swagger schema without operations 0 as 'private' work only if atlas_patch set")
	withCustomAnnotations = flag.Bool("with_custom_annotations", false, "if set, you became available to use custom annotations")
)

func atlasFlags(reg *descriptor.Registry) {
	reg.SetAtlasPatch(*atlasPatch)
	reg.SetWithPrivateOperations(*withPrivate)
	reg.SetWithCustomAnnotations(*withCustomAnnotations)
}
