package bills

import "sort"

func allocate(total Cents, shares []Share) []AllocatedShare {
	type remainder struct {
		position int
		value    int64
	}

	allocated := make([]AllocatedShare, len(shares))
	remainders := make([]remainder, len(shares))
	var allocatedTotal Cents
	for index, share := range shares {
		productWhole := int64(total/10_000) * int64(share.Percentage)
		productRemainder := int64(total%10_000) * int64(share.Percentage)
		amount := Cents(productWhole + productRemainder/10_000)
		allocated[index] = AllocatedShare{
			PersonName: share.PersonName,
			Percentage: share.Percentage,
			Amount:     amount,
		}
		allocatedTotal += amount
		remainders[index] = remainder{position: index, value: productRemainder % 10_000}
	}

	sort.Slice(remainders, func(i, j int) bool {
		if remainders[i].value == remainders[j].value {
			return remainders[i].position < remainders[j].position
		}
		return remainders[i].value > remainders[j].value
	})
	for index := Cents(0); index < total-allocatedTotal; index++ {
		allocated[remainders[index].position].Amount++
	}

	return allocated
}
