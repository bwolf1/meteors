package game

import (
	"math"
	"time"

	"github.com/hajimehoshi/ebiten/v2"

	"github.com/bwolf1/meteors/assets"
)

const (
	shootCooldown     = time.Millisecond * 500
	rotationPerSecond = math.Pi * 1.25 // Radians per second

	bulletSpawnOffset = 50.0
)

type Player struct {
	game *Game

	position Vector
	rotation float64
	sprite   *ebiten.Image

	shootCooldown *Timer
}

func NewPlayer(game *Game) *Player {
	sprite := assets.PlayerSprite

	// Get the sprite bounds to calculate the center of the sprite
	bounds := sprite.Bounds()
	halfWidth := float64(bounds.Dx()) / 2
	halfHeight := float64(bounds.Dy()) / 2
	
	// Center the position of the center of the sprite on the screen
	pos := Vector{
		X: screenWidth/2 - halfWidth,
		Y: screenHeight/2 - halfHeight,
	}

	return &Player{
		game:          game,
		position:      pos,
		rotation:      0,
		sprite:        sprite,
		shootCooldown: NewTimer(shootCooldown),
	}
}

func (p *Player) Update() {
	speed := rotationPerSecond / float64(ebiten.TPS())

	// Rotate left/right
	if ebiten.IsKeyPressed(ebiten.KeyLeft) {
		p.rotation -= speed
	}
	if ebiten.IsKeyPressed(ebiten.KeyRight) {
		p.rotation += speed
	}

	// Move forward
	if ebiten.IsKeyPressed(ebiten.KeyArrowUp) {
		// Move in the direction the player is facing
		p.position.X += math.Sin(p.rotation) * 3
		p.position.Y += math.Cos(p.rotation) * -3

		// Limit the player's position to the screen bounds
		p.position.X = math.Max(0, math.Min(screenWidth-float64(p.sprite.Bounds().Dx()), p.position.X))
		p.position.Y = math.Max(0, math.Min(screenHeight-float64(p.sprite.Bounds().Dy()), p.position.Y))
	}

	// Move backward
	if ebiten.IsKeyPressed(ebiten.KeyArrowDown) {
		// Move in the opposite direction the player is facing
		p.position.X += math.Sin(p.rotation+math.Pi) * 3
		p.position.Y += math.Cos(p.rotation+math.Pi) * -3

		// Limit the player's position to the screen bounds
		p.position.X = math.Max(0, math.Min(screenWidth-float64(p.sprite.Bounds().Dx()), p.position.X))
		p.position.Y = math.Max(0, math.Min(screenHeight-float64(p.sprite.Bounds().Dy()), p.position.Y))
	}

	p.shootCooldown.Update()
	if p.shootCooldown.IsReady() && ebiten.IsKeyPressed(ebiten.KeySpace) {
		p.shootCooldown.Reset()

		bounds := p.sprite.Bounds()
		halfWidth := float64(bounds.Dx()) / 2
		halfHeight := float64(bounds.Dy()) / 2

		spawnPos := Vector{
			p.position.X + halfWidth + math.Sin(p.rotation)*bulletSpawnOffset,
			p.position.Y + halfHeight + math.Cos(p.rotation)*-bulletSpawnOffset,
		}

		bullet := NewBullet(spawnPos, p.rotation)
		p.game.AddBullet(bullet)
	}
}

func (p *Player) Draw(screen *ebiten.Image) {
	bounds := p.sprite.Bounds()
	halfWidth := float64(bounds.Dx()) / 2
	halfHeight := float64(bounds.Dy()) / 2

	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(-halfWidth, -halfHeight)
	op.GeoM.Rotate(p.rotation)
	op.GeoM.Translate(halfWidth, halfHeight)

	op.GeoM.Translate(p.position.X, p.position.Y)

	screen.DrawImage(p.sprite, op)
}

func (p *Player) Collider() Rect {
	bounds := p.sprite.Bounds()

	return NewRect(
		p.position.X,
		p.position.Y,
		float64(bounds.Dx()),
		float64(bounds.Dy()),
	)
}
