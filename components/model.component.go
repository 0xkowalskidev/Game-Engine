package components

import (
	"log"

	"github.com/go-gl/mathgl/mgl32"
	"github.com/udhos/gwob"
)

type ModelComponent struct {
	MeshComponents     []*MeshComponent
	MaterialComponents []*MaterialComponent
	BufferComponents   []*BufferComponent
}

func NewModelComponent(objPath, mtlPath string) *ModelComponent {
	options := &gwob.ObjParserOptions{
		LogStats: false,
		Logger:   func(msg string) {},
	}
	// Load OBJ
	obj, err := gwob.NewObjFromFile(objPath, options)
	if err != nil {
		panic(err)
	}

	// Load material lib
	lib, err := gwob.ReadMaterialLibFromFile(mtlPath, options)
	if err != nil {
		panic(err)
	}

	meshComponents, materialComponents, bufferComponents := ConvertObjToMeshComponents(obj, &lib)

	return &ModelComponent{
		MeshComponents:     meshComponents,
		MaterialComponents: materialComponents,
		BufferComponents:   bufferComponents,
	}
}

func ConvertObjToMeshComponents(obj *gwob.Obj, lib *gwob.MaterialLib) ([]*MeshComponent, []*MaterialComponent, []*BufferComponent) {
	var meshComponents []*MeshComponent
	var materialComponents []*MaterialComponent
	var bufferComponents []*BufferComponent

	// Scan OBJ groups
	for _, g := range obj.Groups {
		var vertices []Vertex
		intIndices := obj.Indices[g.IndexBegin : g.IndexBegin+g.IndexCount]
		indices := make([]uint32, len(intIndices))

		for i, index := range intIndices {
			indices[i] = uint32(index)
		}

		stride := obj.StrideSize / 4
		strideOffsetTex := obj.StrideOffsetTexture / 4
		strideOffsetNorm := obj.StrideOffsetNormal / 4

		for i := 0; i < len(obj.Coord); i += stride {
			// Extract position always assumed to be present at the beginning
			pos := mgl32.Vec3{obj.Coord[i], obj.Coord[i+1], obj.Coord[i+2]}

			var tex mgl32.Vec2
			if obj.TextCoordFound && i+obj.StrideOffsetTexture+1 < len(obj.Coord) {
				tex = mgl32.Vec2{obj.Coord[i+strideOffsetTex], obj.Coord[i+strideOffsetTex+1]}
			}

			var norm mgl32.Vec3
			if obj.NormCoordFound && i+obj.StrideOffsetNormal+2 < len(obj.Coord) {
				norm = mgl32.Vec3{obj.Coord[i+strideOffsetNorm], obj.Coord[i+strideOffsetNorm+1], obj.Coord[i+strideOffsetNorm+2]}
			}

			vertices = append(vertices, Vertex{Position: pos, TexCoords: tex, Normal: norm})
		}

		meshComponents = append(meshComponents, NewMeshComponent(vertices, indices))

		mtl, found := lib.Lib[g.Usemtl]
		if !found {
			log.Fatal("mtl not found")
		}

		materialComponent, err := NewMaterialComponent("assets/models/backpack/"+mtl.MapKd, "assets/models/backpack/"+mtl.MapKs, mtl.Ns)
		if err != nil {
			panic(err)
		}

		materialComponents = append(materialComponents, materialComponent)

		bufferComponents = append(bufferComponents, NewBufferComponent(vertices, indices))
	}

	return meshComponents, materialComponents, bufferComponents
}