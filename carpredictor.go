package main

import "fmt"

type Race struct {
	cars     *[]Car
	distance uint
}

type Car struct {
	model string
	fuel  uint
	hp    uint
}

func RaceWinnerPredictor(race *Race) string {

	perFuelPoint := +1
	perHpPoint := +3
	fuelPerHpRatio := 0.5

	var winnerModel string
	var maximumPointsReached int

	for _, car := range *race.cars {

		fuelEffort := (fuelPerHpRatio * float64(car.hp)) * float64(race.distance)
		if fuelEffort > float64(car.fuel) {
			continue
		}

		pointsSum := 0
		pointsSum = pointsSum + int(car.fuel)*perFuelPoint
		pointsSum = pointsSum + int(car.hp)*perHpPoint

		if pointsSum > maximumPointsReached {
			maximumPointsReached = pointsSum
			winnerModel = car.model
		}
	}
	return winnerModel
}

func main() {

	car1 := Car{model: "BMW", fuel: 100, hp: 200}
	car2 := Car{model: "JAGUAR", fuel: 90, hp: 270}
	race := Race{&[]Car{car1, car2}, 100}
	result := RaceWinnerPredictor(&race)
	fmt.Printf("Winner model is %v !!!", result)
}
