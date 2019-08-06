package gophertile

import (
	"fmt"
	"math"
	"testing"
)

func TestTile_Ul(t *testing.T) {

	tile := Tile{486, 332, 10}
	ll := tile.Ul()
	expected := LngLat{-9.140625, 53.33087298301705}

	assertPrettyClose(t, ll.Lng, expected.Lng)
	assertPrettyClose(t, ll.Lat, expected.Lat)

}

func TestTile_Bounds(t *testing.T) {

	tile := Tile{486, 332, 10}
	expected := LngLatBbox{-9.140625, 53.120405283106564, -8.7890625, 53.33087298301705}
	bbox := tile.Bounds()
	assertPrettyClose(t, expected.East, bbox.East)
	assertPrettyClose(t, expected.West, bbox.West)
	assertPrettyClose(t, expected.North, bbox.North)
	assertPrettyClose(t, expected.South, bbox.South)
}

func TestTile_Parent(t *testing.T) {
	expected := Tile{243, 166, 9}
	tile := Tile{486, 332, 10}
	parent := tile.Parent()
	assertEq(t, expected.X, parent.X)
	assertEq(t, expected.Y, parent.Y)
	assertEq(t, expected.Z, parent.Z)
}

func TestTile_Children(t *testing.T) {

	tile := Tile{246, 166, 9}
	expected := Tile{492, 332, 10}
	children := tile.Children()

	found := false
	for _, t2 := range children {
		if t2.Equals(&expected) {
			found = true
		}
	}
	if !found {
		t.Fail()
	}

}

func TestToXY(t *testing.T) {

	expected := XY{-1017529.7205322663, 7044436.526761846}
	tile := Tile{486, 332, 10}
	//expected := XY{-0.0,0.0}

	ll := tile.Ul()
	//xy := ToXY(&LngLat{0.0,0.0})
	xy := ToXY(ll)

	assertPrettyClose(t, xy.Y, expected.Y)
	assertPrettyClose(t, xy.X, expected.X)

}

func TestGetTile(t *testing.T) {

	tile := GetTile(20.6852, 40.1222, 9)
	expected := Tile{285, 193, 9}

	assertEq(t, tile.Z, expected.Z)
	assertEq(t, tile.Y, expected.Y)
	assertEq(t, tile.X, expected.X)

}

func assertPrettyClose(t *testing.T, x float64, y float64) {

	if x != y {
		if math.Abs(x-y) < 0.00000001 {
			//floats and shit
			return
		}
		t.Logf("Expected: %v", x)
		t.Logf("Actual: %v", y)
		t.Logf("Difference between actual and expected is %v", x-y)
		t.Fail()
	}
}

func assertEq(t *testing.T, x interface{}, y interface{}) {
	if x != y {
		t.Logf("Expected: %v", x)
		t.Logf("Actual: %v", y)
		fmt.Printf("%v is not equal to %v", x, y)
		t.Fail()
	}
}
func TestBboxToTile(t *testing.T) {
	bbox := LngLatBbox{-77.04615354537964,
		38.899967510782346,
		-77.03664779663086,
		38.90728142481329}

	tile := BboxToTile(&bbox)

	if tile.X != 9371 || tile.Y != 12534 || tile.Z != 15 {
		t.Logf("tile: %v", tile)
		t.Fail()
	}

}

func TestPointToTile(t *testing.T) {
	ll := LngLat{
		Lat: 41.26000108568697,
		Lng: -95.93965530395508,
	}
	tile := Tile{X: 119, Y: 191, Z: 9}

	tf := PointToTile(&ll, 9)

	if tile.X != tf.X || tile.Y != tf.Y || tile.Z != tf.Z {
		t.Fail()
	}

}
func TestBounds3857(t *testing.T) {

	t1 := Tile{0, 0, 2}
	correct1 := Bbox{Left: -20037508, Bottom: 10018754, Right: -10018754, Top: 20037508}
	t2 := Tile{16, 29, 6}
	correct2 := Bbox{-10018754, 1252344, -9392582, 1878516}
	resBBox1 := t1.Bounds3857()
	resBBox2 := t2.Bounds3857()

	//truncating the decimals before comparing
	assertEq(t, float64(int(resBBox1.Bottom)), correct1.Bottom)
	assertEq(t, float64(int(resBBox2.Bottom)), correct2.Bottom)
	assertEq(t, float64(int(resBBox1.Top)), correct1.Top)
	assertEq(t, float64(int(resBBox2.Top)), correct2.Top)
	assertEq(t, float64(int(resBBox1.Left)), correct1.Left)
	assertEq(t, float64(int(resBBox2.Left)), correct2.Left)
	assertEq(t, float64(int(resBBox1.Right)), correct1.Right)
	assertEq(t, float64(int(resBBox2.Right)), correct2.Right)

}

func TestPointToFractionalTile(t *testing.T) {
	ll := LngLat{
		Lat: 41.26000108568697,
		Lng: -95.93965530395508,
	}
	tf := tileFraction{X: 119.552490234375, Y: 191.47119140625, Z: 9}
	tile := pointToFractionalTile(&ll, 9)
	if tile.X != tf.X || tile.Y != tf.Y || tile.Z != tf.Z {

		t.Fail()
	}

}
