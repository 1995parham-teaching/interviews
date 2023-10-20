/*
Write a Delivery class (or object) that represents a delivery with a destination
and distance.
Deliveries require different sensors, depending on their distance.

Add a method, getNeededSensors, that returns a mapping of sensor name to the
count of that sensor needed to complete the delivery according to these rules:
- If distance < 10 miles, require 1 gps and 1 temp sensor.
- If 10 <= distance < 100 miles require 1 gps, 2 temp, and 1 weight sensor.
- If distance >= 100 miles require 2 gps, 4 temp, and 2 weight sensors.

Part 2:
Write a Scheduler class (or object) that represents a daily delivery scheduler with a set of available sensors.
Add a method, scheduleDeliveries, that given a list of deliveries as an argument,
returns a list of deliveries that can be made that day.
Assume all deliveries will be leaving at the same time every day so sensors can only
be used once.


Test Cases to consider:
deliveryA = Delivery("A", 9)
deliveryB = Delivery("B", 15)
deliveryC = Delivery("C", 100)
scheduler = Scheduler({"gps": 2, "temp": 4, "weight": 2})
scheduler.scheduleDeliveries([deliveryA, deliveryB, deliveryC]) == [deliveryA, deliveryB]
scheduler.scheduleDeliveries([deliveryA, deliveryC, deliveryB]) == [deliveryA, deliveryB]
scheduler.scheduleDeliveries([deliveryC, deliveryA, deliveryB]) == [deliveryC]

Part 3:
We get paid a flat fee for all deliveries. Modify the scheduleDeliveries
function to maximize the number of deliveries that will be made in a day.

Our previous test:
scheduler.scheduleDeliveries([deliveryC, deliveryA, deliveryB]) ==[deliveryC]
Should now return:
scheduler.scheduleDeliveries([deliveryC, deliveryA, deliveryB]) ==[deliveryA, deliveryB]

Part 4:
We recently purchased a new type of sensor, doorSensor.
A doorSensor can be used in place of 1 weight sensor or in place of 2 temperature
sensors at any time.
Modify our existing functions to maximize the day's deliveries with the new sensor.
*/

package main

import (
	"fmt"
	"log"
	"sort"
)

type SortBySensors []Delivery

func (a SortBySensors) Len() int      { return len(a) }
func (a SortBySensors) Swap(i, j int) { a[i], a[j] = a[j], a[i] }
func (a SortBySensors) Less(i, j int) bool {
	needed1 := a[i].GetNeededSensors()
	needed2 := a[j].GetNeededSensors()

	sum1 := 0
	sum2 := 0

	for _, sn := range SensorNames {
		sum1 += needed1[sn]
		sum2 += needed2[sn]
	}

	return sum1 < sum2
}

type SensorName string

const (
	GPS    SensorName = "gps"
	Temp   SensorName = "temp"
	Weight SensorName = "weight"
	Door   SensorName = "door"
)

var SensorNames = []SensorName{GPS, Temp, Weight, Door}

type Delivery struct {
	destination string
	distance    int
}

func NewDelivery(destination string, distance int) Delivery {
	return Delivery{
		destination: destination,
		distance:    distance,
	}
}

func (d Delivery) GetNeededSensors() map[SensorName]int {
	sensors := make(map[SensorName]int)

	if d.distance < 10 {
		sensors[GPS] = 1
		sensors[Temp] = 1
	} else if d.distance >= 10 && d.distance < 100 {
		sensors[GPS] = 1
		sensors[Temp] = 2
		sensors[Weight] = 1
	} else {
		sensors[GPS] = 2
		sensors[Temp] = 4
		sensors[Weight] = 2
	}

	return sensors
}

type Scheduler struct {
	resources map[SensorName]int
}

func NewScheduler(resources map[SensorName]int) Scheduler {
	return Scheduler{
		resources: resources,
	}
}

func (s Scheduler) ScheduleDeliveries(deliveries []Delivery) []Delivery {
	sort.Sort(SortBySensors(deliveries))

	success := make([]Delivery, 0)
	resources := make(map[SensorName]int)

	for _, sn := range SensorNames {
		resources[sn] = s.resources[sn]
	}

	for _, delivery := range deliveries {
		failed := false

		needed := delivery.GetNeededSensors()
		for _, sn := range SensorNames {
			if needed[sn] > resources[sn] {
				if sn == Weight {
					if (needed[sn] - resources[sn]) <= resources[Door] {
						resources[Door] -= (needed[sn] - resources[sn])
						resources[Weight] += (needed[sn] - resources[sn])
						continue
					}
				} else if sn == Temp {
					if (needed[sn] - resources[sn]) <= 2*resources[Door] {
						resources[Door] -= (needed[sn] - resources[sn]) / 2
						resources[Temp] += (needed[sn] - resources[sn]) / 2
						continue
					}
				}
				failed = true
				break
			}
		}

		if !failed {
			for _, sn := range SensorNames {
				resources[sn] -= needed[sn]
			}
			success = append(success, delivery)
		}
	}

	return success
}

func main() {
	{
		d1 := NewDelivery("NY", 7)
		d2 := NewDelivery("FL", 10)
		d3 := NewDelivery("CA", 158)

		m1 := d1.GetNeededSensors()

		if m1[GPS] != 1 {
			log.Fatal("GetNeededSensors is not working as expected")
		}
		if m1[Temp] != 1 {
			log.Fatal("GetNeededSensors is not working as expected")
		}

		fmt.Println(d2.GetNeededSensors())
		fmt.Println(d3.GetNeededSensors())
	}

	fmt.Println()

	{
		deliveryA := NewDelivery("A", 9)
		deliveryB := NewDelivery("B", 15)
		deliveryC := NewDelivery("C", 100)
		scheduler := NewScheduler(map[SensorName]int{GPS: 2, Temp: 4, Weight: 2})

		fmt.Println(scheduler.ScheduleDeliveries([]Delivery{deliveryA, deliveryB, deliveryC}))
		fmt.Println(scheduler.ScheduleDeliveries([]Delivery{deliveryA, deliveryC, deliveryB}))
		fmt.Println(scheduler.ScheduleDeliveries([]Delivery{deliveryC, deliveryA, deliveryB}))
	}

	fmt.Println()

	{
		deliveryA := NewDelivery("A", 9)
		deliveryB := NewDelivery("B", 15)
		deliveryC := NewDelivery("C", 100)
		scheduler := NewScheduler(map[SensorName]int{GPS: 0, Temp: 0, Weight: 0})

		fmt.Println(scheduler.ScheduleDeliveries([]Delivery{deliveryA, deliveryB, deliveryC}))
		fmt.Println(scheduler.ScheduleDeliveries([]Delivery{deliveryA, deliveryC, deliveryB}))
		fmt.Println(scheduler.ScheduleDeliveries([]Delivery{deliveryC, deliveryA, deliveryB}))
	}

	fmt.Println()

	{
		deliveryA := NewDelivery("A", 9)
		deliveryB := NewDelivery("B", 15)
		deliveryC := NewDelivery("C", 100)
		scheduler := NewScheduler(map[SensorName]int{GPS: 4, Temp: 4, Weight: 2, Door: 2})

		fmt.Println(scheduler.ScheduleDeliveries([]Delivery{deliveryA, deliveryB, deliveryC}))
		fmt.Println(scheduler.ScheduleDeliveries([]Delivery{deliveryA, deliveryC, deliveryB}))
		fmt.Println(scheduler.ScheduleDeliveries([]Delivery{deliveryC, deliveryA, deliveryB}))
	}
}
