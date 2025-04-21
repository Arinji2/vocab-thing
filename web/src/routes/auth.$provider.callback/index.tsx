import { RouteError } from '@/components/error'
import { getServerHeadersAsInstance } from '@/lib/headers'
import { LoginProviders, useSocialLogin } from '@/lib/queries/auth'
import {
  createFileRoute,
  getRouteApi,
  redirect,
  useNavigate,
} from '@tanstack/react-router'
import { Loader2 } from 'lucide-react'
import { useEffect, useState } from 'react'
import { z } from 'zod'

const oauthSearchSchema = z.object({
  state: z.string().catch(''),
  code: z.string().catch(''),
})
type oauthSearch = z.infer<typeof oauthSearchSchema>

export const Route = createFileRoute('/auth/$provider/callback/')({
  component: RouteComponent,
  validateSearch: (search) => oauthSearchSchema.parse(search),
  loaderDeps: ({ search: { state, code } }) => ({ state, code }),
  loader: async ({ params, deps: { state, code } }) => {
    if (!state || !code) {
      throw new Error()
    }
    const validatedProvider = LoginProviders.safeParse(params.provider)
    if (!validatedProvider.success) {
      console.error('Invalid provider:', validatedProvider.error)
      throw redirect({
        to: '/login',
      })
    }

    if (validatedProvider.data === 'guest') {
      throw redirect({
        to: '/login',
      })
    }
    const headerStore = await getServerHeadersAsInstance()
    const userIP =
      headerStore.get('x-vercel-forwarded-for') ||
      headerStore.get('x-forwarded-for') ||
      '0.0.0.0'

    const userAgent = headerStore.get('user-agent') ?? 'Unknown'

    return {
      params,
      state,
      code,
      provider: validatedProvider.data,
      user: {
        ip: userIP,
        agent: userAgent,
      },
    }
  },
})

function RouteComponent() {
  const {
    mutate: loginSocial,
    isPending,
    isError,
    isSuccess,
    error,
    data,
  } = useSocialLogin()
  const routeApi = getRouteApi('/auth/$provider/callback/')
  const { provider, state, code, user } = routeApi.useLoaderData()
  const [timeElapsed, setTimeElapsed] = useState(0)
  const navigate = useNavigate()

  useEffect(() => {
    loginSocial({
      providerType: provider,
      code,
      state,
      fingerprint: user.agent,
      ip: user.ip,
    })
  }, [])

  useEffect(() => {
    if (isSuccess) {
      navigate({
        to: '/dashboard',
      })
    }
  }, [isSuccess])

  useEffect(() => {
    const timer = setTimeout(() => {
      setTimeElapsed(timeElapsed + 1)
    }, 1000)
    if (timeElapsed > 60) {
      setTimeout(() => {
        navigate({
          to: '/login',
        })
      }, 2000)
    }
    return () => clearTimeout(timer)
  }, [timeElapsed])

  return (
    <div className="w-full h-fit flex flex-col items-center gap-6 justify-center">
      <div className="text-brand-text w-full h-fit flex flex-col items-center justify-center gap-6">
        <div className="text-xl font-medium tracking-large w-full h-fit flex flex-row items-center justify-center gap-2">
          <Loader2
            className="animate-spin text-brand-accent"
            size={25}
            strokeWidth={3}
          />{' '}
          Authenticating Securely
        </div>
        <p>Time Elapsed: {timeElapsed} seconds</p>
      </div>
      {error && (
        <div className="text-red-500">
          Authentication failed. Redirecting to login...
        </div>
      )}
      {timeElapsed > 60 && (
        <div className="flex flex-col items-center justify-center">
          <p className="text-xl text-brand-destructive-light">
            Authentication timed out. Redirecting to login...
          </p>
        </div>
      )}
    </div>
  )
}
