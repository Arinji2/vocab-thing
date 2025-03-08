"use client";

export default function HomeClient() {
  return (
    <div className="">
      <button
        onClick={async () => {
          const response = await fetch(
            "http://localhost:8080/oauth/generate-code-url",
            {
              method: "POST",
              credentials: "include",
              headers: {
                "Content-Type": "application/json",
              },
              body: JSON.stringify({
                providerType: "google",
              }),
            },
          );
          const url = await response.text();
          window.location.href = url;
        }}
      >
        Google Login
      </button>
    </div>
  );
}
