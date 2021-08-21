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

func (play Play) amountFor(perf Performance) float64 {
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

func (play Play) volumeCreditsFor(perf Performance) float64 {
	result := math.Max(float64(perf.Audience-30), 0)
	if "comedy" == play.Type {
		result += math.Floor(float64(perf.Audience / 5))
	}
	return result
}

type Rate struct {
	Amount       float64
	VolumeCredit float64
	Audience     int
	Name         string
}

type Bill struct {
	Customer           string
	Rates              []Rate
	TotalAmount        float64
	TotalVolumeCredits float64
}

func (bill Bill) renderText() string {
	result := fmt.Sprintf("Statement for %s\n", bill.Customer)
	for _, r := range bill.Rates {
		result += fmt.Sprintf("  %s: $%.2f (%d seats)\n", r.Name, r.Amount/100, r.Audience)
	}

	result += fmt.Sprintf("Amount owed is $%.2f\n", bill.TotalAmount/100)
	result += fmt.Sprintf("you earned %.0f credits\n", bill.TotalVolumeCredits)
	return result
}

func totalAmountFor(rates []Rate) (amounts float64) {
	for _, r := range rates {
		amounts += r.Amount
	}
	return
}

func totalVolumeCreditsFor(rates []Rate) (credit float64) {
	for _, r := range rates {
		credit += r.VolumeCredit
	}
	return
}

func statement(invoice Invoice, plays Plays) string {
	rates := []Rate{}
	for _, perf := range invoice.Performances {
		play := plays[perf.PlayID]
		rates = append(rates, Rate{
			Amount:       play.amountFor(perf),
			VolumeCredit: play.volumeCreditsFor(perf),
			Name:         play.Name,
			Audience:     perf.Audience,
		})
	}
	bill := Bill{
		Customer:           invoice.Customer,
		Rates:              rates,
		TotalAmount:        totalAmountFor(rates),
		TotalVolumeCredits: totalVolumeCreditsFor(rates),
	}

	return bill.renderText()
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
