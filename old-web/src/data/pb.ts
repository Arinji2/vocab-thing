import z from "zod";
export const PBSchema = z.object({
  items: z.array(z.unknown()),
  page: z.number(),
  perPage: z.number(),
  totalItems: z.number(),
  totalPages: z.number(),
});
