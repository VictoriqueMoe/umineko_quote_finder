package quote

import "embed"

//go:embed data/*.file data/sub/*.ass data/higurashi/*.file data/Ciconia/*.file
var DataFS embed.FS
