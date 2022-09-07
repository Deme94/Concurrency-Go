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

	rightHand *sync.Mutex
	leftHand  *sync.Mutex
}

var wg sync.WaitGroup

var fork12 sync.Mutex
var fork23 sync.Mutex
var fork34 sync.Mutex
var fork45 sync.Mutex
var fork51 sync.Mutex

const lifeDuration = 18 * time.Second
const timeHungry = 17 * time.Second
const timeVeryHungry = 6 * time.Second
const timeToReleaseFork = 2 * time.Second // prevents starvation (deadlock)
const timeZero = 0 * time.Millisecond
const eatDuration = 4 * time.Second
const thinkDuration = 1 * time.Second

var philosopher1 = philosopher{1, lifeDuration, 10, 0, &fork12, &fork51}
var philosopher2 = philosopher{2, lifeDuration, 10, 0, &fork23, &fork12}
var philosopher3 = philosopher{3, lifeDuration, 10, 0, &fork34, &fork23}
var philosopher4 = philosopher{4, lifeDuration, 10, 0, &fork45, &fork34}
var philosopher5 = philosopher{5, lifeDuration, 10, 0, &fork51, &fork45}

func (p *philosopher) thinkAndEat() {
	start := time.Now()
	for {
		if p.mealsDone == p.mealsNeeded {
			fmt.Printf("PHILOSOPHER #%d FINISHED!\n", p.id)
			wg.Done()
			return
		} else {
			if p.timeToStarve < timeHungry {
				// Philosopher wants to eat

				// Mutual exclusion ---
				p.rightHand.Lock() // pick fork
				fmt.Printf("Philosopher #%d picked right fork.\n", p.id)
				wait := timeToReleaseFork
				hasEaten := false
				hasAlerted := false // priority
				for (wait > timeZero || p.timeToStarve < timeVeryHungry) && !hasEaten {
					tryToPick := time.Now()

					if (p.timeToStarve - (wait - timeToReleaseFork)) <= timeZero {
						fmt.Printf("PHILOSOPHER #%d DIED OF STARVATION!!!\n", p.id)
						p.rightHand.Unlock()
						wg.Done()
						return
					} else if (p.timeToStarve-(wait-timeToReleaseFork)) < timeVeryHungry && !hasAlerted {
						fmt.Printf("Philosopher #%d is very hungry! He won't release the fork!\n", p.id)
						hasAlerted = true
					}

					if p.leftHand.TryLock() { // pick fork - TryLock is not a recommended function but it's necessary to prevent starvation and release the forks after timeout
						fmt.Printf("Philosopher #%d picked left fork.\n", p.id)
						p.eat()
						hasEaten = true
						start = time.Now()
						p.leftHand.Unlock() // release fork
					} else {
						wait = wait - time.Since(tryToPick)
					}
					// --------------------
				}
				p.rightHand.Unlock() // release fork
			} else {
				// Philosopher wants to think
				p.think()
			}
		}
		p.timeToStarve = p.timeToStarve - time.Since(start)
	}
}

func (p *philosopher) think() {
	fmt.Printf("Philosopher #%d is thinking...\n", p.id)
	time.Sleep(thinkDuration)
}
func (p *philosopher) eat() {
	fmt.Printf("Philosopher #%d is eating...\n", p.id)
	time.Sleep(eatDuration)
	p.mealsDone = p.mealsDone + 1
	p.timeToStarve = lifeDuration
	fmt.Printf("Philosopher #%d is satisfied (meal %d/%d).\n", p.id, p.mealsDone, p.mealsNeeded)
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
