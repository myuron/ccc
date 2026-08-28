package main

import (
	"encoding/csv"
	"fmt"
	"io"
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
	if len(os.Args) < 2 {
		log.Fatal("usage: ccc <members.csv>")
	}
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

	csv_row_data, err := ParseMembers(csv_file)
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
		switch v.Plan {
		case StandardPlan:
			standard_count++
			if err := table.Append([]string{v.Name, "○", ""}); err != nil {
				return err
			}
		case PremiumPlan:
			premium_count++
			if err := table.Append([]string{v.Name, "", "○"}); err != nil {
				return err
			}
		case Disable:
			continue
		default:
			return fmt.Errorf("unknown plan %q on the %v row", v.Plan, i+1)
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

type Member struct {
	Name string
	Plan string
}

func ParseMembers(r io.Reader) ([]Member, error) {
	rows, err := csv.NewReader(r).ReadAll()
	if err != nil {
		return nil, err
	}
	var members []Member
	for i, v := range rows {
		if i == 0 {
			continue
		}
		if len(v) != 2 {
			return nil, fmt.Errorf("the data in the %v row is missing", i+1)
		}
		members = append(members, Member{Name: v[0], Plan: v[1]})

	}
	return members, nil
}
