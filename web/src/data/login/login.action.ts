"use server";
import { LoginProvidersType, OauthCallbackURLSchema } from "./login";
import { getApiURL } from "@/data/utils/getApiURL";
import { setServerCookies } from "../utils/setServerCookie";
import { cookies } from "next/headers";
import { ErrorResponseSchema, HandleResponseError } from "@/data/errors";

export type OauthCallbackURLActionState =
  | { providerType: LoginProvidersType }
  | {
      providerType: LoginProvidersType;
      success: true;
      data: { codeURL: string };
    }
  | {
      providerType: LoginProvidersType;
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
      credentials: "include",
      body,
    });

    await setServerCookies(res);
    const resError = await HandleResponseError("OAuth callback URL", res);
    if (resError) {
      return {
        ...previousState,
        success: false,
        error: resError.readable,
      };
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
      credentials: "include",
      cache: "no-store",
    });

    await setServerCookies(res);

    const resError = await HandleResponseError("Guest Login", res);
    if (resError) {
      return {
        success: false,
        error: resError.readable,
      };
    }
    return { success: true };
  } catch (error) {
    console.error("Guest login error:", error);
    return {
      success: false,
      error: "Unknown error occurred",
    };
  }
}

export type LoginWithSocialActionState =
  | {
      providerType: LoginProvidersType;
      code: string;
      state: string;
      fingerprint: string;
      ip: string;
    }
  | {
      providerType: LoginProvidersType;
      code: string;
      state: string;
      fingerprint: string;
      ip: string;
      success: true;
      data: { codeURL: string };
    }
  | {
      providerType: LoginProvidersType;
      code: string;
      state: string;
      fingerprint: string;
      ip: string;
      success: false;
      error: string;
    };
export async function LoginWithSocialAction(
  previousState: LoginWithSocialActionState,
) {
  "use server";
  try {
    const { providerType, code, state, fingerprint, ip } = previousState;
    const apiURL = getApiURL();

    const cookieStore = await cookies();
    const reqCookies = cookieStore.toString();

    const body = JSON.stringify({ providerType, code, state, fingerprint, ip });
    const res = await fetch(`${apiURL}/oauth/callback`, {
      method: "POST",
      headers: { "Content-Type": "application/json", cookie: reqCookies },
      credentials: "include",
      body,
    });

    await setServerCookies(res);
    const resError = await HandleResponseError("Social Login", res);
    if (resError) {
      return {
        ...previousState,
        success: false,
        error: resError.readable,
      };
    } else
      return {
        ...previousState,
        success: true,
      };
  } catch (error) {
    console.error("Social login error:", error);
    return {
      ...previousState,
      success: false,
      error: "Unknown error occurred",
    };
  } finally {
    const cookieStore = await cookies();
    if (cookieStore.get("oauth_state")) {
      cookieStore.delete("oauth_state");
    }
    if (cookieStore.get("oauth_session")) {
      cookieStore.delete("oauth_session");
    }
  }
}
