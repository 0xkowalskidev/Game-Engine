package engine

import (
	"0xKowalski/game/components"
	"0xKowalski/game/ecs"
	"0xKowalski/game/systems"
	"0xKowalski/game/window"
	"log"

	"github.com/go-gl/glfw/v3.3/glfw"
)

type Engine struct {
	Window         *window.Window
	ComponentStore *ecs.ComponentStore
	RenderSystem   *systems.RenderSystem
}

func InitEngine() (*Engine, error) {
	winConfig := window.WindowConfig{
		Title:  "Game Window",
		Width:  800,
		Height: 600,
	}

	win, err := window.InitWindow(winConfig)
	if err != nil {
		log.Printf("Failed to create window: %v", err)
		return nil, err
	}

	store := ecs.NewComponentStore()

	rs, err := systems.NewRenderSystem(win, store)
	if err != nil {
		return nil, err
	}

	engine := &Engine{
		Window:         win,
		RenderSystem:   rs,
		ComponentStore: store,
	}

	return engine, nil
}

func (e *Engine) Run() {
	entity := ecs.NewEntity()
	comp := components.NewRenderComponent()
	e.ComponentStore.AddComponent(entity, comp)

	for !e.Window.GlfwWindow.ShouldClose() {
		e.RenderSystem.Update()

		e.Window.GlfwWindow.SwapBuffers() // Swap buffers to display the frame
		glfw.PollEvents()
	}

	e.Cleanup()
}

func (e *Engine) Cleanup() {
	e.Window.Cleanup()
}
