import { useMutation } from '@tanstack/react-query'
import { HandleResponseError } from '../errors'
import { getApiURL } from '../getApiURL'

import z from 'zod'

export const LoginProviders = z.enum(['google', 'discord', 'github', 'guest'])
export type LoginProvidersType = z.infer<typeof LoginProviders>

export const OauthCallbackURLSchema = z.object({
  codeURL: z.string(),
})

export type OauthCallbackURLSchemaType = z.infer<typeof OauthCallbackURLSchema>
export function useGuestLogin() {
  return useMutation({
    mutationFn: async () => {
      const apiURL = getApiURL()
      const res = await fetch(`${apiURL}/user/create/guest`, {
        method: 'POST',
        credentials: 'include',
        cache: 'no-store',
      })

      const resError = await HandleResponseError('Guest Login', res)
      if (resError) {
        throw new Error(resError.readable)
      }

      return { success: true }
    },
  })
}

export function useSocialLogin() {
  return useMutation({
    mutationFn: async ({
      providerType,
      code,
      state,
      fingerprint,
      ip,
    }: {
      providerType: string
      code: string
      state: string
      fingerprint: string
      ip: string
    }) => {
      const apiURL = getApiURL()

      const res = await fetch(`${apiURL}/oauth/callback`, {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        credentials: 'include',
        body: JSON.stringify({ providerType, code, state, fingerprint, ip }),
      })

      const resError = await HandleResponseError('Social Login', res)
      if (resError) {
        throw new Error(resError.readable)
      }

      return true
    },
  })
}

export function useOAuthCallbackURL() {
  return useMutation({
    mutationFn: async (providerType: LoginProvidersType) => {
      console.log('HERE')
      const apiURL = getApiURL()

      const res = await fetch(`${apiURL}/oauth/generate-code-url`, {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        credentials: 'include',
        body: JSON.stringify({ providerType }),
      })

      const resError = await HandleResponseError('OAuth callback URL', res)
      if (resError) {
        throw new Error(resError.readable)
      }

      const data = await res.json()
      const parsed = OauthCallbackURLSchema.parse(data)

      return {
        success: true,
        codeURL: parsed.codeURL,
      }
    },
  })
}
