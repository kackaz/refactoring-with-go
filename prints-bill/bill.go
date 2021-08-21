package main

import (
	"fmt"
	"math"
)

type Play struct {
	Name string
	Type string
}

type Plays map[string]Play

type Performance struct {
	PlayID   string `json:"playID"`
	Audience int    `json:"audience"`
}

type Invoice struct {
	Customer     string        `json:"customer"`
	Performances []Performance `json:"performances"`
}

func amountFor(play Play, perf Performance) float64 {
	result := 0.0

	switch play.Type {
	case "tragedy":
		result = 40000
		if perf.Audience > 30 {
			result += 1000 * (float64(perf.Audience - 30))
		}
	case "comedy":
		result = 30000
		if perf.Audience > 20 {
			result += 10000 + 500*(float64(perf.Audience-20))
		}
		result += 300 * float64(perf.Audience)
	default:
		panic(fmt.Sprintf("unknow type: %s", play.Type))
	}
	return result
}

func volumeCreditsFor(play Play, perf Performance) float64 {
	volumeCredits := math.Max(float64(perf.Audience-30), 0)
	if "comedy" == play.Type {
		volumeCredits += math.Floor(float64(perf.Audience / 5))
	}
	return volumeCredits
}

func totalAmountFor(performances []Performance, plays Plays) float64 {
	totalAmount := 0.0
	for _, perf := range performances {
		play := plays[perf.PlayID]
		totalAmount += amountFor(play, perf)
	}
	return totalAmount
}

func statement(invoice Invoice, plays Plays) string {
	totalAmount := totalAmountFor(invoice.Performances, plays)
	volumeCredits := 0.0
	result := fmt.Sprintf("Statement for %s\n", invoice.Customer)

	// for _, perf := range invoice.Performances {
	// 	play := plays[perf.PlayID]
	// 	totalAmount += amountFor(play, perf)
	// }

	for _, perf := range invoice.Performances {
		play := plays[perf.PlayID]
		volumeCredits += volumeCreditsFor(play, perf)
	}

	for _, perf := range invoice.Performances {
		play := plays[perf.PlayID]
		result += fmt.Sprintf("  %s: $%.2f (%d seats)\n", play.Name, amountFor(play, perf)/100, perf.Audience)
	}

	result += fmt.Sprintf("Amount owed is $%.2f\n", totalAmount/100)
	result += fmt.Sprintf("you earned %.0f credits\n", volumeCredits)
	return result
}

func main() {
	inv := Invoice{
		Customer: "Bigco",
		Performances: []Performance{
			{PlayID: "hamlet", Audience: 55},
			{PlayID: "as-like", Audience: 35},
			{PlayID: "othello", Audience: 40},
		}}
	plays := Plays{
		"hamlet":  {Name: "Hamlet", Type: "tragedy"},
		"as-like": {Name: "As You Like It", Type: "comedy"},
		"othello": {Name: "Othello", Type: "tragedy"},
	}

	bill := statement(inv, plays)
	fmt.Println(bill)
}
