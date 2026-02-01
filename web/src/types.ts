export interface PackSizesResponse {
  pack_sizes: number[]
  message?: string
}

export interface CalculateRequest {
  amount: number
  pack_sizes?: number[]
}

export interface CalculateResponse {
  order_amount: number
  total_items: number
  total_packs: number
  packs: Record<number, number>
  pack_sizes_used: number[]
}

export interface ErrorResponse {
  error: string
  message?: string
}
