package coruntime

import _ "embed"

//go:embed c/coco_runtime.c
var Source []byte
