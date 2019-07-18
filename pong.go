package main
// Xframe rate independence
// score
// Game ove State - win/lose
// 2 player vs playing computer
// Mouse? Joystick
import (
	"fmt"
	"github.com/veandco/go-sdl2/sdl"
	"log"
	"math/rand"
	"time"
)
const winWidth,winHeight = 800,600
var nums = [][][]byte {
	{{1,1,1},
	 {1,0,1},
	 {1,0,1},
	 {1,0,1},
	 {1,1,1}},
	{{0,1,0},
	 {1,1,0},
	 {0,1,0},
	 {0,1,0},
	 {1,1,1}},
	{{1,1,1},
	 {0,0,1},
	 {1,1,1},
	 {1,0,0},
	 {1,1,1}},
	{{1,1,1},
	 {0,0,1},
	 {1,1,1},
	 {0,0,1},
	 {1,1,1}},
	{{1,0,1},
	 {1,0,1},
	 {1,1,1},
	 {0,0,1},
	 {0,0,1}},
	{{1,1,1},
	{1,0,0},
	{1,1,1},
	{0,0,1},
	{1,1,1}},
}
type gameState int
const (
	start gameState = iota
	play
)
var state gameState
type color struct{
	b,g,r byte
}
type position struct{
	x,y float32 // 32 recommended
}
type ball struct{
	position
	radius float32
	vx,vy float32
	color color
}
type paddle struct {
	 position
	 w,h int
	 speed float32
	 color color
	 score int
}

