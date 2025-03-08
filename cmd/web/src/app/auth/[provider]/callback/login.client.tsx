"use client";
import { useQuery } from "@tanstack/react-query";

export function LoginClient({
  provider,
  code,
  state,
  ip,
}: {
  provider: string;
  code: string;
  state: string;
  ip: string;
}) {
  const sendOAuthCallback = async () => {
    const response = await fetch("http://localhost:8080/oauth/callback", {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify({
        providerType: provider,
        code,
        state,
        fingerprint: navigator.userAgent,
        ip,
      }),
      credentials: "include",
    });
    if (!response.ok) throw new Error(response.statusText);
  };

  const { error, isLoading } = useQuery({
    queryKey: ["oauth", provider, code],
    queryFn: sendOAuthCallback,
    enabled: !!code,
    retry: false,
  });

  return (
    <div>
      {isLoading && <p>Authenticating...</p>}
      {error && <p>Error: {error.message}</p>}
      {!isLoading && !error && <p>Authenticated!</p>}
    </div>
  );
}
