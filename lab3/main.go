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

type QueueOfClients struct {
	clients []Client
	mu      sync.Mutex
}

func (barber *Barber) doHaircut(client Client, wg *sync.WaitGroup) {
	defer wg.Done()
	fmt.Printf("Barber started cutting client №%d\n", client.id)
	time.Sleep(2 * time.Second)
	fmt.Printf("Barber finished cutting client №%d\n", client.id)
}

func (barber *Barber) lookForClients(queueOfClients *QueueOfClients, wg *sync.WaitGroup) {
	for {
		queueOfClients.mu.Lock()
		if len(queueOfClients.clients) > 0 {
			client := queueOfClients.clients[0]
			queueOfClients.clients = queueOfClients.clients[1:]
			queueOfClients.mu.Unlock()

			barber.doHaircut(client, wg)
			barber.isSleep = false
		} else {
			fmt.Printf("Barber falls asleep\n")
			barber.isSleep = true
			queueOfClients.mu.Unlock()
		}
	}
}

func (client *Client) arriveToBarbershop(queueOfClients *QueueOfClients) {
	queueOfClients.mu.Lock()
	fmt.Printf("Client №%d arrives at the barbershop\n", client.id)
	if len(queueOfClients.clients) < 3 {
		queueOfClients.clients = append(queueOfClients.clients, *client)
		fmt.Printf("Client №%d joins the queue\n", client.id)
	} else {
		fmt.Printf("Client №%d leaves the barbershop because all chairs are taken\n", client.id)
	}
	queueOfClients.mu.Unlock()
}

func main() {
	const chairs = 3
	queueOfClients := QueueOfClients{clients: make([]Client, 0, chairs)}
	barber := Barber{isSleep: true}
	var wg sync.WaitGroup

	go barber.lookForClients(&queueOfClients, &wg)

	for i := 0; i < 10; i++ {
		client := Client{id: i}
		wg.Add(1)
		go client.arriveToBarbershop(&queueOfClients)
		time.Sleep(time.Duration(rand.Intn(5000)) * time.Millisecond)
	}

	wg.Wait()
	fmt.Println("All clients have been served or left, and the barbershop is now closed.")
}
