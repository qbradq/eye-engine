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
	ret := &types.Model{}
	for i := range o.f {
		f := types.Face{}
		for j := range o.f[i].v {
			v := o.f[i].v[j]
			n := o.f[i].n[j]
			u := o.f[i].u[j]
			var vv mgl32.Vec3
			var nv mgl32.Vec3
			var uv mgl32.Vec2
			if v > 0 {
				vv = o.v[v-1]
			}
			if n > 0 {
				nv = o.n[n-1]
			}
			if u > 0 {
				uv = o.u[u-1]
			}
			f.Vertexes = append(f.Vertexes, types.Vertex{
				Position: vv,
				Normal:   nv,
				UV:       uv,
			})
		}
		ret.Faces = append(ret.Faces, f)
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
	return obj.ToModel()
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
