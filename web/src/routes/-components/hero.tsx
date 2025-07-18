import { ErrorWrapper } from '@/components/error'
import { Button } from '@/components/ui/button'
import OptimizedImage from '@/lib/image'
import { getRouteApi, Link } from '@tanstack/react-router'
import type { RequestHeaders } from '@tanstack/react-start/server'
import { Suspense } from 'react'

import { Words } from './word'

export default function Hero() {
  return (
    <div className="xl:h-full-navbar relative flex h-fit w-full flex-col items-start justify-stretch">
      <div className="relative z-20 flex h-full w-full flex-col items-start justify-start gap-14 py-8 md:gap-6 xl:py-12">
        <h1 className="text-2xl font-medium tracking-large text-brand-text md:text-4xl">
          Save words and phrases you find on the internet, <br /> and use them
          in the future effortlessly
        </h1>
        <div className="flex h-fit w-fit flex-row items-center justify-start gap-6">
          <Button asChild variant={'default'}>
            <Link to="/login">Get Started</Link>
          </Button>
          <Button variant={'secondary'}>Add Extension</Button>
        </div>
        <div className="flex h-fit w-fit flex-row items-center justify-start gap-6">
          <ul className="flex h-fit w-fit flex-row flex-wrap items-center justify-start gap-6">
            <li className="relative whitespace-nowrap pl-7 text-lg tracking-small text-brand-text before:absolute before:left-0 before:top-1/2 before:size-5 before:-translate-y-1/2 before:rounded-full before:bg-brand-accent">
              AI Integration
            </li>
            <li className="relative whitespace-nowrap pl-7 text-lg tracking-small text-brand-text before:absolute before:left-0 before:top-1/2 before:size-5 before:-translate-y-1/2 before:rounded-full before:bg-brand-accent">
              Chrome and Firefox Extension
            </li>
            <li className="relative whitespace-nowrap pl-7 text-lg tracking-small text-brand-text before:absolute before:left-0 before:top-1/2 before:size-5 before:-translate-y-1/2 before:rounded-full before:bg-brand-accent">
              Unlimited Vocab
            </li>
          </ul>
        </div>
        <div className="mt-auto flex h-full max-h-[250px] w-full">
          <div className="w-full max-w-full overflow-hidden">
            <ErrorWrapper>
              <Words />
            </ErrorWrapper>
          </div>
        </div>
      </div>

      <OptimizedImage
        fill
        srcLocation="/home/words/words"
        alt=""
        aria-hidden
        decoding="async"
        priority
      />
    </div>
  )
}
