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
