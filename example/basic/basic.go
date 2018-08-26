package main

import (
	"github.com/gopherjs/gopherjs/js"
	"github.com/gopherjs/vecty"
	"github.com/gopherjs/vecty/elem"
	"github.com/lngramos/three"
	"github.com/vecty/vthree"
)

func main() {
	vecty.RenderBody(&PageView{})
}

// PageView is our main page component.
type PageView struct {
	vecty.Core

	scene    *three.Scene
	camera   three.PerspectiveCamera
	renderer *three.WebGLRenderer
	mesh     *three.Mesh
}

// Render implements the vecty.Component interface.
func (p *PageView) Render() vecty.ComponentOrHTML {
	return elem.Body(
		elem.Heading1(
			vecty.Markup(
				vecty.Style("margin", "0"),
			),
			vecty.Text("Vecty + three.js = ♡"),
		),
		vthree.WebGLRenderer(vthree.WebGLOptions{
			Init:     p.init,
			Shutdown: p.shutdown,
		}),
	)
}

func (p *PageView) shutdown(renderer *three.WebGLRenderer) {
	// After shutdown, we shouldn't use any of these anymore.
	p.scene = nil
	p.camera = three.PerspectiveCamera{}
	p.renderer = nil
	p.mesh = nil
}

func (p *PageView) init(renderer *three.WebGLRenderer) {
	p.renderer = renderer

	windowWidth := js.Global.Get("innerWidth").Float()
	windowHeight := js.Global.Get("innerHeight").Float()
	devicePixelRatio := js.Global.Get("devicePixelRatio").Float()

	p.camera = three.NewPerspectiveCamera(70, windowWidth/windowHeight, 1, 1000)
	p.camera.Position.Set(0, 0, 400)

	p.scene = three.NewScene()

	light := three.NewDirectionalLight(three.NewColor(126, 255, 255), 0.5)
	light.Position.Set(256, 256, 256).Normalize()
	p.scene.Add(light)

	p.renderer.SetPixelRatio(devicePixelRatio)
	p.renderer.SetSize(windowWidth, windowHeight, true)

	// Create cube
	geometry := three.NewBoxGeometry(&three.BoxGeometryParameters{
		Width:  128,
		Height: 128,
		Depth:  128,
	})

	// geometry2 := three.NewCircleGeometry(three.CircleGeometryParameters{
	// 	Radius:      50,
	// 	Segments:    20,
	// 	ThetaStart:  0,
	// 	ThetaLength: 2,
	// })

	materialParams := three.NewMaterialParameters()
	materialParams.Color = three.NewColor(0, 123, 211)
	materialParams.Shading = three.SmoothShading
	materialParams.Side = three.FrontSide
	material := three.NewMeshBasicMaterial(materialParams)
	// material := three.NewMeshLambertMaterial(materialParams)
	// material := three.NewMeshPhongMaterial(materialParams)
	p.mesh = three.NewMesh(geometry, material)
	p.scene.Add(p.mesh)

	// Begin animating.
	p.animate()
}

func (p *PageView) animate() {
	if p.renderer == nil {
		// We shutdown, stop animation.
		return
	}
	js.Global.Call("requestAnimationFrame", p.animate)

	pos := p.mesh.Object.Get("rotation")
	pos.Set("x", pos.Get("x").Float()+float64(0.01))
	pos.Set("y", pos.Get("y").Float()+float64(0.01))

	p.renderer.Render(p.scene, p.camera)
}
