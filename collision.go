package geom

// CollisionRectangles checks if the given rectangles collide.
func CollisionRectangles[T Number](rect1 Rectangle[T], rect2 Rectangle[T]) bool {
	min1, max1 := rect1.Min(), rect1.Max()
	min2, max2 := rect2.Min(), rect2.Max()

	return min1.X <= max2.X && min2.X <= max1.X && min1.Y <= max2.Y && min2.Y <= max1.Y
}

// CollisionCircles checks if the given circles collide. Touching circles collide,
// the same closed convention as CollisionRectangles.
func CollisionCircles[T Number](circle1 Circle[T], circle2 Circle[T]) bool {
	distance := circle1.Center.Subtract(circle2.Center)
	threshold := circle1.Radius + circle2.Radius

	return distance.LessOrEqual(threshold)
}

// CollisionRectangleCircle checks if the given rectangle and circle collide: the point of the
// rectangle closest to the circle center lies within the radius. Touching shapes collide, and
// the rectangle bounds are the same Min and Max that Contains and CollisionRectangles use.
func CollisionRectangleCircle[T Number](rect Rectangle[T], circle Circle[T]) bool {
	closest := rect.Clamp(circle.Center)

	return closest.Subtract(circle.Center).LessOrEqual(circle.Radius)
}
