export interface Bill {
  id: number
  description: string
  totalAmount: string
}

export interface Share {
  personName: string
  percentage: string
  amount: string
}

export interface BillShares {
  bill: Bill
  shares: Share[]
}

export interface ShareInput {
  personName: string
  percentage: string
}

export interface DraftShare extends ShareInput {
  key: number
}

export interface RowErrors {
  personName?: string
  percentage?: string
}
