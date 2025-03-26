import { createFileRoute } from "@tanstack/react-router";
import { Suspense } from "react";
import { ErrorWrapper } from "~/components/ErrorWrapper";
import OptimizedImage from "~/components/image";
import { Words, WordsLoading } from "~/components/routes/home/word";
import { Button } from "~/components/ui/button";

export const Route = createFileRoute("/")({
  component: Home,
});

function Home() {
  return (
    <div className="flex h-fit w-full min-w-1 flex-col items-center justify-start">
      <div className="xl:h-full-navbar relative flex h-fit w-full flex-col items-start justify-stretch">
        <div className="screen-padding relative z-20 flex h-full w-full flex-col items-start justify-start gap-14 py-8 md:gap-6 xl:py-12">
          <h1 className="text-2xl font-medium tracking-large text-brand-text md:text-4xl">
            TANSTACK START
          </h1>
          <div className="flex h-fit w-fit flex-row items-center justify-start gap-6">
            <Button variant={"default"}>Get Started</Button>
            <Button variant={"secondary"}>Add Extension</Button>
          </div>
          <div className="flex h-fit w-fit flex-row items-center justify-start gap-6">
            <ul className="flex h-fit w-fit flex-row flex-wrap items-center justify-start gap-6">
              <li className="relative whitespace-nowrap pl-7 text-lg tracking-small text-brand-text before:absolute before:left-0 before:top-1/2 before:size-5 before:-translate-y-1/2 before:rounded-full before:bg-brand-accent">
                Offline Support
              </li>
              <li className="relative whitespace-nowrap pl-7 text-lg tracking-small text-brand-text before:absolute before:left-0 before:top-1/2 before:size-5 before:-translate-y-1/2 before:rounded-full before:bg-brand-accent">
                Chrome and Firefox Extension
              </li>
              <li className="relative whitespace-nowrap pl-7 text-lg tracking-small text-brand-text before:absolute before:left-0 before:top-1/2 before:size-5 before:-translate-y-1/2 before:rounded-full before:bg-brand-accent">
                Unlimited Vocab
              </li>
            </ul>
          </div>
          <div className="mt-auto flex h-full max-h-[250px] w-full">
            <div className="w-full max-w-full overflow-hidden">
              <ErrorWrapper>
                <Suspense fallback={<WordsLoading />}>
                  <Words />
                </Suspense>
              </ErrorWrapper>
            </div>
          </div>
        </div>

        <OptimizedImage
          fill
          srcLocation="/words/words"
          alt="Word Background"
          decoding="async"
        />
      </div>

      <div className="h-[100svh]"></div>
    </div>
  );
}
