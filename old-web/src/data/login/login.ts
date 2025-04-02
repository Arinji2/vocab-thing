import z from "zod";
export const LoginProviders = z.enum(["google", "discord", "github", "guest"]);
export type LoginProvidersType = z.infer<typeof LoginProviders>;

export const OauthCallbackURLSchema = z.object({
  codeURL: z.string(),
});

export type OauthCallbackURLSchemaType = z.infer<typeof OauthCallbackURLSchema>;
