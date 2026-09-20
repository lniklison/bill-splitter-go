package bills

import (
	"context"
	"fmt"
	"strings"
)

type Repository interface {
	GetBill(ctx context.Context, billID int64) (Bill, error)
	ReplaceShares(ctx context.Context, billID int64, shares []Share) (Bill, error)
}

type ShareInput struct {
	PersonName string
	Percentage string
}

type Service struct {
	repository Repository
}

func NewService(repository Repository) *Service {
	return &Service{repository: repository}
}

func (s *Service) GetShares(ctx context.Context, billID int64) (Allocation, error) {
	bill, err := s.repository.GetBill(ctx, billID)
	if err != nil {
		return Allocation{}, err
	}
	return allocationFor(bill)
}

func (s *Service) ReplaceShares(ctx context.Context, billID int64, inputs []ShareInput) (Allocation, error) {
	shares, err := validateShares(inputs)
	if err != nil {
		return Allocation{}, err
	}

	bill, err := s.repository.ReplaceShares(ctx, billID, shares)
	if err != nil {
		return Allocation{}, err
	}
	return allocationFor(bill)
}

func validateShares(inputs []ShareInput) ([]Share, error) {
	if len(inputs) == 0 {
		return nil, validationError("shares_required", "At least one share is required")
	}

	shares := make([]Share, 0, len(inputs))
	seenNames := make(map[string]struct{}, len(inputs))
	total := Percentage(0)
	for position, input := range inputs {
		personName := strings.TrimSpace(input.PersonName)
		if personName == "" {
			return nil, validationError("invalid_person_name", "Person name must not be blank")
		}
		normalizedName := strings.ToLower(personName)
		if _, exists := seenNames[normalizedName]; exists {
			return nil, validationError("duplicate_person_name", "Person names must be unique")
		}

		percentage, err := ParsePercentage(input.Percentage)
		if err != nil {
			return nil, err
		}
		seenNames[normalizedName] = struct{}{}
		total += percentage
		shares = append(shares, Share{
			Position:       position,
			PersonName:     personName,
			NormalizedName: normalizedName,
			Percentage:     percentage,
		})
	}
	if total != 10_000 {
		return nil, validationError("invalid_percentage_total", "Percentages must total exactly 100.00")
	}
	return shares, nil
}

func allocationFor(bill Bill) (Allocation, error) {
	if len(bill.Shares) == 0 {
		return Allocation{}, fmt.Errorf("bill %d has no shares", bill.ID)
	}
	totalPercentage := Percentage(0)
	for _, share := range bill.Shares {
		totalPercentage += share.Percentage
	}
	if totalPercentage != 10_000 {
		return Allocation{}, fmt.Errorf("bill %d percentages total %s", bill.ID, totalPercentage.String())
	}

	return Allocation{
		BillID:      bill.ID,
		Description: bill.Description,
		TotalAmount: bill.TotalAmount,
		Shares:      allocate(bill.TotalAmount, bill.Shares),
	}, nil
}
