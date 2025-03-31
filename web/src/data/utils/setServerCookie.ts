import "server-only";
import { cookies } from "next/headers";
import { parse } from "cookie";

export async function setServerCookies(response: Response) {
  const setCookies = response.headers.getSetCookie();
  if (!setCookies || setCookies.length === 0) return;

  const cookieStore = await cookies();

  setCookies.forEach((cookieString) => {
    const parsed = parse(cookieString);
    const entries = Object.entries(parsed);
    if (entries.length === 0) return;

    const [name, value] = entries[0];

    const attributes: Record<string, string | number | boolean | Date> = {};

    for (const [key, val] of entries) {
      if (key !== name && val !== undefined) {
        // Convert standard cookie attributes to the format expected by Next.js
        if (key.toLowerCase() === "max-age" && val) {
          attributes["maxAge"] = parseInt(val, 10);
        } else if (key.toLowerCase() === "expires" && val) {
          attributes["expires"] = new Date(val);
        } else if (
          key.toLowerCase() === "secure" ||
          key.toLowerCase() === "httponly"
        ) {
          attributes[key.toLowerCase()] = true;
        } else if (key.toLowerCase() === "samesite" && val) {
          attributes["sameSite"] = val.toLowerCase();
        } else {
          attributes[key.toLowerCase()] = val;
        }
      }
    }

    cookieStore.set({
      name,
      value: value || "",
      ...attributes,
    });
  });
}
