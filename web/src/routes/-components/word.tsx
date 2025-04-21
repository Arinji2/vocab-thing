import { Button } from '@/components/ui/button'
import { Skeleton } from '@/components/ui/skeleton'
import {
  wordsQueryOptions,
  type WordSchemaType,
} from '@/lib/queries/fetchWords'
import { useQuery } from '@tanstack/react-query'
import { getRouteApi, useNavigate } from '@tanstack/react-router'
import { toast } from 'sonner'

export function Words() {
  const { data, isLoading, error } = useQuery(wordsQueryOptions)
  if (error) {
    throw error
  }
  return (
    <div className="flex h-full w-full flex-col items-start justify-start gap-2">
      <p className="text-lg font-medium tracking-small text-brand-text">
        Powered By{' '}
        <a
          href="https://sense.arinji.com"
          target="_blank"
          rel="noreferrer"
          className="underline decoration-brand-offwhite underline-offset-4"
        >
          <span className="text-green-500">Sense</span> Or{' '}
          <span className="text-red-500">Nonsense</span>
        </a>
      </p>
      <div className="flex h-full w-full flex-row items-center snap-x snap-proximity justify-start gap-10 overflow-x-auto">
        {isLoading
          ? Array.from({ length: 3 }, (_, i) => <WordSuspenseCard key={i} />)
          : data?.map((d) => <WordCard key={d.id} data={d} />)}
      </div>
    </div>
  )
}

function WordSuspenseCard() {
  return (
    <article className="flex h-full min-h-[200px] w-full snap-center md:w-[350px] shrink-0 flex-col items-start justify-start gap-3 rounded-xl bg-brand-secondary-dark p-4">
      <div className="flex h-fit w-full flex-col items-start justify-start gap-1">
        <p className="text-sm tracking-small text-brand-text">title</p>
        <Skeleton className="h-5 w-full" />
      </div>
      <div className="flex h-fit w-full flex-col items-start justify-start gap-1">
        <p className="text-sm tracking-small text-brand-text">description</p>
        <Skeleton className="h-[18px] w-full" />
      </div>
      <Button variant={'default'} disabled className="mt-auto" size={'sm'}>
        Add to Vocab
      </Button>
    </article>
  )
}

export function WordCard({ data }: { data: WordSchemaType }) {
  const navigate = useNavigate()
  const routeApi = getRouteApi('__root__')
  const loaderData = routeApi.useLoaderData()
  const { isLoggedIn } = loaderData
  return (
    <article className="snap-center flex h-full min-h-[200px] w-full md:w-[350px] shrink-0 flex-col items-start justify-start gap-3 rounded-xl bg-brand-secondary-dark p-4">
      <div className="flex h-fit w-full flex-col items-start justify-start">
        <p className="text-sm tracking-small text-brand-text">title</p>
        <p className="line-clamp-1 text-xl font-semibold tracking-small text-brand-primary">
          {data.word}
        </p>
      </div>
      <div className="flex h-fit w-full flex-col items-start justify-start">
        <p className="text-sm tracking-small text-brand-text">description</p>
        <p className="line-clamp-2 text-lg font-semibold tracking-small text-brand-primary">
          {data.definition}
        </p>
      </div>
      <Button
        onClick={() => {
          toast.success(`${data.word} Added to Vocab`)
          localStorage.setItem('homepage-word', data.word)
          localStorage.setItem('homepage-definition', data.definition)
          if (isLoggedIn) {
            navigate({
              to: '/dashboard',
            })
          } else {
            navigate({
              to: '/login',
            })
          }
        }}
        variant={'default'}
        className="mt-auto"
        size={'sm'}
      >
        Add to Vocab
      </Button>
    </article>
  )
}
