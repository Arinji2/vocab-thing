import { Button } from "@/components/ui/button";
import Link from "next/link";

export function Navbar() {
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
      <Button asChild variant={"secondary"}>
        <Link href="/login">Get Started</Link>
      </Button>
    </div>
  );
}
