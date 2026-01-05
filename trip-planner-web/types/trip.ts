export type TripRequest = {
  destination: string
  start_date?: string
  days: number
  budget?: 'low' | 'mid' | 'high'
  interests?: string[]
  pace?: 'chill' | 'balanced' | 'fast'
  notes?: string
}

export type TripPlan = {
  id: string
  input_hash: string
  request: TripRequest
  itinerary: any
  weather?: any
  places?: any
  created_at: number
  updated_at: number
}
