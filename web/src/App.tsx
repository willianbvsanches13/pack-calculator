import { useState } from 'react'
import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query'
import { fetchPackSizes, addPackSize, removePackSize, calculatePacks } from './api'
import type { CalculateResponse } from './types'

interface PackDisplay {
  size: number
  qty: number
}

function App() {
  const queryClient = useQueryClient()
  const [newSize, setNewSize] = useState('')
  const [amount, setAmount] = useState('')
  const [result, setResult] = useState<CalculateResponse | null>(null)
  const [error, setError] = useState('')

  const { data: packSizes = [], isLoading: loadingPackSizes } = useQuery({
    queryKey: ['packSizes'],
    queryFn: fetchPackSizes,
  })

  const addMutation = useMutation({
    mutationFn: addPackSize,
    onSuccess: (data) => {
      queryClient.setQueryData(['packSizes'], data)
      setNewSize('')
    },
    onError: (err: Error) => alert(err.message),
  })

  const removeMutation = useMutation({
    mutationFn: removePackSize,
    onSuccess: (data) => {
      queryClient.setQueryData(['packSizes'], data)
    },
    onError: (err: Error) => alert(err.message),
  })

  const calculateMutation = useMutation({
    mutationFn: calculatePacks,
    onSuccess: (data) => {
      setResult(data)
      setError('')
    },
    onError: (err: Error) => {
      setError(err.message)
      setResult(null)
    },
  })

  function handleAddPackSize(e: React.FormEvent) {
    e.preventDefault()
    const size = parseInt(newSize)
    if (!size || size <= 0) return
    addMutation.mutate(size)
  }

  function handleRemovePackSize(size: number) {
    if (!confirm(`Remove ${size}?`)) return
    removeMutation.mutate(size)
  }

  function handleCalculate(e: React.FormEvent) {
    e.preventDefault()
    const amt = parseInt(amount)
    if (!amt || amt <= 0) {
      setError('Enter a valid amount')
      return
    }
    calculateMutation.mutate(amt)
  }

  const sortedPacks: PackDisplay[] = result
    ? Object.entries(result.packs)
        .map(([size, qty]) => ({ size: parseInt(size), qty }))
        .sort((a, b) => b.size - a.size)
    : []

  return (
    <div className="min-h-screen bg-gradient-to-br from-indigo-500 to-purple-600 p-5">
      <div className="max-w-3xl mx-auto">
        <h1 className="text-4xl font-bold text-white text-center mb-8">
          Pack Calculator
        </h1>

        {/* Pack Sizes */}
        <div className="bg-white rounded-xl p-6 mb-5 shadow-xl">
          <h2 className="text-xl font-semibold text-gray-800 border-b-2 border-indigo-500 pb-2 mb-4">
            Pack Sizes
          </h2>

          <div className="flex flex-wrap gap-2 mb-4">
            {loadingPackSizes ? (
              <span className="text-gray-500">Loading...</span>
            ) : (
              [...packSizes].sort((a, b) => a - b).map(size => (
                <span
                  key={size}
                  className="bg-indigo-500 text-white px-4 py-2 rounded-full flex items-center gap-2 font-medium"
                >
                  {size.toLocaleString()}
                  <button
                    onClick={() => handleRemovePackSize(size)}
                    disabled={removeMutation.isPending}
                    className="bg-white/30 hover:bg-white/50 disabled:opacity-50 w-5 h-5 rounded-full text-sm"
                  >
                    ×
                  </button>
                </span>
              ))
            )}
          </div>

          <form onSubmit={handleAddPackSize} className="flex gap-3">
            <input
              type="number"
              value={newSize}
              onChange={e => setNewSize(e.target.value)}
              placeholder="New pack size"
              className="flex-1 px-4 py-3 border-2 border-gray-200 rounded-lg focus:border-indigo-500 focus:outline-none"
            />
            <button
              type="submit"
              disabled={addMutation.isPending}
              className="bg-green-500 hover:bg-green-600 disabled:bg-green-300 text-white px-6 py-3 rounded-lg font-semibold transition-colors"
            >
              {addMutation.isPending ? 'Adding...' : 'Add'}
            </button>
          </form>
        </div>

        {/* Calculator */}
        <div className="bg-white rounded-xl p-6 shadow-xl">
          <h2 className="text-xl font-semibold text-gray-800 border-b-2 border-indigo-500 pb-2 mb-4">
            Calculate Packs
          </h2>

          <form onSubmit={handleCalculate} className="flex gap-3 mb-4">
            <input
              type="number"
              value={amount}
              onChange={e => setAmount(e.target.value)}
              placeholder="Order amount"
              className="flex-1 px-4 py-3 border-2 border-gray-200 rounded-lg focus:border-indigo-500 focus:outline-none"
            />
            <button
              type="submit"
              disabled={calculateMutation.isPending}
              className="bg-indigo-500 hover:bg-indigo-600 disabled:bg-indigo-300 text-white px-6 py-3 rounded-lg font-semibold transition-colors"
            >
              {calculateMutation.isPending ? 'Calculating...' : 'Calculate'}
            </button>
          </form>

          {error && (
            <div className="bg-red-50 border-l-4 border-red-500 p-4 rounded-lg text-red-700">
              {error}
            </div>
          )}

          {result && (
            <div className="bg-green-50 border-l-4 border-green-500 p-4 rounded-lg">
              <div className="grid grid-cols-3 gap-4 mb-4">
                <div className="text-center bg-white p-4 rounded-lg shadow">
                  <div className="text-gray-500 text-sm">Order</div>
                  <div className="text-2xl font-bold text-gray-800">
                    {result.order_amount.toLocaleString()}
                  </div>
                </div>
                <div className="text-center bg-white p-4 rounded-lg shadow">
                  <div className="text-gray-500 text-sm">Items</div>
                  <div className="text-2xl font-bold text-gray-800">
                    {result.total_items.toLocaleString()}
                  </div>
                </div>
                <div className="text-center bg-white p-4 rounded-lg shadow">
                  <div className="text-gray-500 text-sm">Packs</div>
                  <div className="text-2xl font-bold text-gray-800">
                    {result.total_packs.toLocaleString()}
                  </div>
                </div>
              </div>

              <div className="flex flex-wrap gap-3 justify-center">
                {sortedPacks.map(p => (
                  <div
                    key={p.size}
                    className="bg-gradient-to-br from-indigo-500 to-purple-600 text-white px-5 py-3 rounded-lg text-center min-w-[100px]"
                  >
                    <div className="text-2xl font-bold">{p.qty}×</div>
                    <div className="text-sm opacity-90">{p.size.toLocaleString()}</div>
                  </div>
                ))}
              </div>
            </div>
          )}
        </div>
      </div>
    </div>
  )
}

export default App
