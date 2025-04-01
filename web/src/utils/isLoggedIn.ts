import { cookies } from "next/headers";
import { cache } from "react";

export async function isLoggedIn(): Promise<boolean> {
  return await cache(async () => {
    const cookieStore = await cookies();
    return cookieStore.get("oauth_session") !== null;
  })();
}
