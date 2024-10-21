package main

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

type Barber struct {
	isSleep bool
}

type Client struct {
	id int
}

func (barber *Barber) doHaircut(client Client, queueOfClients chan Client, wg *sync.WaitGroup) {
	fmt.Printf("Barber started cutting client №%d\n", client.id)
	time.Sleep(2 * time.Second)
	fmt.Printf("Barber finished cutting client №%d\n", client.id)
	barber.lookForClients(queueOfClients, wg)
}

func (barber *Barber) lookForClients(queueOfClients chan Client, wg *sync.WaitGroup) {
	select {
	case client := <-queueOfClients:
		barber.doHaircut(client, queueOfClients, wg)
	default:
		fmt.Printf("Barber falls asleep\n")
		barber.isSleep = true
	}
}

func (client *Client) arriveToBarbershop(barber *Barber, queueOfClients chan Client, wg *sync.WaitGroup) {
	defer wg.Done()
	fmt.Printf("Client №%d arrives at the barbershop\n", client.id)
	if barber.isSleep {
		fmt.Printf("Barber wakes up\n")
		barber.isSleep = false
		barber.doHaircut(*client, queueOfClients, wg)
	} else {
		select {
		case queueOfClients <- *client:
			fmt.Printf("Client №%d joins the queue\n", client.id)
		default:
			fmt.Printf("Client №%d leaves the barbershop because all chairs are taken\n", client.id)
		}
	}
}

func main() {
	const chairs = 3
	queueOfClients := make(chan Client, chairs)
	barber := Barber{isSleep: true}
	var wg sync.WaitGroup

	for i := 0; i < 10; i++ {
		client := Client{id: i}
		wg.Add(1)
		go client.arriveToBarbershop(&barber, queueOfClients, &wg)
		time.Sleep(time.Duration(rand.Intn(2000)) * time.Millisecond)
	}

	wg.Wait()
	fmt.Println("All clients have been served or left, and the barbershop is now closed.\n")
}
