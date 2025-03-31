import LoginButton from "@/app/login/button.client";
import { ErrorWrapper } from "@/components/ui/error-boundary";
import { CheckCircle2 } from "lucide-react";

export default function Page() {
  return (
    <div className="h-fit flex flex-col items-center justify-center py-4 gap-10 w-full xl:h-full-navbar screen-padding ">
      <h1 className="text-3xl text-center md:text-5xl font-bold text-brand-text tracking-large">
        Continue To Vocab Thing
      </h1>
      <div className="w-full gap-10 h-fit flex flex-col xl:flex-row items-center xl:items-start  justify-between">
        <div className="w-full xl:w-[550px] md:w-[80%] rounded-lg shadow-black shadow-lg h-fit xl:h-[400px] gap-8 bg-brand-secondary-dark flex flex-col items-start justify-start px-8 py-8">
          <h2 className="text-2xl font-medium text-brand-text tracking-large">
            Login With Socials
          </h2>
          <div className="flex flex-col items-start justify-start gap-3">
            <div className="w-full h-fit flex flex-row items-center justify-start gap-3">
              <CheckCircle2 className="text-brand-accent" />
              <p className="text-brand-text  text-lg">Free Forever</p>
            </div>
            <div className="w-full h-fit flex flex-row items-center justify-start gap-3">
              <CheckCircle2 className="text-brand-accent" />
              <p className="text-brand-text  text-lg">Access To All Features</p>
            </div>
            <div className="w-full h-fit flex flex-row items-center justify-start gap-3">
              <CheckCircle2 className="text-brand-accent" />
              <p className="text-brand-text  text-lg">15 AI Usage Per Day</p>
            </div>
          </div>
          <div className="w-full h-fit flex flex-col items-center justify-center">
            <ErrorWrapper>
              <div className="w-full h-fit grid md:grid-cols-2 grid-cols-1 mt-auto gap-4">
                <LoginButton provider="google" />
                <LoginButton provider="discord" />
                <LoginButton provider="github" className="md:col-span-2" />
              </div>
            </ErrorWrapper>
          </div>
        </div>
        <div className="md:w-[80%] w-full xl:w-[550px] rounded-lg shadow-black shadow-lg h-fit gap-8 bg-brand-offwhite-dark flex flex-col items-start justify-start px-8 py-8">
          <h2 className="text-2xl font-medium text-brand-text tracking-large">
            Login As Guest
          </h2>
          <div className="flex flex-col items-start justify-start gap-3">
            <div className="w-full h-fit flex flex-row items-center justify-start gap-3">
              <CheckCircle2 className="text-brand-accent" />
              <p className="text-brand-text  text-lg">Free Forever</p>
            </div>
            <div className="w-full h-fit flex flex-row items-center justify-start gap-3">
              <CheckCircle2 className="text-brand-accent" />
              <p className="text-brand-text  text-lg">Access To All Features</p>
            </div>
            <div className="w-full h-fit flex flex-row items-center justify-start gap-3">
              <CheckCircle2 className="text-brand-accent" />
              <p className="text-brand-text  text-lg">2 AI Usage Per Day</p>
            </div>
          </div>
          <ErrorWrapper>
            <div className="w-full h-fit flex flex-col items-center justify-center">
              <LoginButton provider="guest" />
            </div>
          </ErrorWrapper>
        </div>
      </div>
    </div>
  );
}
