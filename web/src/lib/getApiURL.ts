import { env } from '@/env'

export function getApiURL() {
  const envURL = env.VITE_API_URL
  if (envURL) {
    return envURL
  } else {
    console.error('VITE_API_URL is not set')
    return 'http://localhost:8080'
  }
}
