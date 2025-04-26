import { createServerFn } from '@tanstack/react-start'
import { getCookie, getEvent } from '@tanstack/react-start/server'
import { parse } from 'cookie'

/**
 * Quick check: just verifies if an `session` cookie exists.
 */
export const checkSessionCookie = createServerFn({ method: 'GET' }).handler(
  () => {
    const event = getEvent()
    const cookieHeader = event.node.req.headers.cookie
    console.log('full cookie header:', cookieHeader)

    const cookie = getCookie(event, 'session')
    return Boolean(cookie)
  },
)

/**
 * Deep check: validate session token with DB (placeholder for now).
 */
export const validateSessionCookie = async (req: Request): Promise<boolean> => {
  // TODO: Actually implement this
  await new Promise((resolve) => setTimeout(resolve, 10))
  const cookieHeader = req.headers.get('cookie') ?? ''
  const cookies = parse(cookieHeader)
  const sessionToken = cookies['session']
  if (!sessionToken) return false
  return true
}
