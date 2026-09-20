import type { BillShares, ShareInput } from "../types"

const genericErrorMessage = "Something went wrong. Please try again."

interface ErrorPayload {
  error?: {
    code?: unknown
    message?: unknown
  }
}

export class BillsApiError extends Error {
  constructor(
    message: string,
    readonly code = "unexpected_error",
  ) {
    super(message)
    this.name = "BillsApiError"
  }
}

async function request(path: string, init?: RequestInit): Promise<BillShares> {
  let response: Response
  try {
    response = await fetch(path, init)
  } catch {
    throw new BillsApiError(genericErrorMessage)
  }

  if (!response.ok) {
    throw await errorFor(response)
  }

  try {
    const payload: unknown = await response.json()
    if (isBillShares(payload)) {
      return payload
    }
    throw new BillsApiError(genericErrorMessage)
  } catch {
    throw new BillsApiError(genericErrorMessage)
  }
}

function isBillShares(value: unknown): value is BillShares {
  if (typeof value !== "object" || value === null) {
    return false
  }

  const candidate = value as Partial<BillShares>
  return (
    typeof candidate.bill?.id === "number" &&
    typeof candidate.bill.description === "string" &&
    typeof candidate.bill.totalAmount === "string" &&
    Array.isArray(candidate.shares) &&
    candidate.shares.every(
      (share) =>
        typeof share?.personName === "string" &&
        typeof share.percentage === "string" &&
        typeof share.amount === "string",
    )
  )
}

async function errorFor(response: Response): Promise<BillsApiError> {
  try {
    const payload = (await response.json()) as ErrorPayload
    if (
      typeof payload.error?.code === "string" &&
      typeof payload.error.message === "string"
    ) {
      return new BillsApiError(payload.error.message, payload.error.code)
    }
  } catch {
    // A non-JSON error response is handled by the generic fallback below.
  }
  return new BillsApiError(genericErrorMessage)
}

export function getBillShares(billID: number): Promise<BillShares> {
  return request(`/bills/${billID}/shares`)
}

export function replaceBillShares(
  billID: number,
  shares: ShareInput[],
): Promise<BillShares> {
  return request(`/bills/${billID}/shares`, {
    method: "PUT",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify({ shares }),
  })
}
