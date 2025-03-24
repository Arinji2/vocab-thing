import { QueryErrorResetBoundary } from "@tanstack/react-query";
import { ErrorBoundary } from "react-error-boundary";
import { Button } from "~/components/ui/button";
export function ErrorWrapper({ children }: { children: React.ReactNode }) {
  return (
    <QueryErrorResetBoundary>
      {({ reset }) => (
        <ErrorBoundary
          onReset={reset}
          fallbackRender={({ error, resetErrorBoundary }) => (
            <div className="flex h-full w-full flex-col items-center justify-center gap-4 rounded-xl bg-brand-destructive-dark/50 p-4">
              <h2 className="text-2xl font-semibold tracking-small text-brand-text">
                Something went wrong!
              </h2>
              <p className="text-lg text-brand-text">{error?.message}</p>
              <Button variant={"default"} onClick={resetErrorBoundary}>
                Try Again
              </Button>
            </div>
          )}
        >
          {children}
        </ErrorBoundary>
      )}
    </QueryErrorResetBoundary>
  );
}
