import { Link } from '@tanstack/react-router'

import { Button } from './ui/button'

export function GlobalNotFound() {
  return (
    <div className="flex h-full-navbar w-full flex-col items-center justify-center gap-4 rounded-xl p-4">
      <h2 className="text-2xl font-semibold tracking-small text-brand-text">
        404 – Not Found
      </h2>
      <p className="text-lg text-brand-text">
        The page you're looking for doesn't exist.
      </p>
      <Link to="/">
        <Button variant="default">Go Back Home</Button>
      </Link>
    </div>
  )
}
