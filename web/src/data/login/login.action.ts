"use server";
import { LoginProviders, OauthCallbackURLSchema } from "./login";
import { getApiURL } from "@/data/utils/getApiURL";
import { setServerCookies } from "../utils/setServerCookie";

export type OauthCallbackURLActionState =
  | { providerType: LoginProviders }
  | {
      providerType: LoginProviders;
      success: true;
      data: { codeURL: string };
    }
  | {
      providerType: LoginProviders;
      success: false;
      error: string;
    };

export async function OauthCallbackURLAction(
  previousState: OauthCallbackURLActionState,
) {
  "use server";
  const provider = previousState.providerType;
  try {
    const apiURL = getApiURL();
    const body = JSON.stringify({ providerType: provider });
    const res = await fetch(`${apiURL}/oauth/generate-code-url`, {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body,
    });

    await setServerCookies(res);

    if (!res.ok) {
      throw new Error("Failed to get callback url");
    }

    const data = await res.json();

    return {
      ...previousState,
      success: true,
      data: OauthCallbackURLSchema.parse(data),
    };
  } catch (error) {
    console.error("OAuth callback URL error:", error);
    return {
      ...previousState,
      success: false,
      error: error instanceof Error ? error.message : "Unknown error occurred",
    };
  }
}

export async function LoginAsGuestAction() {
  "use server";
  try {
    const apiURL = getApiURL();

    const res = await fetch(`${apiURL}/user/create/guest`, {
      method: "POST",
      cache: "no-store",
    });

    await setServerCookies(res);

    if (!res.ok) {
      throw new Error("Failed to login as guest");
    }

    return { success: true };
  } catch (error) {
    console.error("Guest login error:", error);
    return {
      success: false,
      error: error instanceof Error ? error.message : "Unknown error occurred",
    };
  }
}
