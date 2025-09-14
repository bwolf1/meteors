package game

import (
	"math"

	"github.com/hajimehoshi/ebiten/v2"

	"github.com/bwolf1/meteors/assets"
)

const (
	bulletSpeedPerSecond = 350.0
)

type Bullet struct {
	position Vector
	rotation float64
	sprite   *ebiten.Image
}

// NewBullet creates a new bullet at the given position and rotation.
func NewBullet(pos Vector, rotation float64) *Bullet {
	sprite := assets.LaserSprite

	// Offset the position by half the sprite size to center it
	bounds := sprite.Bounds()
	halfWidth := float64(bounds.Dx()) / 2
	halfHeight := float64(bounds.Dy()) / 2

	pos.X -= halfWidth
	pos.Y -= halfHeight

	b := &Bullet{
		position: pos,
		rotation: rotation,
		sprite:   sprite,
	}

	return b
}

func (b *Bullet) Update() {
	speed := bulletSpeedPerSecond / float64(ebiten.TPS())

	// Move in the direction of the rotation
	b.position.X += math.Sin(b.rotation) * speed
	b.position.Y += math.Cos(b.rotation) * -speed
}

func (b *Bullet) Draw(screen *ebiten.Image) {
	// Center the rotation on the sprite
	bounds := b.sprite.Bounds()
	halfWidth := float64(bounds.Dx()) / 2
	halfHeight := float64(bounds.Dy()) / 2

	// Draw the sprite with rotation
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(-halfWidth, -halfHeight)
	op.GeoM.Rotate(b.rotation)
	op.GeoM.Translate(halfWidth, halfHeight)

	op.GeoM.Translate(b.position.X, b.position.Y)

	screen.DrawImage(b.sprite, op)
}

func (b *Bullet) Collider() Rect {
	bounds := b.sprite.Bounds()

	return NewRect(
		b.position.X,
		b.position.Y,
		float64(bounds.Dx()),
		float64(bounds.Dy()),
	)
}
