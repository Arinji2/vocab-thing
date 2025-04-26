import GlobalLoading from '@/components/loading.global'
import { HomeNavbar } from '@/components/navbar/home'
import { GlobalNotFound } from '@/components/not-found.global'
import { Toaster } from '@/components/ui/sonner'
import { checkSessionCookie } from '@/lib/is-logged-in'
import { generateDescription, generateTitle } from '@/lib/metadata'
import type { QueryClient } from '@tanstack/react-query'
import {
  createRootRouteWithContext,
  HeadContent,
  Outlet,
  Scripts,
} from '@tanstack/react-router'
import { TanStackRouterDevtools } from '@tanstack/react-router-devtools'

import TanstackQueryLayout from '../integrations/tanstack-query/layout'
import appCss from '../styles.css?url'
import { RouteError } from '@/components/error'
import { SyncingProvider } from '@/lib/is-syncing'

interface MyRouterContext {
  queryClient: QueryClient
}

export const Route = createRootRouteWithContext<MyRouterContext>()({
  head: () => ({
    meta: [
      { charSet: 'utf-8' },
      {
        name: 'viewport',
        content: 'width=device-width, initial-scale=1',
      },
      {
        name: 'description',
        content: generateDescription(),
      },
      {
        title: generateTitle(),
      },
      {
        name: 'keywords',
        content:
          'vocab, vocabthing, arinji, arinji.com, arinjay dhar, save words, phrases',
      },
      {
        name: 'theme-color',
        content: '#89DFE9',
      },
    ],
    links: [
      {
        rel: 'preload',
        href: '/fonts/Tektur-Regular.ttf',
        as: 'font',
        type: 'font/ttf',
        crossOrigin: 'anonymous',
      },

      {
        rel: 'preload',
        href: '/fonts/Tektur-Medium.ttf',
        as: 'font',
        type: 'font/ttf',
        crossOrigin: 'anonymous',
      },

      {
        rel: 'preload',
        href: '/fonts/Tektur-Bold.ttf',
        as: 'font',
        type: 'font/ttf',
        crossOrigin: 'anonymous',
      },
      {
        rel: 'stylesheet',
        href: appCss,
      },
      {
        rel: 'icon',
        href: '/metadata/favicon-16x16.png',
        sizes: '16x16',
        type: 'image/png',
      },
      {
        rel: 'icon',
        href: '/metadata/favicon-32x32.png',
        sizes: '32x32',
        type: 'image/png',
      },
      {
        rel: 'icon',
        href: '/metadata/favicon-96x96.png',
        sizes: '96x96',
        type: 'image/png',
      },
      { rel: 'icon', href: '/metadata/favicon.svg', type: 'image/svg+xml' },
      { rel: 'shortcut icon', href: '/metadata/favicon.ico' },
      {
        rel: 'apple-touch-icon',
        href: '/metadata/apple-touch-icon.png',
        sizes: '180x180',
      },
      { rel: 'manifest', href: '/metadata/site.webmanifest' },
    ],
  }),
  pendingComponent: () => <GlobalLoading />,
  notFoundComponent: () => <GlobalNotFound />,
  errorComponent: () => <RouteError />,
  loader: async () => {
    const isLoggedIn = await checkSessionCookie()
    return {
      isLoggedIn,
    }
  },

  component: () => (
    <RootDocument>
      <Outlet />
      <TanStackRouterDevtools />
      <TanstackQueryLayout />
    </RootDocument>
  ),
})

function RootDocument({ children }: { children: React.ReactNode }) {
  return (
    <html lang="en" className="bg-brand-background">
      <head>
        <HeadContent />
      </head>
      <body
        className={`flex h-full w-full flex-col items-center justify-start antialiased`}
      >
        <div className="flex w-full max-w-[1280px] flex-col items-center justify-start">
          <HomeNavbar />
          <SyncingProvider>{children}</SyncingProvider>
          <Toaster />
        </div>
        <Scripts />
      </body>
    </html>
  )
}
