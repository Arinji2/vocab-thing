import { isLoggedInQuick } from '@/lib/isLoggedIn'
import { wordsQueryOptions } from '@/lib/queries/fetchWords'
import { createFileRoute } from '@tanstack/react-router'

import Features from './-components/features'
import Footer from './-components/footer'
import Hero from './-components/hero'
import Info from './-components/info'
import Works from './-components/works'

export const Route = createFileRoute('/')({
  component: App,
  loader: async ({ context }) => {
    const isLoggedIn = await isLoggedInQuick()
    context.queryClient.ensureQueryData(wordsQueryOptions)
    return {
      isLoggedIn,
    }
  },
})

function App() {
  return (
    <div className="flex h-fit w-full pb-10 gap-20 screen-padding flex-col items-center justify-start">
      <Hero />
      <Works />
      <Features />
      <Info />
      <Footer />
    </div>
  )
}
