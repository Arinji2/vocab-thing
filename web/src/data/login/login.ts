import z from "zod";
export type LoginProviders = "google" | "discord" | "github" | "guest";

export const OauthCallbackURLSchema = z.object({
  codeURL: z.string(),
});

export type OauthCallbackURLSchemaType = z.infer<typeof OauthCallbackURLSchema>;
