import { headers } from "next/headers";
import { LoginClient } from "./login.client";

export default async function Page({
  params,
  searchParams,
}: {
  params: Promise<{ provider: string }>;
  searchParams: Promise<{
    code: string;
    state: string;
  }>;
}) {
  const { provider } = await params;
  const { code, state } = await searchParams;
  const headerStore = await headers();
  const IP =
    (headerStore.get("x-forwarded-for") ||
      headerStore.get("cf-connecting-ip")) ??
    "";

  return (
    <div className="w-full h-[100svh] flex flex-col items-center justify-center">
      <LoginClient provider={provider} code={code} state={state} ip={IP} />
    </div>
  );
}
