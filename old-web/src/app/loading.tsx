import { Loader2 } from "lucide-react";

export default function Loading() {
  return (
    <div className="w-full h-full flex flex-col items-center justify-center bg-brand-background">
      <Loader2 className="animate-spin text-brand-accent xl:size-28 md:size-24 size-20" />
    </div>
  );
}
