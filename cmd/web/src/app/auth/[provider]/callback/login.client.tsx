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
    return response.json();
  };

  const { data, error, isLoading } = useQuery({
    queryKey: ["oauth", provider, code],
    queryFn: sendOAuthCallback,
    enabled: !!code,
  });

  return (
    <div>
      {isLoading && <p>Authenticating...</p>}
      {error && <p>Error: {error.message}</p>}
      {data && (
        <div>
          <h3>OAuth Data:</h3>
          <pre>{JSON.stringify(data, null, 2)}</pre>
        </div>
      )}
    </div>
  );
}
