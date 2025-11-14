package types

import (
	"github.com/qbradq/eye-engine/internal/util"
)

// pointI2DPool is the global memory pool for PointI2D slices.
var pointI2DPool *util.SlicePool[PointI2D] = util.NewSlicePool[PointI2D](512)

// polyFillEdgeDefPool is the global memory pool of polyFillEdgeDef slices.
var polyFillEdgeDefPool = util.NewSlicePool[polyFillEdgeDef](16)
