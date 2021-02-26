package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"strconv"
	"strings"
)

type genInfo struct {
	loopCount         int
	streetCount       int
	carCount          int
	intersectionCount int
	points            int
}

type trafficLight struct {
	fromStreet *street
	time       int
}

type intersection struct {
	id            int
	streets       []*street
	trafficLights []*trafficLight
}

type car struct {
	paths      []string
	currentPos int
}

type street struct {
	name              string
	duration          int
	cars              []*car
	trafficLightState bool
	startIntersection *intersection
	endIntersection   *intersection
}

var (
	allIntersections = make(map[int]*intersection)
	allStreets       = make(map[string]*street)
)

func main() {
	file, err := ioutil.ReadFile("input.txt")
	if err != nil {
		log.Fatal(err)
	}
	parsed := strings.Split(string(file), "\n")

	generalInfo := strings.Split(parsed[0], " ")
	loop, _ := strconv.Atoi(generalInfo[0])
	scount, _ := strconv.Atoi(generalInfo[1])
	ccount, _ := strconv.Atoi(generalInfo[2])
	icounf, _ := strconv.Atoi(generalInfo[3])
	points, _ := strconv.Atoi(generalInfo[4])

	gi := genInfo{
		loopCount:         loop,
		streetCount:       scount,
		carCount:          ccount,
		intersectionCount: icounf,
		points:            points,
	}

	for _, s := range parsed[1:scount] {
		streetLine := strings.Split(s, " ")
		startIntersectionID, _ := strconv.Atoi(streetLine[0])
		startIntersection := &intersection{
			id:            startIntersectionID,
			streets:       []*street{},
			trafficLights: []*trafficLight{},
		}

		endIntersectionID, _ := strconv.Atoi(streetLine[1])
		endIntersection := &intersection{
			id:            endIntersectionID,
			streets:       []*street{},
			trafficLights: []*trafficLight{},
		}

		var st *street
		if s, ok := allStreets[streetLine[3]]; ok {
			st = s
		}

		if st == nil {
			streetDur, _ := strconv.Atoi(streetLine[3])
			st = &street{
				name:              streetLine[2],
				duration:          streetDur,
				cars:              []*car{},
				trafficLightState: false,
				startIntersection: nil,
				endIntersection:   nil,
			}
			allStreets[st.name] = st
		}

		startIntersection.streets = append(startIntersection.streets, st)
		st.startIntersection = startIntersection
		st.endIntersection = endIntersection

		if ic, ok := allIntersections[startIntersectionID]; !ok {
			allIntersections[startIntersectionID] = startIntersection
		} else {
			ic.streets = append(ic.streets, st)
		}

		if ic, ok := allIntersections[endIntersectionID]; !ok {
			allIntersections[endIntersectionID] = endIntersection
		} else {
			ic.streets = append(ic.streets, st)
		}
	}

	for _, carLine := range parsed[scount:] {
		line := strings.Split(carLine, " ")
		nCar := car{currentPos: 1}
		for _, streetName := range line[1:] {
			nCar.paths = append(nCar.paths, streetName)
		}
		st, _ := allStreets[line[1]]
		st.cars = append(st.cars, &nCar)
	}
	ioutil.WriteFile("result.txt", []byte(run(gi)), 777)
	fmt.Print()
}

func setSchedule() {
	for _, i := range allIntersections {
		var max int
		var StreetWithMostCars *street
		for _, s := range i.streets {
			if len(s.cars) > max {
				max = len(s.cars)
				StreetWithMostCars = s
			}
		}
		if StreetWithMostCars != nil {
			i.trafficLights = []*trafficLight{{
				fromStreet: StreetWithMostCars,
				time:       len(StreetWithMostCars.cars),
			}}
		}
	}
}

func run(g genInfo) string {
	var points int
	var op string
	var trafficLightCount int
	for i := 0; i < g.loopCount; i++ {
		setSchedule()
		for _, intersec := range allIntersections {
			if len(intersec.streets[0].cars) == 0 || len(intersec.trafficLights) == 0 {
				continue
			}
			trafficLightCount++
			op += fmt.Sprintf("%d\n", intersec.streets[0].endIntersection.id)
			op += fmt.Sprintf("%d\n", len(intersec.trafficLights))
			for _, tl := range intersec.trafficLights {
				if len(tl.fromStreet.cars) == 0 {
					continue
				}
				op += fmt.Sprintf("%s %d\n", tl.fromStreet.name, tl.time)
				for _, car := range tl.fromStreet.cars {
					if car.currentPos-1 == 0 {
						if len(car.paths) == 1 {
							tl.fromStreet.cars = tl.fromStreet.cars[1:]
							points += g.points + g.loopCount - i
						} else {

							car.currentPos = allStreets[car.paths[1]].duration + len(allStreets[car.paths[1]].cars)
							allStreets[car.paths[1]].cars = append(allStreets[car.paths[1]].cars, car)
							allStreets[car.paths[0]].cars = allStreets[car.paths[0]].cars[1:]
							car.paths = car.paths[1:]
						}
					} else {
						car.currentPos -= 1
					}
				}
			}
		}
	}
	return fmt.Sprintf("%d\n", trafficLightCount) + op
}
