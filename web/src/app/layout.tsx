import { Tektur } from "next/font/google";
import type React from "react";
import "./globals.css";

const tektur = Tektur({
  subsets: ["latin"],
  weight: ["400", "500", "600", "700", "800", "900"],
});

export const metadata = {
  title: "VocabThing",
  description:
    "Save words and phrases you find on the internet, and use them in the future effortlessly",
  keywords:
    "vocab, vocabthing, arinji, arinji.com, arinjay dhar, save words, phrases",
  themeColor: "#89DFE9",
  icons: {
    icon: [
      { url: "/metadata/favicon-16x16.png", sizes: "16x16", type: "image/png" },
      { url: "/metadata/favicon-32x32.png", sizes: "32x32", type: "image/png" },
      { url: "/metadata/favicon-96x96.png", sizes: "96x96", type: "image/png" },
      { url: "/metadata/favicon.svg", type: "image/svg+xml" },
    ],
    shortcut: [{ url: "/metadata/favicon.ico" }],
    apple: [{ url: "/metadata/apple-touch-icon.png", sizes: "180x180" }],
  },
  manifest: "/metadata/site.webmanifest",
};

export default function RootLayout({
  children,
}: Readonly<{
  children: React.ReactNode;
}>) {
  return (
    <html lang="en">
      <body className={`${tektur.className} antialiased`}>{children}</body>
    </html>
  );
}
