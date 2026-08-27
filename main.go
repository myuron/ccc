package main

import (
	"encoding/csv"
	"fmt"
	"log"
	"os"

	"github.com/olekukonko/tablewriter"
)

const (
	JpyPerUsd      = 159.32
	CcTeamStandard = 22
	CcTeamPremium  = 110
	StandardPlan   = "standard"
	PremiumPlan    = "premium"
	Disable        = "disable"
)

func main() {
	members_file := os.Args[1]
	if err := run(members_file); err != nil {
		log.Fatal(err)
	}
}

func run(members_file string) error {
	csv_file, err := os.Open(members_file)
	if err != nil {
		return err
	}
	defer func() {
		err := csv_file.Close()
		if err != nil {
			log.Print(err)
		}
	}()

	csv_data := csv.NewReader(csv_file)
	csv_row_data, err := csv_data.ReadAll()
	if err != nil {
		return err
	}

	var standard_count int
	var premium_count int
	table := tablewriter.NewWriter(os.Stdout)
	members := [][]string{
		{"Name", "Standard", "Premium"},
	}
	for i, v := range csv_row_data {
		if i == 0 {
			continue
		}
		if len(v) != 2 {
			return fmt.Errorf("the data in the %v row is missing", i+1)
		}
		switch v[1] {
		case StandardPlan:
			standard_count++
			if err := table.Append([]string{v[0], "○", ""}); err != nil {
				return err
			}
		case PremiumPlan:
			premium_count++
			if err := table.Append([]string{v[0], "", "○"}); err != nil {
				return err
			}
		case Disable:
			continue
		default:
			return fmt.Errorf("unknown plan %q on the %v row", v[1], i+1)
		}
	}

	standard_cost := float64(standard_count) * float64(CcTeamStandard) * JpyPerUsd
	premium_cost := float64(premium_count) * float64(CcTeamPremium) * JpyPerUsd
	total_cost := standard_cost + premium_cost

	table.Header(members[0])
	if err := table.Render(); err != nil {
		return err
	}

	fmt.Printf("standard cost: %.2fJPY\n", standard_cost)
	fmt.Printf("premium cost: %.2fJPY\n", premium_cost)
	fmt.Printf("total cost: %.2fJPY\n", total_cost)

	return nil
}
