"use client";

export default function HomeClient() {
  return (
    <div className=" flex flex-col items-center justify-center gap-10">
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
                providerType: "discord",
              }),
            },
          );
          const url = await response.text();
          window.location.href = url;
        }}
      >
        Discord Login
      </button>
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
                providerType: "github",
              }),
            },
          );
          const url = await response.text();
          window.location.href = url;
        }}
      >
        Github Login
      </button>
      <button
        onClick={async () => {
          await fetch("http://localhost:8080/user/create/guest", {
            method: "POST",
            credentials: "include",
            headers: {
              "Content-Type": "application/json",
            },
          });
        }}
      >
        Guest Login
      </button>
    </div>
  );
}
