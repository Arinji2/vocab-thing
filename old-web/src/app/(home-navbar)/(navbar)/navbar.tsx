import { Suspense } from "react";
import { LoginButton } from "./login";
import Link from "next/link";
import { Button } from "@/components/ui/button";

export function HomeNavbar() {
  return (
    <div className="h-navbar screen-padding sticky top-0 z-50 flex w-full flex-row items-center justify-between border-b-[4px] border-white/10 bg-brand-background">
      <Link href="/">
        <img
          src="/logo/logo.png"
          alt="logo"
          width={180}
          height={40}
          fetchPriority="high"
          loading="eager"
          className="hidden object-cover md:block"
        />

        <img
          src="/logo/logo-icon.png"
          alt="logo"
          width={50}
          height={30}
          fetchPriority="high"
          loading="eager"
          className="block object-cover md:hidden"
        />
      </Link>
      <Suspense
        fallback={
          <Button asChild disabled variant={"secondary"}>
            Get Started
          </Button>
        }
      >
        <LoginButton />
      </Suspense>
    </div>
  );
}
