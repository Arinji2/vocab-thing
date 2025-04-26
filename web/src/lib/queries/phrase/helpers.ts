import { dexieDB } from '@/lib/dexie/db'
import type { PhraseSchemaType, PhraseTagSchemaType } from './schema'

export const updateDexiePhrase = async (
  id: string,
  phrase: PhraseSchemaType,
) => {
  await dexieDB.phrase.put({ ...phrase, id })
}

export const getDexiePhrase = async (id: string) => {
  const phraseData = await dexieDB.phrase.get(id)
  const tags = await getAllDexiePhraseTags(id)
  return { phrase: phraseData, tag: tags }
}

export const getDexiePhrases = async () => {
  const phraseData = await dexieDB.phrase.toArray()
  return await Promise.all(
    phraseData.map(async (phrase) => {
      const tags = await getAllDexiePhraseTags(phrase.id)
      return { phrase, tag: tags }
    }),
  )
}

export const updateDexiePhraseTag = async (
  id: string,
  tag: PhraseTagSchemaType,
) => {
  await dexieDB.tag.put({ ...tag, id })
}

export const getDexiePhraseTag = async (id: string) => {
  return await dexieDB.tag.get(id)
}

export const getAllDexiePhraseTags = async (phraseID: string) => {
  return await dexieDB.tag.where('phrase_id').equals(phraseID).toArray()
}
