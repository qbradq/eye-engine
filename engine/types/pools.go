package types

import (
	"github.com/qbradq/eye-engine/internal/util"
)

// PointI2DPool is the global memory pool for PointI2D slices.
var PointI2DPool *util.SlicePool[PointI2D] = util.NewSlicePool[PointI2D](512)
