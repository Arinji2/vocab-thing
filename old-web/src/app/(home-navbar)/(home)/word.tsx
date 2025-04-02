"use cache";
import { WordCard } from "@/app/(home-navbar)/(home)/word.client";
import { Button } from "@/components/ui/button";
import { Skeleton } from "@/components/ui/skeleton";
import { fetchWords } from "@/data/words";
import { unstable_cacheLife as cacheLife } from "next/cache";

export async function Words({ isLoggedIn }: { isLoggedIn: boolean }) {
  cacheLife("hours");
  const data = await fetchWords();
  return (
    <div className="flex h-full w-full flex-col items-start justify-start gap-2">
      <p className="text-lg font-medium tracking-small text-brand-text">
        Powered By{" "}
        <a
          href="https://sense.arinji.com"
          target="_blank"
          rel="noreferrer"
          className="underline decoration-brand-offwhite underline-offset-4"
        >
          <span className="text-green-500">Sense</span> Or{" "}
          <span className="text-red-500">Nonsense</span>
        </a>
      </p>
      <div className="flex h-full w-full flex-row items-center snap-x snap-proximity justify-start gap-10 overflow-x-auto">
        {data.map((d) => (
          <WordCard data={d} key={d.id} isLoggedIn={isLoggedIn} />
        ))}
      </div>
    </div>
  );
}

export async function WordsLoading() {
  return (
    <div className="flex h-full w-full flex-col items-start justify-start gap-2">
      <p className="text-lg font-medium tracking-small text-brand-text">
        Powered By{" "}
        <a
          href="https://sense.arinji.com"
          target="_blank"
          rel="noreferrer"
          className="underline decoration-brand-offwhite underline-offset-4"
        >
          <span className="text-green-500">Sense</span> Or{" "}
          <span className="text-red-500">Nonsense</span>
        </a>
      </p>
      <div className="flex h-full w-full flex-row items-center snap-x snap-proximity justify-start gap-10 overflow-x-auto">
        {[...Array(3)].map((_, i) => (
          <WordSuspenseCard key={i} />
        ))}
      </div>
    </div>
  );
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
      <Button variant={"default"} disabled className="mt-auto" size={"sm"}>
        Add to Vocab
      </Button>
    </article>
  );
}
