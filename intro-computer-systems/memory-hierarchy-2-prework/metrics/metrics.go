package metrics

import (
	"encoding/csv"
	"log"
	"math"
	"os"
	"strconv"
)

type Payment struct {
	cents uint32
}

type User struct {
	age uint8
}

func AverageAge(users []*User) float64 {
	var sum uint64
	count := 0
	for _, u := range users {
		count++
		sum += uint64(u.age)
	}
	return float64(sum) / float64(count)
}

func AveragePaymentAmount(payments []*Payment) float64 {
	var sum uint64
	count := 0
	for _, p := range payments {
		count++
		sum += uint64(p.cents)
	}
	return 0.01 * float64(sum) / float64(count)
}

// Compute the standard deviation of payment amounts
func StdDevPaymentAmount(payments []*Payment) float64 {
	mean := AveragePaymentAmount(payments) * 100
	squaredDiffs := 0.0
	count := 0
	for _, p := range payments {
		count += 1
		diff := float64(p.cents) - mean
		squaredDiffs += diff * diff
	}
	return math.Sqrt(squaredDiffs/float64(count)) / 100
}

func LoadData() ([]*User, []*Payment) {
	f, err := os.Open("users.csv")
	if err != nil {
		log.Fatalln("Unable to read users.csv", err)
	}
	reader := csv.NewReader(f)
	userLines, err := reader.ReadAll()
	if err != nil {
		log.Fatalln("Unable to parse users.csv as csv", err)
	}

	var users []*User
	for _, line := range userLines {
		age, _ := strconv.Atoi(line[2])
		users = append(users, &User{uint8(age)})
	}

	f, err = os.Open("payments.csv")
	if err != nil {
		log.Fatalln("Unable to read payments.csv", err)
	}
	reader = csv.NewReader(f)
	paymentLines, err := reader.ReadAll()
	if err != nil {
		log.Fatalln("Unable to parse payments.csv as csv", err)
	}

	var payments []*Payment
	for _, line := range paymentLines {
		paymentCents, _ := strconv.Atoi(line[0])
		payments = append(payments, &Payment{
			uint32(paymentCents),
		})
	}

	return users, payments
}
