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
	result := math.Max(float64(perf.Audience-30), 0)
	if "comedy" == play.Type {
		result += math.Floor(float64(perf.Audience / 5))
	}
	return result
}

func totalAmountFor(performances []Performance, plays Plays) float64 {
	result := 0.0
	for _, perf := range performances {
		play := plays[perf.PlayID]
		result += amountFor(play, perf)
	}
	return result
}

func totalVolumeCreditsFor(performances []Performance, plays Plays) float64 {
	result := 0.0
	for _, perf := range performances {
		play := plays[perf.PlayID]
		result += volumeCreditsFor(play, perf)
	}
	return result
}

type Rate struct {
	Name     string
	Amount   float64
	Audience int
}

type Bill struct {
	Customer           string
	TotalAmount        float64
	TotalVolumeCredits float64
	Rates              []Rate
}

func ratesFor(performances []Performance, plays Plays) []Rate {
	var result []Rate
	for _, perf := range performances {
		play := plays[perf.PlayID]
		rate := Rate{
			Name:     play.Name,
			Amount:   amountFor(play, perf),
			Audience: perf.Audience,
		}
		result = append(result, rate)
	}
	return result
}

func renderText(bill Bill) string {
	result := fmt.Sprintf("Statement for %s\n", bill.Customer)
	for _, r := range bill.Rates {
		result += fmt.Sprintf("  %s: $%.2f (%d seats)\n", r.Name, r.Amount/100, r.Audience)
	}

	result += fmt.Sprintf("Amount owed is $%.2f\n", bill.TotalAmount/100)
	result += fmt.Sprintf("you earned %.0f credits\n", bill.TotalVolumeCredits)
	return result
}

func statement(invoice Invoice, plays Plays) string {

	bill := Bill{
		Customer:           invoice.Customer,
		TotalAmount:        totalAmountFor(invoice.Performances, plays),
		TotalVolumeCredits: totalVolumeCreditsFor(invoice.Performances, plays),
		Rates:              ratesFor(invoice.Performances, plays),
	}

	return renderText(bill)
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
