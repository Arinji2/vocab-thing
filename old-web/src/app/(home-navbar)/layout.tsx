import { HomeNavbar } from "@/app/(home-navbar)/(navbar)/navbar";

export default function Layout({ children }: { children: React.ReactNode }) {
  return (
    <div className="w-full h-full flex flex-col items-center justify-start">
      <HomeNavbar />
      {children}
    </div>
  );
}
