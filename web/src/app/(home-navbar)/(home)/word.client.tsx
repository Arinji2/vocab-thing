"use client";

import { Button } from "@/components/ui/button";
import { WordSchemaType } from "@/data/words";
import { useRouter } from "next/navigation";

export function WordCard({
  data,
  isLoggedIn,
}: {
  data: WordSchemaType;
  isLoggedIn: boolean;
}) {
  const router = useRouter();
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
          localStorage.setItem("homepage-word", data.word);
          localStorage.setItem("homepage-definition", data.definition);
          if (isLoggedIn) {
            router.push("/dashboard");
          } else {
            router.push("/login");
          }
        }}
        variant={"default"}
        className="mt-auto"
        size={"sm"}
      >
        Add to Vocab
      </Button>
    </article>
  );
}
