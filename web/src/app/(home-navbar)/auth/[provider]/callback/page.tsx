import Login from "./login.client";
import { LoginProviders } from "@/data/login/login";
import { formatCapitalize } from "@/utils/format";
import { headers } from "next/headers";
import { redirect } from "next/navigation";
import { unstable_noStore as noStore } from "next/cache";
import { Suspense } from "react";
export default function Page({
  params,
  searchParams,
}: {
  params: Promise<{ provider: string }>;
  searchParams: Promise<{ state: string; code: string }>;
}) {
  return (
    <Suspense>
      <SuspensedPage params={params} searchParams={searchParams} />
    </Suspense>
  );
}

async function SuspensedPage({
  params,
  searchParams,
}: {
  params: Promise<{ provider: string }>;
  searchParams: Promise<{ state: string; code: string }>;
}) {
  noStore();
  const { provider } = await params;
  const { state, code } = await searchParams;
  if (!state || !code) {
    return redirect("/login");
  }
  const validatedProvider = LoginProviders.safeParse(provider);
  if (!validatedProvider.success) {
    console.error("Invalid provider:", validatedProvider.error);
    return redirect("/login");
  }
  if (validatedProvider.data === "guest") {
    return redirect("/login");
  }

  const headerStore = await headers();
  const userIP =
    headerStore.get("x-vercel-forwarded-for") ||
    headerStore.get("x-forwarded-for") ||
    "0.0.0.0";
  const userAgent = headerStore.get("user-agent") ?? "Unknown";
  return (
    <div className="flex flex-col items-center justify-center py-4 gap-10 w-full h-full-navbar screen-padding ">
      <h1 className="text-3xl text-center md:text-5xl font-bold text-brand-text tracking-large">
        Logging In With{" "}
        <span className="text-brand-accent">
          {formatCapitalize(validatedProvider.data)}
        </span>
      </h1>
      <Login
        providerType={validatedProvider.data}
        code={code}
        state={state}
        ip={userIP}
        fingerprint={userAgent}
      />
    </div>
  );
}
