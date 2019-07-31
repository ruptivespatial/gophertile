package gophertile

import (
	"math"
)

const threeSixty float64 = 360.0
const oneEighty float64 = 180.0
const radius float64 = 6378137.0
const d2r float64 = math.Pi / 180
const r2d = 180 / math.Pi

//Tile struct is the main object we deal with, represents a standard X/Y/Z tile
type Tile struct {
	X, Y, Z int
}
type tileFraction struct {
	X, Y float64
	Z    int
}

//LngLat holds a standard geographic coordinate pair in decimal degrees
type LngLat struct {
	Lng, Lat float64
}

//LngLatBbox bounding box of a tile, in decimal degrees
type LngLatBbox struct {
	West, South, East, North float64
}

//Bbox holds Spherical Mercator bounding box of a tile
type Bbox struct {
	Left, Bottom, Right, Top float64
}

//XY holds a Spherical Mercator point
type XY struct {
	X, Y float64
}

func deg2rad(deg float64) float64 {
	return deg * (math.Pi / oneEighty)
}
func rad2deg(rad float64) float64 {
	return rad * (oneEighty / math.Pi)
}

//GetTile returns a tile for a given longitude latitude and zoom level
func GetTile(lng float64, lat float64, zoom int) *Tile {

	latRad := deg2rad(lat)
	n := math.Pow(2.0, float64(zoom))
	x := int(math.Floor((lng + oneEighty) / threeSixty * n))
	y := int(math.Floor((1.0 - math.Log(math.Tan(latRad)+(1.0/math.Cos(latRad)))/math.Pi) / 2.0 * n))

	return &Tile{x, y, zoom}

}

//Equals compares 2 tiles
func (tile *Tile) Equals(t2 *Tile) bool {

	return tile.X == t2.X && tile.Y == t2.Y && tile.Z == t2.Z

}

//Ul returns the upper left corner of the tile decimal degrees
func (tile *Tile) Ul() *LngLat {

	n := math.Pow(2.0, float64(tile.Z))
	lonDeg := float64(tile.X)/n*threeSixty - oneEighty
	latRad := math.Atan(math.Sinh(math.Pi * float64(1-(2*float64(tile.Y)/n))))
	latDeg := rad2deg(latRad)

	return &LngLat{lonDeg, latDeg}
}

//Bounds returns a LngLatBbox for a given tile
func (tile *Tile) Bounds() *LngLatBbox {
	a := tile.Ul()
	shifted := Tile{tile.X + 1, tile.Y + 1, tile.Z}
	b := shifted.Ul()
	return &LngLatBbox{a.Lng, b.Lat, b.Lng, a.Lat}
}

//Parent returns the tile above (i.e. at a lower zoon number) the given tile
func (tile *Tile) Parent() *Tile {

	if tile.Z == 0 && tile.X == 0 && tile.Y == 0 {
		return tile
	}

	if math.Mod(float64(tile.X), 2) == 0 && math.Mod(float64(tile.Y), 2) == 0 {
		return &Tile{tile.X / 2, tile.Y / 2, tile.Z - 1}
	}
	if math.Mod(float64(tile.X), 2) == 0 {
		return &Tile{tile.X / 2, (tile.Y - 1) / 2, tile.Z - 1}
	}
	if math.Mod(float64(tile.X), 2) != 0 && math.Mod(float64(tile.Y), 2) != 0 {
		return &Tile{(tile.X - 1) / 2, (tile.Y - 1) / 2, tile.Z - 1}
	}
	if math.Mod(float64(tile.X), 2) != 0 && math.Mod(float64(tile.Y), 2) == 0 {
		return &Tile{(tile.X - 1) / 2, tile.Y / 2, tile.Z - 1}
	}
	return nil
}

//Children returns the 4 tiles below (i.e. at a higher zoom number) the given tile
func (tile *Tile) Children() []*Tile {

	kids := []*Tile{
		{tile.X * 2, tile.Y * 2, tile.Z + 1},
		{tile.X*2 + 1, tile.Y * 2, tile.Z + 1},
		{tile.X*2 + 1, tile.Y*2 + 1, tile.Z + 1},
		{tile.X * 2, tile.Y*2 + 1, tile.Z + 1},
	}
	return kids
}

//ToXY transforms WGS84 DD to Spherical Mercator meters
func ToXY(ll *LngLat) *XY {

	x := radius * deg2rad(ll.Lng)
	intrx := (math.Pi * 0.25) + (0.5 * deg2rad(ll.Lat))
	y := radius * math.Log(math.Tan(intrx))

	return &XY{x, y}
}

//BboxToTile returns the smallest tile which will fit the entire bounding box
func BboxToTile(box *LngLatBbox) *Tile {

	min := PointToTile(&LngLat{box.West, box.South}, 32)
	max := PointToTile(&LngLat{box.East, box.North}, 32)

	tilePoints := []int{min.X, min.Y, max.X, max.Y}
	z := getBBoxZoom(tilePoints)
	if z == 0 {
		return &Tile{0, 0, 0}
	}
	x := tilePoints[0] >> uint(32-z)
	y := tilePoints[1] >> uint(32-z)
	return &Tile{X: x, Y: y, Z: z}
}

//getBBoxZoom returns the lowest zoom level that will constrain the provided tile numbers
func getBBoxZoom(tc []int) int {

	maxZoom := 28
	for z := 0; z < maxZoom; z++ {
		mask := 1 << uint(32-(z+1))
		if (tc[0]&mask != tc[2]&mask) || (tc[1]&mask != tc[3]&mask) {
			return z
		}
	}
	return maxZoom

}

//PointToTile returns a tile at the giver lat/lng and zoom level
func PointToTile(ll *LngLat, z int) *Tile {
	tile := pointToFractionalTile(ll, z)
	return &Tile{X: int(math.Floor(tile.X)),
		Y: int(math.Floor(tile.Y)),
		Z: z}
}

//PointToFractionalTile returns a tile for the giver lat/lng and zoom -- however it also returns tile decimals
//which might not be useful, will perhaps mark this unexported in the furture
func pointToFractionalTile(ll *LngLat, z int) *tileFraction {
	sin := math.Sin(ll.Lat * d2r)
	z2 := math.Pow(2, float64(z))
	x := z2 * (ll.Lng/360 + 0.5)
	y := z2 * (0.5 - 0.25*math.Log((1+sin)/(1-sin))/math.Pi)

	// Wrap Tile X
	x = math.Mod(x, z2)

	if x < 0 {
		x = x + z2
	}

	return &tileFraction{X: x, Y: y, Z: z}

}
