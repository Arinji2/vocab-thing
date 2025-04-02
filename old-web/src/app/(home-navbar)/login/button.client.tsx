"use client";
import { Button } from "@/components/ui/button";
import {
  LoginAsGuestAction,
  OauthCallbackURLAction,
  OauthCallbackURLActionState,
} from "@/data/login/login.action";
import { LoginProvidersType } from "@/data/login/login";
import { cn } from "@/utils/cn";
import { formatCapitalize } from "@/utils/format";
import { ClassValue } from "clsx";
import { useRouter } from "next/navigation";
import { startTransition, useActionState, useEffect } from "react";
import { Loader2 } from "lucide-react";

export default function LoginButton({
  provider,
  className,
}: {
  provider: LoginProvidersType;
  className?: ClassValue;
}) {
  const router = useRouter();
  const [oauthState, oauthAction, isOauthLoading] =
    useActionState<OauthCallbackURLActionState>(OauthCallbackURLAction, {
      providerType: provider,
    });

  const [guestState, guestAction, isGuestLoading] = useActionState(
    LoginAsGuestAction,
    null,
  );

  useEffect(() => {
    if (!oauthState && !guestState) {
      return;
    }
    if (oauthState && "success" in oauthState) {
      if (oauthState.success) {
        router.push(oauthState.data.codeURL);
      } else {
        throw new Error(oauthState.error);
      }
    }

    if (guestState && "success" in guestState) {
      if (guestState.success) {
        router.push("/dashboard");
      } else {
        throw new Error(guestState.error);
      }
    }
  }, [oauthState, guestState, router]);
  return (
    <Button
      onClick={() => {
        startTransition(() => {
          if (provider === "guest") guestAction();
          else oauthAction();
        });
      }}
      className={cn(
        "relative flex w-full flex-row items-center justify-center overflow-hidden gap-2 text-base bg-brand-primary-dark text-brand-text ",
        className,
      )}
    >
      <div
        className={cn(
          " flex flex-col items-center justify-center w-full h-full absolute top-0 left-0 transition-all ease-in-out duration-200 -translate-y-full bg-brand-primary-dark",
          {
            "translate-y-0": isGuestLoading || isOauthLoading,
          },
        )}
      >
        <Loader2 className="text-brand-text size-9 animate-spin" />
      </div>
      Login {provider === "guest" ? "as" : "with"} {formatCapitalize(provider)}
    </Button>
  );
}
