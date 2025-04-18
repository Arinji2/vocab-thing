"use client";
import { NextFont } from "next/dist/compiled/@next/font";
import { Toaster as Sonner } from "sonner";

type ToasterProps = React.ComponentProps<typeof Sonner> & {
  font: NextFont;
};

const Toaster = ({ ...props }: ToasterProps) => {
  return (
    <Sonner
      theme={"dark"}
      className="toaster group "
      toastOptions={{
        classNames: {
          toast: `${props.font.className} group !shadow-lg !shadow-black !border-2 !border-black !text-brand-text !tracking-small !text-base`,
          default: "!bg-brand-primary-dark",
          success: "!bg-green-900",
          error: "!bg-red-900",
        },
      }}
      {...props}
    />
  );
};

export { Toaster };
