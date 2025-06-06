import { useQuery } from '@tanstack/react-query'
import {
  getDexiePhrase,
  getDexiePhrases,
  updateDexiePhrase,
  updateDexiePhraseTag,
} from './helpers'
import { fetchPhraseFromAPI, fetchPhrasesFromAPI } from './api'
import type { PhraseResponseSchemaType } from './schema'

type PhraseByIDResult =
  | { success: true; data: PhraseResponseSchemaType }
  | { success: false; error: string }
const fetchPhraseByID = async (
  id: string,
  isSyncing: boolean,
): Promise<PhraseByIDResult> => {
  if (!isSyncing) {
    const local = await getDexiePhrase(id)
    if (local.phrase)
      return {
        success: true,
        data: {
          phrase: local.phrase,
          tag: local.tag,
        },
      }
  }
  try {
    const data = await fetchPhraseFromAPI(id)
    if (!isSyncing) {
      void (async () => {
        try {
          await updateDexiePhrase(id, data.phrase)
          await Promise.all(
            data.tag.map((tag) => updateDexiePhraseTag(tag.id, tag)),
          )
        } catch (err) {
          console.warn('Dexie update failed:', err)
        }
      })()
    }

    return { success: true, data }
  } catch (e: any) {
    return {
      success: false,
      error: e?.message ?? 'Unexpected error',
    }
  }
}
export const usePhraseByID = (id: string, isSyncing: boolean) =>
  useQuery({
    queryKey: ['phrase', id],
    enabled: !!id,
    queryFn: () => fetchPhraseByID(id, isSyncing),
    refetchOnWindowFocus: false,
    refetchOnMount: false,
  })

type PhrasesResult =
  | { success: true; data: Array<PhraseResponseSchemaType> }
  | { success: false; error: string }
const fetchPhrases = async (isSyncing: boolean): Promise<PhrasesResult> => {
  if (!isSyncing) {
    const local = await getDexiePhrases()
    if (local.length > 0)
      return {
        success: true,
        data: local,
      }
  }
  try {
    const data = await fetchPhrasesFromAPI()
    if (!isSyncing) {
      void (async () => {
        try {
          await Promise.all(
            data.map(async (item) => {
              await updateDexiePhrase(item.phrase.id, item.phrase)
              await Promise.all(
                item.tag.map((tag) => updateDexiePhraseTag(tag.id, tag)),
              )
            }),
          )
        } catch (err) {
          console.warn('Dexie update failed:', err)
        }
      })()
    }

    return { success: true, data }
  } catch (e: any) {
    return {
      success: false,
      error: e?.message ?? 'Unexpected error',
    }
  }
}
export const usePhrases = (isSyncing: boolean) =>
  useQuery({
    queryKey: ['phrase'],
    queryFn: () => fetchPhrases(isSyncing),
    refetchOnWindowFocus: false,
    refetchOnMount: false,
  })
