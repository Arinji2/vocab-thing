import { queryOptions } from "@tanstack/react-query";
import { createServerFn } from "@tanstack/react-start";
import z from "zod";
import { PBSchema } from "~/utils/DAL/pb";

export const WordSchema = z.object({
  id: z.string(),
  word: z.string(),
  definition: z.string(),
});

export type WordSchemaType = z.infer<typeof WordSchema>;

export const fetchWords = createServerFn({ method: "GET" }).handler(
  async () => {
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
  },
);

export const wordsQueryOptions = () =>
  queryOptions({
    queryKey: ["words"],
    queryFn: fetchWords,
    staleTime: Infinity,
    gcTime: 1000 * 60 * 60 * 24, // 24 hours
    retry: false,
  });
