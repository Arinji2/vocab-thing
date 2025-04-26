import { getApiURL } from '@/lib/getApiURL'
import { PhraseResponseSchema, PhrasesResponseSchema } from './schema'

export const fetchPhraseFromAPI = async (id: string) => {
  try {
    const res = await fetch(`${getApiURL()}/phrase/${id}`, {
      method: 'GET',
      credentials: 'include',
    })

    if (!res.ok) throw new Error('Failed to fetch phrase')

    const data = await res.json()
    return PhraseResponseSchema.parse(data)
  } catch (e: any) {
    if (e instanceof TypeError && !navigator.onLine) {
      throw new Error('NETWORK-ERROR')
    }

    if (e instanceof Error) {
      throw e
    } else {
      throw new Error('Unexpected error, please try again')
    }
  }
}

export const fetchPhrasesFromAPI = async () => {
  try {
    const res = await fetch(`${getApiURL()}/phrase/`, {
      method: 'GET',
      credentials: 'include',
    })

    if (!res.ok) throw new Error('Failed to fetch phrase')

    const data = await res.json()
    return PhrasesResponseSchema.parse(data)
  } catch (e: any) {
    if (e instanceof TypeError && !navigator.onLine) {
      throw new Error('NETWORK-ERROR')
    }

    if (e instanceof Error) {
      throw e
    } else {
      throw new Error('Unexpected error, please try again')
    }
  }
}
