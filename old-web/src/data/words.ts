import z from "zod";
import { PBSchema } from "./pb";

export const WordSchema = z.object({
  id: z.string(),
  word: z.string(),
  definition: z.string(),
});

export type WordSchemaType = z.infer<typeof WordSchema>;

export const fetchWords = async () => {
  try {
    const res = await fetch(
      "https://db-word.arinji.com/api/collections/real_words/records?perPage=3&sort=@random,level",
    );

    if (!res.ok) {
      console.error("Fetch failed with status:", res.status);
      throw new Error("Failed to fetch words");
    }

    const rawData = await res.json();
    const data = PBSchema.parse(rawData);

    return data.items.map((rawWord) => WordSchema.parse(rawWord));
  } catch (error) {
    console.error("Error in fetchWords:", error);
    throw error;
  }
};
