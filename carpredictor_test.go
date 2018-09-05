package main

import "testing"

func TestRaceWinnerPredictor(t *testing.T) {

	// Must win Jaguar due to BUGATTI is out of fuel.
	car1 := Car{model: "BMW", fuel: 50000, hp: 200}
	car2 := Car{model: "JAGUAR", fuel: 67500, hp: 270}
	car3 := Car{model: "BUGATTI", fuel: 124999, hp: 500}
	race := Race{&[]Car{car1, car2, car3}, 500}

	result := RaceWinnerPredictor(&race)

	if "JAGUAR" != result {
		t.Errorf("This car is not expected to win: %v", result)
	}
}
