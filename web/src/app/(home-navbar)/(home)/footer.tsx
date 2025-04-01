import { Button } from "@/components/ui/button";
import OptimizedImage from "@/utils/image";
import Link from "next/link";

export default function Footer() {
  return (
    <div className="gap-10 h-full-navbar relative flex  w-full flex-col items-center justify-center">
      <div className="w-full h-fit flex flex-col items-center z-10 justify-center gap-10">
        <h1 className="text-4xl font-semibold tracking-large text-brand-primary md:text-6xl">
          VOCAB THING
        </h1>
        <p className="text-brand-text font-medium text-xl md:text-2xl text-center">
          {'"Snag cool words now, flex them later"'}
        </p>
        <div className="flex h-fit w-fit flex-row items-center justify-start gap-6">
          <Button asChild variant={"default"} size={"lg"}>
            <Link href="/login">Get Started</Link>
          </Button>
          <Button variant={"secondary"} size={"lg"}>
            Add Extension
          </Button>
        </div>
      </div>

      <OptimizedImage
        fill
        srcLocation="/home/words/words"
        alt="Word Background"
        decoding="async"
      />
    </div>
  );
}
