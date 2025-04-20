import { cn } from '@/lib/utils'

function Skeleton({
  className,
  ...props
}: React.HTMLAttributes<HTMLDivElement>) {
  return (
    <div
      className={cn('bg-brand-offwhite/30 animate-pulse rounded-md', className)}
      {...props}
    />
  )
}

export { Skeleton }
