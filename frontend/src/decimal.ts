const decimalPattern = /^\d+(?:\.\d{1,2})?$/

function parseMinorUnits(value: string): bigint | null {
  if (!decimalPattern.test(value)) {
    return null
  }

  const [whole, fraction = ""] = value.split(".")
  return BigInt(whole) * 100n + BigInt(fraction.padEnd(2, "0"))
}

export function parsePercentage(value: string): number | null {
  const basisPoints = parseMinorUnits(value)
  if (basisPoints === null || basisPoints === 0n || basisPoints > 10_000n) {
    return null
  }
  return Number(basisPoints)
}

export function parseAmount(value: string): bigint | null {
  return parseMinorUnits(value)
}

export function formatPercentage(basisPoints: number): string {
  const whole = Math.floor(basisPoints / 100)
  const fraction = basisPoints % 100
  return `${whole}.${fraction.toString().padStart(2, "0")}`
}

export function formatAmount(cents: bigint): string {
  const whole = cents / 100n
  const fraction = cents % 100n
  return `${whole}.${fraction.toString().padStart(2, "0")}`
}
