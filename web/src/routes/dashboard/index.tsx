import { RouteError } from '@/components/error'
import { useSyncing } from '@/lib/is-syncing'
import { usePhrases } from '@/lib/queries/phrase/index'
import { useQueryClient } from '@tanstack/react-query'
import { createFileRoute } from '@tanstack/react-router'
import { useEffect } from 'react'

export const Route = createFileRoute('/dashboard/')({
  component: RouteComponent,
  errorComponent: () => <RouteError />,
})

function RouteComponent() {
  const { isSyncing } = useSyncing()
  const { isLoading, isError, data, error } = usePhrases(isSyncing)
  const queryClient = useQueryClient()
  useEffect(() => {
    console.log('data', data)
  }, [data])
  return (
    <div className="w-full h-full-navbar flex flex-col items-center justify-center">
      {data && data.success && <div className="text-white"></div>}
    </div>
  )
}
