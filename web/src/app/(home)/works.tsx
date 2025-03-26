export default function Works() {
  return (
    <div className="w-full h-[50svh] flex flex-col items-center justify-start gap-10">
      <h2 className="font-semibold text-4xl text-center md:text-5xl tracking-large text-brand-text">
        How It Works
      </h2>
      <div className="w-full h-fit flex flex-row items-center justify-center gap-20 md:flex-nowrap flex-wrap">
        <WorkCard
          index={1}
          description={
            "Find cool words, phrases and other vocabulary from all over the internet"
          }
        />

        <WorkCard
          index={2}
          description={
            "Save them to Vocab Thing with a description for yourself to use in the future"
          }
        />

        <WorkCard
          index={3}
          description={
            "Search through your saved vocab whenever you want, blazingly fast"
          }
        />
      </div>
    </div>
  );
}

function WorkCard({
  index,
  description,
}: {
  index: number;
  description: string;
}) {
  return (
    <div className="flex py-4 min-h-[200px] rounded-xl bg-brand-secondary-dark relative w-full flex-col items-center justify-center px-3 gap-10 ">
      <div className="size-10 md:size-12 rounded-full bg-brand-primary flex flex-col items-center absolute -top-5 md:-top-6 justify-center">
        {
          <p className="font-medium text-brand-background tracking-small text-2xl md:text-3xl">
            {index}
          </p>
        }
      </div>
      <p className="text-brand-text text-lg md:text-xl  tracking-small text-center">
        {description}
      </p>
    </div>
  );
}
