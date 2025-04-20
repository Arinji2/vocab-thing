import { createServerFn } from '@tanstack/react-start'
import { getCookie, getEvent } from '@tanstack/react-start/server'
import { parse } from 'cookie'

/**
 * Quick check: just verifies if an `oauth_session` cookie exists.
 */
export const isLoggedInQuick = createServerFn({ method: 'GET' }).handler(() => {
  const event = getEvent()
  const cookie = getCookie(event, 'oauth_session')
  const cookies = parse(cookie ?? '')
  return Boolean(cookies['oauth_session'])
})

/**
 * Deep check: validate session token with DB (placeholder for now).
 */
export const isLoggedInDeep = async (req: Request): Promise<boolean> => {
  // TODO: Actually implement this
  await new Promise((resolve) => setTimeout(resolve, 10))
  const cookieHeader = req.headers.get('cookie') ?? ''
  const cookies = parse(cookieHeader)
  const sessionToken = cookies['oauth_session']
  if (!sessionToken) return false
  return true
}
