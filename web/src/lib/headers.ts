import { createServerFn } from '@tanstack/react-start'
import { getHeaders } from '@tanstack/react-start/server'

function serializeHeaderObject(
  headers: Partial<Record<string, string | undefined>>,
): Record<string, string> {
  const serialized: Record<string, string> = {}
  for (const key in headers) {
    const value = headers[key]
    if (typeof value === 'string') {
      serialized[key] = value
    }
  }
  return serialized
}

const fetchSerializedHeadersFromServer = createServerFn().handler(() => {
  const headers = getHeaders()
  return serializeHeaderObject(headers)
})

export const getServerHeadersAsInstance = async () => {
  const rawHeaders = await fetchSerializedHeadersFromServer()
  return new Headers(rawHeaders)
}
