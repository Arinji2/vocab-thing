import Dexie, { type EntityTable } from 'dexie'
import type {
  PhraseSchemaType,
  PhraseTagSchemaType,
} from '../queries/phrase/schema'

const dexieDB = new Dexie('VocabDB') as Dexie & {
  phrase: EntityTable<PhraseSchemaType, 'id'>
  tag: EntityTable<PhraseTagSchemaType, 'id'>
}

dexieDB.version(1).stores({
  phrase:
    'id, user_id, phrase, phrase_definition, pinned, found_in, public, usageCount, created_at, updated_at, deleted_at',
  tag: 'id, phrase_id, tag_name, tag_color, created_at',
})

export { dexieDB }