func (paddle *paddle) draw(pixels []byte,negative bool){
	startX:= int(paddle.x) - paddle.w/2
	startY:= int(paddle.y) - paddle.h/2
	for y:=0;y<paddle.h;y++{
		for x:=0;x<paddle.w;x++{
			if !negative {
				setPixel(startX+x, startY+y, paddle.color, pixels)
			}else{
				setPixel(startX+x, startY+y, color{0,0,0}, pixels)
			}
		}
	}
	numX:= lerp(paddle.x,getCenter().x,20)
	if paddle.score <= 5 {
		drawText(position{numX, 35}, paddle.color, 10, paddle.score,nums, pixels)
	}
}
func (paddle *paddle) update(keyState []uint8,pixels []byte,elapsedTime float32){
	if keyState[sdl.SCANCODE_UP]!=0{
		if int(paddle.y - paddle.speed*elapsedTime) - paddle.h/2 > 0{
			paddle.draw(pixels, true)
			paddle.y -= paddle.speed*elapsedTime
		}
	}else if keyState[sdl.SCANCODE_DOWN]!=0{
		if int(paddle.y  + paddle.speed*elapsedTime) + paddle.h/2 < winHeight {
			paddle.draw(pixels,true)
			paddle.y += paddle.speed*elapsedTime
		}
	}

}
func (paddle *paddle) aiUpdate(ball *ball,pixels []byte,elapsedTime float32){
	rand.Seed(int64(time.Now().Second()))
	random := rand.Intn(200)
	if ball.x + ball.radius > float32(winWidth-400 + random) {
		if paddle.y < ball.y {
			if int(paddle.y+paddle.speed * elapsedTime)+paddle.h/2 < 600 {
				paddle.draw(pixels, true)
				paddle.y += paddle.speed*elapsedTime
			}
		} else if paddle.y > ball.y {
			if int(paddle.y-paddle.speed*elapsedTime)-paddle.h/2 > 0 {
				paddle.draw(pixels, true)
				paddle.y -= paddle.speed*elapsedTime
			}
		}
	}
}
func (ball *ball) draw(pixels []byte,negative bool){
	for y:=-ball.radius;y<ball.radius;y++{
		for x:=-ball.radius;x<ball.radius;x++ {
			if x*x+y*y < ball.radius*ball.radius { // squaring is cpu expensive
				if !negative {
					setPixel(int(ball.x+x), int(ball.y+y), ball.color, pixels)
				} else {
					setPixel(int(ball.x+x), int(ball.y+y), color{0, 0, 0}, pixels)
				}
			}
		}
	}
}
func (ball *ball) update(pixels []byte , player1,player2 *paddle, elapsedTime float32){
	ball.draw(pixels,true)
	ball.x+=ball.vx*elapsedTime
	ball.y+=ball.vy*elapsedTime
	if ball.y - ball.radius < 0 || ball.y + ball.radius > winHeight{
		ball.vy = - ball.vy
	}
	if ball.x < 0{
		player1.score++
		ball.position = getCenter()
		state =start
	}else if ball.x>winWidth{
		player2.score++
		ball.position = getCenter()
		state = start
	}
	// The minimum translation vector is the shortest distance that
	// the colliding object can be moved in order to no longer be colliding with the collidee.
		if int(ball.x-ball.radius) < int(player1.x)+player1.w/2 || int(ball.x+ball.radius) > int(player2.x)-player2.w/2 {
			if int(ball.y) > int(player1.y)-player1.h/2 && int(ball.y) < int(player1.y)+player1.h/2 {
					ball.vx = - ball.vx
					ball.x = player1.x + float32(player1.w/2) + ball.radius
			} else if int(ball.y) > int(player2.y)-player2.h/2 && int(ball.y) < int(player2.y)+player2.h/2{
				ball.vx = - ball.vx
				ball.x = player2.x - float32(player2.w/2) - ball.radius
			}
		}
}
func drawText(pos position,clr color,size,num int ,elements [][][]byte ,pixels []byte){
	startX := int(pos.x) - (size*3)/2
	startY := int(pos.y) - (size*5)/2

	for i:= 0;i<5;i++{
		for j := 0;j<3;j++{
			if elements[num][i][j]==1{
				for y:= 0;y<size;y++{
					for x:=0;x<size;x++{
						setPixel(startX+x,startY+y,clr,pixels)
					}
				}
			}else if elements[num][i][j]==0 {
				for y := 0; y < size; y++ {
					for x := 0; x < size; x++ {
						setPixel(startX+x, startY+y, color{0,0,0}, pixels)
					}
				}
			}
			startX+=size
		}
		startX-=3*size
		startY+=size
	}
}
//func drawStart(pixels []byte){
//	var letters = [][][]byte{
//		{
//			{1,1,1},
//			{1,0,0},
//			{1,1,1},
//			{0,0,1},
//			{1,1,1},
//		},
//		{
//			{1,1,1},
//			{1,1,1},
//			{0,0,0},
//			{1,1,1},
//			{1,1,1},
//		},
//		{
//			{0,0,0},
//			{1,1,1},
//			{1,1,1},
//			{1,1,1},
//			{0,0,0},
//		},
//	}
//	second:= lerp(getCenter().x,winWidth,25)
//	third := lerp(second,winWidth,25)
//	drawText(position{getCenter().x-100,winHeight/2},color{132,222,2},20,0,letters,pixels)
//	drawText(position{second,winHeight/2},color{132,222,2},10,1,letters,pixels)
//	drawText(position{third,winHeight/2},color{132,222,2},20,2,letters,pixels)
//}
func lerp(a,b float32, pct float32) float32{
	return a+(b-a)*pct*0.01
}
func setPixel(x,y int,c color,pixels []byte) {
	index := (y*winWidth + x) * 4
	// making sure there is enough room for the ones to be set.
	if index < len(pixels)-4 && index>=0 {
		pixels[index] = c.r
		pixels[index+1] = c.g
		pixels[index+2] = c.b
	}
}
func getCenter() position{
	return position{winWidth/2 , winHeight/2}
}
func getPixelColor(x,y int, pixels []byte) color{
	var color color
	index := (y*winWidth + x) * 4
	if index < len(pixels)-4 && index>=0 {
		color.r = pixels[index]
		color.g = pixels[index+1]
		color.b = pixels[index+2]
	}
	return color
}
func main() {
	window, err := sdl.CreateWindow("Start", sdl.WINDOWPOS_UNDEFINED, sdl.WINDOWPOS_UNDEFINED, winWidth, winHeight, sdl.WINDOW_SHOWN)
	if err!=nil{
		fmt.Println(err)
		return
	}
	defer window.Destroy()
	renderer,err:= sdl.CreateRenderer(window,-1,sdl.RENDERER_ACCELERATED)
	if err!=nil{
		log.Println(err)
		return
	}
	defer renderer.Destroy()
	tex,err:= renderer.CreateTexture(sdl.PIXELFORMAT_ARGB8888,sdl.TEXTUREACCESS_STREAMING,winWidth,winHeight)
	if err!=nil{
		log.Println(err)
		return
	}
	defer tex.Destroy()
	pixels := make([]byte, winWidth*winHeight*4)
	player1 := paddle{position{10,50},20,100,400,color{b:38,g:0,r:230},0}
	player2 := paddle{position{790,50},20,100,300,color{b:175,g:0,r:42},0}
	ball := ball{position{300,300},20,400,400,color{b:255,g:191,r:0}}
	keyState := sdl.GetKeyboardState()
	var frameStart time.Time
	var elapsedTime float32
	for {
		frameStart = time.Now()
		for event := sdl.PollEvent(); event!=nil;event = sdl.PollEvent() {
			switch event.(type){
			case *sdl.QuitEvent :
				return
			}
		}
		if state == play {
			player1.update(keyState, pixels, elapsedTime)
			player2.aiUpdate(&ball, pixels, elapsedTime)
			ball.update(pixels, &player1, &player2, elapsedTime)

			player1.draw(pixels, false)
			player2.draw(pixels, false)
			ball.draw(pixels, false)
		}else if state == start{
			if keyState[sdl.SCANCODE_SPACE] !=0{
				if player2.score == 5 || player1.score==5{
					player1.score = 0
					player2.score = 0
				}
				state = play
			}
		}
		_ =tex.Update(nil,pixels,winWidth*4)
		_ = renderer.Copy(tex,nil,nil)
		renderer.Present()
		elapsedTime = float32(time.Since(frameStart).Seconds())
		if elapsedTime < 0.005{
			sdl.Delay(5-uint32(elapsedTime*1000))
			elapsedTime = float32(time.Since(frameStart).Seconds())
		}
	}
}