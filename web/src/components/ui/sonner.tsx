import { useTheme } from 'next-themes'
import { Toaster as Sonner } from 'sonner'
import type { ToasterProps } from 'sonner'

const Toaster = ({ ...props }: ToasterProps) => {
  const { theme = 'system' } = useTheme()

  return (
    <Sonner
      theme={'dark'}
      className="toaster group "
      toastOptions={{
        classNames: {
          toast: `group !shadow-lg !shadow-black !border-2 !border-black !text-brand-text !tracking-small !text-base`,
          default: '!bg-brand-primary-dark',
          success: '!bg-green-900',
          error: '!bg-red-900',
        },
      }}
      {...props}
    />
  )
}

export { Toaster }
