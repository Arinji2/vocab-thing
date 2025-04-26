import { createContext, useContext, useState } from 'react'

type SyncingContextType = {
  isSyncing: boolean
  setIsSyncing: (isSyncing: boolean) => void
}

const SyncingContext = createContext<SyncingContextType | undefined>(undefined)

export const SyncingProvider = ({
  children,
}: {
  children: React.ReactNode
}) => {
  const [isSyncing, setIsSyncing] = useState(false)

  return (
    <SyncingContext.Provider value={{ isSyncing, setIsSyncing }}>
      {children}
    </SyncingContext.Provider>
  )
}

export const useSyncing = () => {
  const context = useContext(SyncingContext)
  if (!context) {
    throw new Error('useSyncing must be used within a SyncingProvider')
  }
  return context
}
