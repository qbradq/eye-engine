package formats

import (
	"bufio"
	"errors"
	"io"
	"log/slog"
	"strconv"
	"strings"

	"github.com/go-gl/mathgl/mgl32"
	"github.com/qbradq/eye-engine/engine/types"
)

// objFace defines a single face in an OBJ model.
type objFace struct {
	v []int // Vertex indexes
	u []int // UV indexes
	n []int // Normal indexes
}

// ParseVert parses a single vertex definition in OBJ format and adds that
// vertex to the face.
func (f *objFace) ParseVert(s string) {
	fn := func(s string) int {
		v, err := strconv.ParseInt(s, 10, 32)
		if err != nil {
			slog.Error("error parsing OBJ face vertex", "error", err)
			return 0
		}
		return int(v)
	}
	parts := strings.Split(s, "/")
	vi := int(0)
	ui := int(0)
	ni := int(0)
	if len(parts) > 0 {
		vi = fn(parts[0])
	}
	if len(parts) > 1 {
		ui = fn(parts[1])
	}
	if len(parts) > 2 {
		ni = fn(parts[2])
	}
	f.v = append(f.v, vi)
	f.u = append(f.u, ui)
	f.n = append(f.n, ni)
}

// objModel represents a 3D model in OBJ format.
type objModel struct {
	v []mgl32.Vec3 // Vertex positions
	n []mgl32.Vec3 // Vertex normals
	u []mgl32.Vec2 // Vertex UVs
	f []objFace    // Faces
}

// ToModel returns a new engine.types.Model from this model.
func (o *objModel) ToModel() *types.Model {
	idx := func(i, l int) int {
		if i == 0 {
			return 0
		}
		if i > 0 {
			return i - 1
		}
		return l + i
	}
	ret := &types.Model{}
	for i := range o.f {
		// Fan triangulation
		f := o.f[i]
		ia := 0
		for j := 2; j < len(f.v); j++ {
			ib := j - 1
			ic := j - 0
			ret.Faces = append(ret.Faces, types.Face{
				Vertexes: [3]types.Vertex{
					{
						Position: o.v[idx(f.v[ia], len(o.v))],
						Normal:   o.n[idx(f.n[ia], len(o.n))],
						UV:       o.u[idx(f.u[ia], len(o.u))],
					},
					{
						Position: o.v[idx(f.v[ib], len(o.v))],
						Normal:   o.n[idx(f.n[ib], len(o.n))],
						UV:       o.u[idx(f.u[ib], len(o.u))],
					},
					{
						Position: o.v[idx(f.v[ic], len(o.v))],
						Normal:   o.n[idx(f.n[ic], len(o.n))],
						UV:       o.u[idx(f.u[ic], len(o.u))],
					},
				},
			})
		}
	}
	return ret
}

// LoadModelFromOBJ loads a new Model object from a .obj file in r.
func LoadModelFromOBJ(r io.Reader) *types.Model {
	// Load all data from the file
	fn := func(s string) float32 {
		v, err := strconv.ParseFloat(s, 32)
		if err != nil {
			slog.Error("error parsing OBJ model", "error", err)
			return 0
		}
		return float32(v)
	}
	scanner := bufio.NewScanner(r)
	verts := []mgl32.Vec3{}
	normals := []mgl32.Vec3{}
	uvs := []mgl32.Vec2{}
	faces := []objFace{}
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if len(line) == 0 || line[0] == '#' {
			// Skip empty and comment lines
			continue
		}
		parts := strings.Fields(line)
		switch parts[0] {
		case "v":
			if len(parts) < 4 || len(parts) > 5 {
				slog.Error(
					"error parsing OBJ model",
					"error",
					errors.New("expected 3 or 4 arguments to v line"),
				)
				return nil
			}
			verts = append(verts,
				mgl32.Vec3{
					fn(parts[1]),
					fn(parts[2]),
					fn(parts[3]),
				},
			)
			if len(parts) == 5 {
				w := fn(parts[4])
				if w != 1.0 {
					slog.Error(
						"error parsing OBJ model",
						"error",
						errors.New("the only supported w value is 1.0"),
					)
					return nil
				}
			}
		case "vn":
			if len(parts) != 4 {
				slog.Error(
					"error parsing OBJ model",
					"error",
					errors.New("expected 3 arguments to vn line"),
				)
				return nil
			}
			normals = append(normals,
				mgl32.Vec3{
					fn(parts[1]),
					fn(parts[2]),
					fn(parts[3]),
				},
			)
		case "vt":
			if len(parts) < 2 || len(parts) > 4 {
				slog.Error(
					"error parsing OBJ model",
					"error",
					errors.New("expected 1 to 3 arguments to vt line"),
				)
				return nil
			}
			uv := mgl32.Vec2{fn(parts[1])}
			if len(parts) > 2 {
				uv[1] = fn(parts[1])
			}
			uvs = append(uvs, uv)
		case "f":
			if len(parts) < 4 {
				slog.Error(
					"error parsing OBJ model",
					"error",
					errors.New("expected at least 3 vertexes in f line"),
				)
				return nil
			}
			f := objFace{}
			for _, p := range parts[1:] {
				f.ParseVert(p)
			}
			faces = append(faces, f)
		}
	}
	if err := scanner.Err(); err != nil {
		slog.Error("error scanning OBJ model", "error", err)
		return nil
	}
	// Back-patch negative indexes
	for iFace := range faces {
		for iVert := range faces[iFace].v {
			i := faces[iFace].v[iVert]
			if i < 0 {
				faces[iFace].v[iVert] = len(verts) + i
			}
		}
		for iUV := range faces[iFace].u {
			i := faces[iFace].u[iUV]
			if i < 0 {
				faces[iFace].u[iUV] = len(uvs) + i
			}
		}
		for iNormal := range faces[iFace].n {
			i := faces[iFace].n[iNormal]
			if i < 0 {
				faces[iFace].n[iNormal] = len(normals) + i
			}
		}
	}
	// Construct and return the new model
	obj := NewModelFromData(verts, normals, uvs, faces)
	if obj == nil {
		return nil
	}
	ret := obj.ToModel()
	ret.UpdateGeometry()
	return ret
}

// NewModelFromData returns a new model constructed from data arrays.
func NewModelFromData(verts, normals []mgl32.Vec3, uvs []mgl32.Vec2, faces []objFace) *objModel {
	return &objModel{
		v: verts,
		n: normals,
		u: uvs,
		f: faces,
	}
}
