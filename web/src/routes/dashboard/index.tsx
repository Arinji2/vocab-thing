import { RouteError } from '@/components/error'
import { usePhrases } from '@/lib/queries/phrase/index'
import { useQueryClient } from '@tanstack/react-query'
import { createFileRoute } from '@tanstack/react-router'

export const Route = createFileRoute('/dashboard/')({
  component: RouteComponent,
  errorComponent: () => <RouteError />,
})

function RouteComponent() {
  const { isLoading, isError, data, error } = usePhrases()
  const queryClient = useQueryClient()
  return (
    <div className="w-full h-full-navbar flex flex-col items-center justify-center">
      {data && data.success && <div className="text-white"></div>}
    </div>
  )
}
