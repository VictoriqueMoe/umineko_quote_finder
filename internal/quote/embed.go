package quote

import "embed"

//go:embed data/*.file data/sub/*.ass data/higurashi/*.file
var DataFS embed.FS
