import { queryOptions } from '@tanstack/react-query'
import z from 'zod'

export const WordSchema = z.object({
  id: z.string(),
  word: z.string(),
  definition: z.string(),
})

export type WordSchemaType = z.infer<typeof WordSchema>

export const PBSchema = z.object({
  items: z.array(z.unknown()),
  page: z.number(),
  perPage: z.number(),
  totalItems: z.number(),
  totalPages: z.number(),
})

export const fetchWords = async (): Promise<Array<WordSchemaType>> => {
  const res = await fetch(
    'https://db-word.arinji.com/api/collections/real_words/records?perPage=3&sort=@random,level',
  )

  if (!res.ok) {
    throw new Error('Failed to fetch words')
  }

  const rawData = await res.json()
  const parsed = PBSchema.parse(rawData)

  return parsed.items.map((raw) => WordSchema.parse(raw))
}

export const wordsQueryOptions = queryOptions({
  queryKey: ['words'],
  queryFn: fetchWords,
  staleTime: Infinity,
  gcTime: 1000 * 60 * 60, // 1 hour
  refetchOnWindowFocus: false,
  refetchOnMount: false,
})
