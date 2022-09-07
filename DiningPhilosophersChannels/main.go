package main

import (
	"fmt"
	"sync"
	"time"
)

type philosopher struct {
	id int

	timeToStarve time.Duration
	mealsNeeded  int
	mealsDone    int

	rightHand chan bool
	leftHand  chan bool
}

var wg sync.WaitGroup

var fork12 = make(chan bool, 1)
var fork23 = make(chan bool, 1)
var fork34 = make(chan bool, 1)
var fork45 = make(chan bool, 1)
var fork51 = make(chan bool, 1)

const lifeDuration = 10 * time.Second
const timeHungry = 9 * time.Second
const timeVeryHungry = 5 * time.Second
const timeToReleaseFork = 2 * time.Second // prevents starvation (deadlock)
const timeZero = 0 * time.Millisecond
const eatDuration = 4 * time.Second
const thinkDuration = 1 * time.Second

var philosopher1 = philosopher{1, lifeDuration, 10, 0, fork12, fork51}
var philosopher2 = philosopher{2, lifeDuration, 10, 0, fork23, fork12}
var philosopher3 = philosopher{3, lifeDuration, 10, 0, fork34, fork23}
var philosopher4 = philosopher{4, lifeDuration, 10, 0, fork45, fork34}
var philosopher5 = philosopher{5, lifeDuration, 10, 0, fork51, fork45}

func (p *philosopher) thinkAndEat() {

	for {
		start := time.Now()
		if p.mealsDone == p.mealsNeeded {
			fmt.Printf("PHILOSOPHER #%d FINISHED!\n", p.id)
			wg.Done()
			return
		} else {
			if p.timeToStarve < timeHungry {
				// Philosopher wants to eat
				// time.Sleep(500 * time.Millisecond) // Prevents starvation (optional)

				select {
				// CASE HE PICKS RIGHT FIRST
				case p.rightHand <- true: // pick fork
					fmt.Printf("Philosopher #%d picked right fork.\n", p.id)

					hasEaten, died := p.tryPick(p.leftHand)
					<-p.rightHand // release fork
					if died {
						wg.Done()
						return
					}
					if hasEaten {
						start = time.Now()
					}

				// CASE HE PICKS LEFT FIRST
				case p.leftHand <- true: // pick fork
					fmt.Printf("Philosopher #%d picked left fork.\n", p.id)

					hasEaten, died := p.tryPick(p.rightHand)
					<-p.leftHand // release fork
					if died {
						wg.Done()
						return
					}
					if hasEaten {
						start = time.Now()
					}
				}
			} else {
				// Philosopher wants to think
				p.think()
			}
		}
		p.timeToStarve = p.timeToStarve - time.Since(start)
	}
}

func (p *philosopher) tryPick(fork chan bool) (bool, bool) {
	wait := timeToReleaseFork
	hasEaten := false
	hasAlerted := false // priority
	died := false
	for (wait > timeZero || p.timeToStarve < timeVeryHungry) && !hasEaten && !died {
		tryPick := time.Now()
		if (p.timeToStarve - (wait - timeToReleaseFork)) <= timeZero {
			died = true
			fmt.Printf("PHILOSOPHER #%d DIED OF STARVATION!!!\n", p.id)
		} else if (p.timeToStarve-(wait-timeToReleaseFork)) < timeVeryHungry && !hasAlerted {
			fmt.Printf("Philosopher #%d is very hungry! He won't release the fork!\n", p.id)
			hasAlerted = true
		}
		select {
		case fork <- true: // pick fork
			if p.rightHand == fork {
				fmt.Printf("Philosopher #%d picked right fork.\n", p.id)
			} else {
				fmt.Printf("Philosopher #%d picked left fork.\n", p.id)
			}
			p.eat()
			hasEaten = true
			<-fork // release fork
		default:
			wait = wait - time.Since(tryPick)
		}
	}

	return hasEaten, died
}
func (p *philosopher) eat() {
	fmt.Printf("Philosopher #%d is eating...\n", p.id)
	time.Sleep(eatDuration)
	p.mealsDone = p.mealsDone + 1
	p.timeToStarve = lifeDuration
	fmt.Printf("Philosopher #%d is satisfied (meal %d/%d).\n", p.id, p.mealsDone, p.mealsNeeded)
}
func (p *philosopher) think() {
	fmt.Printf("Philosopher #%d is thinking...\n", p.id)
	time.Sleep(thinkDuration)
}

func main() {
	fmt.Println("Philosophers are awake!")

	wg.Add(5)

	go philosopher1.thinkAndEat()
	go philosopher2.thinkAndEat()
	go philosopher3.thinkAndEat()
	go philosopher4.thinkAndEat()
	go philosopher5.thinkAndEat()

	wg.Wait()

	fmt.Println("Philosophers went to sleep.")
}
