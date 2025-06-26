import { useQuery } from '@tanstack/react-query'
import { fetchPhraseFromAPI, fetchPhrasesFromAPI } from './api'
import type { PhraseResponseSchemaType } from './schema'

type PhraseByIDResult =
  | { success: true; data: PhraseResponseSchemaType }
  | { success: false; error: string }
const fetchPhraseByID = async (id: string): Promise<PhraseByIDResult> => {
  try {
    const data = await fetchPhraseFromAPI(id)

    return { success: true, data }
  } catch (e: any) {
    return {
      success: false,
      error: e?.message ?? 'Unexpected error',
    }
  }
}
export const usePhraseByID = (id: string) =>
  useQuery({
    queryKey: ['phrase', id],
    enabled: !!id,
    queryFn: () => fetchPhraseByID(id),
    refetchOnWindowFocus: false,
    refetchOnMount: false,
  })

type PhrasesResult =
  | { success: true; data: Array<PhraseResponseSchemaType> }
  | { success: false; error: string }
const fetchPhrases = async (): Promise<PhrasesResult> => {
  try {
    const data = await fetchPhrasesFromAPI()
    return { success: true, data }
  } catch (e: any) {
    return {
      success: false,
      error: e?.message ?? 'Unexpected error',
    }
  }
}
export const usePhrases = () =>
  useQuery({
    queryKey: ['phrase'],
    queryFn: () => fetchPhrases(),
    refetchOnWindowFocus: false,
    refetchOnMount: false,
  })
