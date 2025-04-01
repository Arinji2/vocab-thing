import { Button } from "@/components/ui/button";
import { cookies } from "next/headers";
import Link from "next/link";

export async function LoginButton() {
  const cookieStore = await cookies();
  const sessionExists = cookieStore.has("session");
  return (
    <Button asChild variant={"secondary"}>
      <Link href={sessionExists ? "/dashboard" : "/login"}>
        {sessionExists ? "Dashboard" : "Get Started"}
      </Link>
    </Button>
  );
}
