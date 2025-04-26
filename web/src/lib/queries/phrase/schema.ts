import { z } from 'zod'

export const PhraseSchema = z.object({
  id: z.string(),
  user_id: z.string(),
  phrase: z.string(),
  phrase_definition: z.string(),
  pinned: z.boolean(),
  found_in: z.string(),
  public: z.boolean(),
  usage_count: z.number(),
  created_at: z.string(),
  updated_at: z.string(),
  deleted_at: z.string().nullable(),
})

export const PhraseTagSchema = z.object({
  id: z.string(),
  phrase_id: z.string(),
  tag_name: z.string(),
  tag_color: z.string(),
  created_at: z.string(),
})

export const PhraseResponseSchema = z.object({
  phrase: PhraseSchema,
  tag: z.array(PhraseTagSchema),
})
export const PhrasesResponseSchema = z.array(PhraseResponseSchema)

export type PhraseSchemaType = z.infer<typeof PhraseSchema>
export type PhraseTagSchemaType = z.infer<typeof PhraseTagSchema>
export type PhraseResponseSchemaType = z.infer<typeof PhraseResponseSchema>
