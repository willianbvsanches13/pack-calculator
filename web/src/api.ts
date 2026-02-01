import type { CalculateResponse } from './types'

const API_BASE = '/api'

export async function fetchPackSizes(): Promise<number[]> {
  const res = await fetch(`${API_BASE}/pack-sizes`)
  if (!res.ok) throw new Error('Failed to fetch pack sizes')
  const data = await res.json()
  return data.pack_sizes || []
}

export async function addPackSize(size: number): Promise<number[]> {
  const res = await fetch(`${API_BASE}/pack-sizes/add`, {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify({ size })
  })
  const data = await res.json()
  if (!res.ok) throw new Error(data.message || data.error)
  return data.pack_sizes
}

export async function removePackSize(size: number): Promise<number[]> {
  const res = await fetch(`${API_BASE}/pack-sizes/remove`, {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify({ size })
  })
  const data = await res.json()
  if (!res.ok) throw new Error(data.message || data.error)
  return data.pack_sizes
}

export async function calculatePacks(amount: number): Promise<CalculateResponse> {
  const res = await fetch(`${API_BASE}/calculate`, {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify({ amount })
  })
  const data = await res.json()
  if (!res.ok) throw new Error(data.message || data.error)
  return data
}
