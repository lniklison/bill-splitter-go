import { formatAmount, parseAmount, parsePercentage } from "./decimal"
import type { DraftShare } from "./types"

interface Remainder {
  position: number
  value: bigint
}

export function allocateDraft(
  totalAmount: string,
  shares: DraftShare[],
): Array<string | null> {
  const totalCents = parseAmount(totalAmount)
  if (totalCents === null) {
    return shares.map(() => null)
  }

  const percentages = shares.map((share) => parsePercentage(share.percentage))
  const amounts: Array<bigint | null> = []
  const remainders: Remainder[] = []
  let allocatedTotal = 0n

  percentages.forEach((percentage, position) => {
    if (percentage === null) {
      amounts.push(null)
      return
    }

    const product = totalCents * BigInt(percentage)
    const amount = product / 10_000n
    amounts.push(amount)
    allocatedTotal += amount
    remainders.push({ position, value: product % 10_000n })
  })

  const percentageTotal = percentages.reduce<number>(
    (total, percentage) => total + (percentage ?? 0),
    0,
  )
  if (percentageTotal === 10_000 && remainders.length === shares.length) {
    remainders.sort((left, right) => {
      if (left.value === right.value) {
        return left.position - right.position
      }
      return left.value > right.value ? -1 : 1
    })
    const remainingCents = Number(totalCents - allocatedTotal)
    for (let index = 0; index < remainingCents; index += 1) {
      const position = remainders[index].position
      amounts[position] = (amounts[position] ?? 0n) + 1n
    }
  }

  return amounts.map((amount) => (amount === null ? null : formatAmount(amount)))
}
