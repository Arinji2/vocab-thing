import { ErrorWrapper } from '@/components/error'
import { Button } from '@/components/ui/button'
import { generateDescription, generateTitle } from '@/lib/metadata'
import {
  useGuestLogin,
  useOAuthCallbackURL,
  type LoginProvidersType,
} from '@/lib/queries/auth'
import { cn } from '@/lib/utils'
import { createFileRoute, useNavigate } from '@tanstack/react-router'
import type { ClassValue } from 'clsx'
import { CheckCircle2, Loader2 } from 'lucide-react'
import { title } from 'radash'
import { startTransition, useEffect } from 'react'
import { toast } from 'sonner'

export const Route = createFileRoute('/login/')({
  component: RouteComponent,
  head: () => ({
    title: generateTitle('Login'),
    meta: [
      {
        name: 'description',
        content: generateDescription(
          'Login to Vocab Thing with Socials Or as a Guest. Supported Socials: Google, Discord, Github',
        ),
      },
    ],
  }),
})

function RouteComponent() {
  return (
    <div className="h-fit flex flex-col items-center justify-center py-4 gap-10 w-full xl:h-full-navbar screen-padding ">
      <h1 className="text-3xl text-center md:text-5xl font-bold text-brand-text tracking-large">
        Continue To Vocab Thing
      </h1>
      <div className="w-full gap-10 h-fit flex flex-col xl:flex-row items-center xl:items-start  justify-between">
        <div className="w-full xl:w-[550px] md:w-[80%] rounded-lg shadow-black shadow-lg h-fit xl:h-[400px] gap-8 bg-brand-secondary-dark flex flex-col items-start justify-start px-8 py-8">
          <h2 className="text-2xl font-medium text-brand-text tracking-large">
            Login With Socials
          </h2>
          <div className="flex flex-col items-start justify-start gap-3">
            <div className="w-full h-fit flex flex-row items-center justify-start gap-3">
              <CheckCircle2 className="text-brand-accent" />
              <p className="text-brand-text  text-lg">Free Forever</p>
            </div>
            <div className="w-full h-fit flex flex-row items-center justify-start gap-3">
              <CheckCircle2 className="text-brand-accent" />
              <p className="text-brand-text  text-lg">Access To All Features</p>
            </div>
            <div className="w-full h-fit flex flex-row items-center justify-start gap-3">
              <CheckCircle2 className="text-brand-accent" />
              <p className="text-brand-text  text-lg">15 AI Usage Per Day</p>
            </div>
          </div>
          <div className="w-full h-fit flex flex-col items-center justify-center">
            <ErrorWrapper>
              <div className="w-full h-fit grid md:grid-cols-2 grid-cols-1 mt-auto gap-4">
                <LoginButton provider="google" />
                <LoginButton provider="discord" />
                <LoginButton provider="github" className="md:col-span-2" />
              </div>
            </ErrorWrapper>
          </div>
        </div>
        <div className="md:w-[80%] w-full xl:w-[550px] rounded-lg shadow-black shadow-lg h-fit gap-8 bg-brand-offwhite-dark flex flex-col items-start justify-start px-8 py-8">
          <h2 className="text-2xl font-medium text-brand-text tracking-large">
            Login As Guest
          </h2>
          <div className="flex flex-col items-start justify-start gap-3">
            <div className="w-full h-fit flex flex-row items-center justify-start gap-3">
              <CheckCircle2 className="text-brand-accent" />
              <p className="text-brand-text  text-lg">Free Forever</p>
            </div>
            <div className="w-full h-fit flex flex-row items-center justify-start gap-3">
              <CheckCircle2 className="text-brand-accent" />
              <p className="text-brand-text  text-lg">Access To All Features</p>
            </div>
            <div className="w-full h-fit flex flex-row items-center justify-start gap-3">
              <CheckCircle2 className="text-brand-accent" />
              <p className="text-brand-text  text-lg">2 AI Usage Per Day</p>
            </div>
          </div>
          <ErrorWrapper>
            <div className="w-full h-fit flex flex-col items-center justify-center">
              <LoginButton provider="guest" />
            </div>
          </ErrorWrapper>
        </div>
      </div>
    </div>
  )
}

function LoginButton({
  provider,
  className,
}: {
  provider: LoginProvidersType
  className?: ClassValue
}) {
  const navigate = useNavigate()

  const {
    mutate: oauthMutation,
    isPending: isOAuthPending,
    isSuccess: isOAuthSuccess,
    isError: isOAuthError,
    error: oauthError,
    data: oauthData,
  } = useOAuthCallbackURL()
  const {
    mutate: guestLogin,
    isPending: isGuestPending,
    isSuccess: isGuestSuccess,
    isError: isGuestError,
    error: guestError,
  } = useGuestLogin()

  useEffect(() => {
    if (isGuestSuccess) {
      navigate({
        to: '/dashboard',
      })
    }
    if (isOAuthSuccess && oauthData.codeURL) {
      window.location.href = oauthData.codeURL
    }

    if (isOAuthError) {
      toast.error(oauthError.message)
    }

    if (isGuestError) {
      toast.error(guestError.message)
    }
  }, [
    isGuestSuccess,
    isOAuthSuccess,
    oauthData,
    isGuestError,
    isOAuthError,
    guestError,
    oauthError,
  ])
  return (
    <Button
      onClick={() => {
        startTransition(() => {
          if (provider === 'guest') guestLogin()
          else oauthMutation(provider)
        })
      }}
      className={cn(
        'relative flex w-full flex-row items-center justify-center overflow-hidden gap-2 text-base bg-brand-primary-dark text-brand-text ',
        className,
      )}
    >
      <div
        className={cn(
          'flex flex-col items-center justify-center w-full h-full absolute top-0 left-0 transition-all ease-in-out duration-200 -translate-y-full bg-brand-primary-dark',
          {
            'translate-y-0': isGuestPending || isOAuthPending,
          },
        )}
      >
        <Loader2 className="text-brand-text size-9 animate-spin" />
      </div>
      Login {provider === 'guest' ? 'as' : 'with'} {title(provider)}
    </Button>
  )
}
