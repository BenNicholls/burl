package burl

import "errors"
import "runtime"
import "github.com/veandco/go-sdl2/sdl"
import "github.com/bennicholls/burl/console"

var gameState State

//Initializes the game State. Call before running the game loop.
func InitState(m State) {
	gameState = m
}

//The Big Enchelada! This is the gameloop that runs everything. Make sure to run burl.InitMode() and console.Setup() before beginning the game!
func GameLoop() error {
	//TODO: implement that horrible thread job queue thing from the go-sdl2 package
	runtime.LockOSThread() //fixes some kind of go-sdl2 based thread release bug.

	if !console.Ready {
		return errors.New("Console not set up. Run burl.console.Setup() before starting game loop!")
	}

	if gameState == nil {
		return errors.New("No gameState initialized. Run burl.InitState() before starting game loop!")
	}

	var event sdl.Event
	running := true

	for running {
		for event = sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
			switch t := event.(type) {
			case *sdl.QuitEvent:
				running = false
			case *sdl.WindowEvent:
				if t.Event == sdl.WINDOWEVENT_RESTORED {
					console.ForceRedraw()
				}
			// case *sdl.MouseMotionEvent:
			// 	fmt.Printf("[%d ms] MouseMotion\ttype:%d\tid:%d\tx:%d\ty:%d\txrel:%d\tyrel:%d\n",
			// 		t.Timestamp, t.Type, t.Which, t.X, t.Y, t.XRel, t.YRel)
			// case *sdl.MouseButtonEvent:
			// 	fmt.Printf("[%d ms] MouseButton\ttype:%d\tid:%d\tx:%d\ty:%d\tbutton:%d\tstate:%d\n",
			// 		t.Timestamp, t.Type, t.Which, t.X, t.Y, t.Button, t.State)
			// case *sdl.MouseWheelEvent:
			// 	fmt.Printf("[%d ms] MouseWheel\ttype:%d\tid:%d\tx:%d\ty:%d\n",
			// 		t.Timestamp, t.Type, t.Which, t.X, t.Y)
			case *sdl.KeyUpEvent:
				//fmt.Printf("[%d ms] Keyboard\ttype:%d\tsym:%c\tmodifiers:%d\tstate:%d\trepeat:%d\n",
				//	t.Timestamp, t.Type, t.Keysym.Sym, t.Keysym.Mod, t.State, t.Repeat)
				gameState.HandleKeypress(t.Keysym.Sym)
			}
		}

		gameState.Update()

		//TODO: get console.Render() running in another thread (i think this is a good idea... maybe?)
		gameState.Render()
		console.Render()
	}

	return nil
}

//Defines a game state (level, menu, anything that can take input, update itself, render to screen.)
type State interface {
	HandleKeypress(sdl.Keycode)
	Update()
	Render()
}

//base state object, compose states around this if you want
type BurlState struct {
	tick int //update ticks since init
}

func (b BurlState) HandleKeypress(key sdl.Keycode) {

}

func (b *BurlState) Update() {
	b.tick++
}

func (b BurlState) Render() {

}
