import {
  HeadContent,
  Outlet,
  Scripts,
  createRootRouteWithContext,
} from '@tanstack/react-router'
import { TanStackRouterDevtools } from '@tanstack/react-router-devtools'

import TanstackQueryLayout from '../integrations/tanstack-query/layout'
import appCss from '../styles.css?url'

import type { QueryClient } from '@tanstack/react-query'
import { generateDescription, generateTitle } from '@/lib/metadata'

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
    <html lang="en">
      <head>
        <HeadContent />
      </head>
      <body>
        {children}
        <Scripts />
      </body>
    </html>
  )
}
