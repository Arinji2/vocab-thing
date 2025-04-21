import { cn } from '@/lib/utils'
import type { ClassValue } from 'clsx'
import { useState } from 'react'
import { ErrorBoundary } from 'react-error-boundary'
import { Button } from './ui/button'
import { Link } from '@tanstack/react-router'

export function ErrorWrapper({
  children,
  className,
}: {
  children: React.ReactNode
  className?: ClassValue
}) {
  const [errorKey, setErrorKey] = useState(0) // Track key to force remount on error

  return (
    <ErrorBoundary
      onReset={() => setErrorKey((prev) => prev + 1)}
      fallbackRender={({ error, resetErrorBoundary }) => (
        <div className="flex h-full w-full flex-col items-center justify-center gap-4 rounded-xl bg-brand-destructive-dark/50 p-4">
          <h2 className="text-2xl font-semibold tracking-small text-brand-text">
            Something went wrong!
          </h2>
          <p className="text-lg text-brand-text">{error?.message}</p>
          <Button variant="default" onClick={resetErrorBoundary}>
            Try Again
          </Button>
        </div>
      )}
    >
      <div key={errorKey} className={cn('w-full h-full', className)}>
        {children}
      </div>
    </ErrorBoundary>
  )
}

export function RouteError() {
  return (
    <div className="flex min-h-screen w-full flex-col items-center justify-center gap-6 bg-brand-destructive-dark text-brand-text px-4">
      <h1 className="text-3xl font-bold tracking-large">
        Oops! Something went wrong.
      </h1>
      <p className="text-lg text-brand-text/80 max-w-md text-center">
        The page failed to load properly or some required data was missing.
      </p>
      <Link
        to="/"
        className="rounded-xl bg-brand-primary px-6 py-3 text-black text-lg font-medium transition hover:bg-brand-primary-dark"
      >
        Go back home
      </Link>
    </div>
  )
}
