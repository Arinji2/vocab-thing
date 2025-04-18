"use client";

import type { LoginProvidersType } from "@/data/login/login";
import {
  LoginWithSocialAction,
  type LoginWithSocialActionState,
} from "@/data/login/login.action";
import { Loader2 } from "lucide-react";
import { useRouter } from "next/navigation";
import { startTransition, useActionState, useEffect, useState } from "react";

export default function Login({
  providerType,
  code,
  state,
  ip,
  fingerprint,
}: {
  providerType: LoginProvidersType;
  code: string;
  state: string;
  ip: string;
  fingerprint: string;
}) {
  const router = useRouter();
  const [timeElapsed, setTimeElapsed] = useState(0);
  const [actionState, action, isPending] =
    useActionState<LoginWithSocialActionState>(LoginWithSocialAction, {
      providerType: providerType,
      code: code,
      state: state,
      ip: ip,
      fingerprint: fingerprint,
    });

  useEffect(() => {
    const shouldTriggerAction =
      !isPending && !("success" in actionState) && !("error" in actionState);

    if (shouldTriggerAction) {
      startTransition(() => {
        action();
      });
    }

    if ("success" in actionState) {
      if (actionState.success) {
        router.push("/dashboard");
      } else {
        router.push("/login");
      }
    }
  }, [router, action, actionState, isPending]);

  useEffect(() => {
    const timer = setTimeout(() => {
      setTimeElapsed(timeElapsed + 1);
    }, 1000);
    if (timeElapsed > 60) {
      setTimeout(() => {
        router.push("/login");
      }, 2000);
    }
    return () => clearTimeout(timer);
  }, [timeElapsed, router]);

  return (
    <div className="w-full h-fit flex flex-col items-center gap-6 justify-center">
      <div className="text-brand-text w-full h-fit flex flex-col items-center justify-center gap-6">
        <div className="text-xl font-medium tracking-large w-full h-fit flex flex-row items-center justify-center gap-2">
          <Loader2
            className="animate-spin text-brand-accent"
            size={25}
            strokeWidth={3}
          />{" "}
          Authenticating Securely
        </div>
        <p>Time Elapsed: {timeElapsed} seconds</p>
      </div>
      {!isPending && "error" in actionState && (
        <div className="text-red-500">
          Authentication failed. Redirecting to login...
        </div>
      )}
      {timeElapsed > 60 && (
        <div className="flex flex-col items-center justify-center">
          <p className="text-xl text-brand-destructive-light">
            Authentication timed out. Redirecting to login...
          </p>
        </div>
      )}
    </div>
  );
}
