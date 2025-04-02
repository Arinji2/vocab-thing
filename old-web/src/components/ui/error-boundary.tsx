"use client";

import { Button } from "@/components/ui/button";
import { cn } from "@/utils/cn";
import { ClassValue } from "clsx";
import { useState } from "react";
import { ErrorBoundary } from "react-error-boundary";

export function ErrorWrapper({
  children,
  className,
}: {
  children: React.ReactNode;
  className?: ClassValue;
}) {
  const [errorKey, setErrorKey] = useState(0); // Track key to force remount on error

  return (
    <ErrorBoundary
      onReset={() => setErrorKey((prev) => prev + 1)}
      fallbackRender={({ error, resetErrorBoundary }) => (
        <div className="flex h-full w-full flex-col items-center justify-center gap-4 rounded-xl bg-brand-destructive-dark/50 p-4">
          <h2 className="text-2xl font-semibold tracking-small text-brand-text">
            Something went wrong!
          </h2>
          <p className="text-lg text-brand-text">{error?.message}</p>
          <Button variant="default" onClick={resetErrorBoundary}>
            Try Again
          </Button>
        </div>
      )}
    >
      <div key={errorKey} className={cn("w-full h-full", className)}>
        {children}
      </div>
    </ErrorBoundary>
  );
}
