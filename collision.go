package geom

// CollisionRectangles checks if the given rectangles collide.
func CollisionRectangles[T Number](rect1 Rectangle[T], rect2 Rectangle[T]) bool {
	min1, max1 := rect1.Min(), rect1.Max()
	min2, max2 := rect2.Min(), rect2.Max()

	return min1.X <= max2.X && min2.X <= max1.X && min1.Y <= max2.Y && min2.Y <= max1.Y
}

// CollisionCircles checks if the given circles collide.
func CollisionCircles[T Number](circle1 Circle[T], circle2 Circle[T]) bool {
	distance := circle1.Center.Subtract(circle2.Center)
	threshold := circle1.Radius + circle2.Radius

	return distance.Less(threshold)
}

// CollisionRectangleCircle checks if the given rectangle and circle collide.
func CollisionRectangleCircle[T Number](rect Rectangle[T], circle Circle[T]) bool {
	distance := circle.Center.Subtract(rect.Center).Abs()
	extends := rect.Size.Scale(0.5).Vector()

	// circle center is more than size outside rectangle borders
	if distance.X > extends.X+circle.Radius || distance.Y > extends.Y+circle.Radius {
		return false
	}

	// circle center is in rectangle
	if distance.X <= extends.X || distance.Y <= extends.Y {
		return true
	}

	// circle center is less than its size outside nearest border
	return distance.Subtract(extends).Less(circle.Radius)
}
