import Hero from "@/app/(home)/hero";
import Works from "@/app/(home)/works";

export default function Page() {
  return (
    <div className="flex h-fit w-full pb-10 gap-20 px-4 md:px-2 2xl:px-0 flex-col items-center justify-start">
      <Hero />
      <Works />
    </div>
  );
}
