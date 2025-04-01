import Features from "@/app/(home-navbar)/(home)/features";
import Footer from "@/app/(home-navbar)/(home)/footer";
import Hero from "@/app/(home-navbar)/(home)/hero";
import Info from "@/app/(home-navbar)/(home)/info";
import Works from "@/app/(home-navbar)/(home)/works";
import Loading from "@/app/loading";
import { Suspense } from "react";

export default async function Page() {
  return (
    <div className="flex h-fit w-full pb-10 gap-20 screen-padding flex-col items-center justify-start">
      <Suspense fallback={<Loading />}>
        <Hero />
      </Suspense>
      <Works />
      <Features />
      <Info />
      <Footer />
    </div>
  );
}
