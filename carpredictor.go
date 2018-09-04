package main

import "fmt"

type Car struct {
	model string
	sites uint
	fuel  uint
}

func raceWinnerPredictor(cars *[2]Car) string {

	var topFuel uint
	var winnerModel string

	for _, car := range cars {
		if car.fuel > topFuel {
			topFuel = car.fuel
			winnerModel = car.model
		}
	}
	return winnerModel
}

func main() {

	car1 := Car{model: "BMW", sites: 4, fuel: 100}
	car2 := Car{model: "JAGUAR", sites: 2, fuel: 90}
	var race [2]Car
	race[0] = car1
	race[1] = car2
	result := raceWinnerPredictor(&race)
	fmt.Printf("Winner model is %v !!!", result)
}
